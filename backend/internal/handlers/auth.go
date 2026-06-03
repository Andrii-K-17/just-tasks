package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Andrii-K-17/just-tasks/internal/middleware"
	"github.com/Andrii-K-17/just-tasks/internal/response"
	"github.com/Andrii-K-17/just-tasks/internal/services"
)

// AuthHandler manages HTTP authentication endpoints.
type AuthHandler struct {
	svc       *services.AuthService
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthHandler initializes and returns a new AuthHandler.
func NewAuthHandler(
	svc *services.AuthService,
	jwtSecret string,
	jwtExpiry time.Duration,
) *AuthHandler {
	return &AuthHandler{svc: svc, jwtSecret: jwtSecret, jwtExpiry: jwtExpiry}
}

// issueTokenCookie generates a JWT and sets it as an HTTP-only cookie.
func (h *AuthHandler) issueTokenCookie(w http.ResponseWriter, userID int) error {
	signed, err := services.IssueJWT(userID, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
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
func clearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
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

	user, err := h.svc.Register(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUsernameTaken) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
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

	user, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
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

	user, err := h.svc.GetByID(userID)
	if err != nil {
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

	if err := h.svc.DeleteAccount(userID); err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	clearTokenCookie(w)
	response.JSON(w, http.StatusOK, map[string]string{"message": "account deleted successfully"})
}
