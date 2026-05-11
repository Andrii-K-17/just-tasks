package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andriik17/just-tasks/internal/middleware"
	"github.com/andriik17/just-tasks/internal/models"
	"github.com/andriik17/just-tasks/internal/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler manages authentication logic, database access, and JWT configuration.
type AuthHandler struct {
	db        *sqlx.DB
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthHandler initializes and returns a new AuthHandler.
func NewAuthHandler(db *sqlx.DB, jwtSecret string, jwtExpiry time.Duration) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: jwtSecret, jwtExpiry: jwtExpiry}
}

// issueTokenCookie generates a JWT and sets it as an HTTP-only cookie.
func (h *AuthHandler) issueTokenCookie(r http.ResponseWriter, userID int) error {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(h.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return err
	}

	http.SetCookie(r, &http.Cookie{
		Name:     "token",
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.jwtExpiry.Seconds()),
	})
	return nil
}

// clearTokenCookie removes the authentication cookie by expiring it.
func clearTokenCookie(r http.ResponseWriter) {
	http.SetCookie(r, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// authRequest represents the login or registration payload.
type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register creates a new user and returns an authentication token.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Username) < 3 || len(req.Password) < 8 {
		response.Error(w, http.StatusUnprocessableEntity,
			"username must be at least 3 characters and password at least 8 characters long")
		return
	}

	var exists bool
	err := h.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", req.Username)
	if err != nil {
	    response.Error(w, http.StatusInternalServerError, "database error")
	    return
	}
	if exists {
		response.Error(w, http.StatusConflict, "this username is already taken")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	var user models.User
	err = h.db.QueryRowx(
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)
		 RETURNING id, username, created_at`,
		req.Username, string(hash),
	).StructScan(&user)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := h.issueTokenCookie(w, user.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Login authenticates a user and provides a session cookie.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var user models.User
	if err := h.db.Get(&user,
		"SELECT id, username, password_hash FROM users WHERE username=$1", req.Username,
	); err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := h.issueTokenCookie(w, user.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Logout clears the user session cookie.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearTokenCookie(w)
	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// Me retrieves and returns the currently authenticated user's data.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var user models.User
	if err := h.db.Get(&user,
		"SELECT id, username FROM users WHERE id=$1", userID,
	); err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"id":       user.ID,
		"username": user.Username,
	})
}

// DeleteAccount removes a user from the database and clears their session.
func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if _, err := h.db.Exec("DELETE FROM users WHERE id=$1", userID); err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	clearTokenCookie(w)
	response.JSON(w, http.StatusOK, map[string]string{"message": "account deleted successfully"})
}
