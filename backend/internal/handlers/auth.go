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
	svc           *services.AuthService
	jwtSecret     string
	jwtExpiry     time.Duration
	refreshExpiry time.Duration
	cookieSecure  bool
}

// NewAuthHandler initializes and returns a new AuthHandler.
func NewAuthHandler(
	svc *services.AuthService,
	jwtSecret string,
	jwtExpiry time.Duration,
	refreshExpiry time.Duration,
	cookieSecure bool,
) *AuthHandler {
	return &AuthHandler{
		svc:           svc,
		jwtSecret:     jwtSecret,
		jwtExpiry:     jwtExpiry,
		refreshExpiry: refreshExpiry,
		cookieSecure:  cookieSecure,
	}
}

// issueTokenCookies sets both the access JWT and refresh token as HTTP-only cookies.
func (h *AuthHandler) issueTokenCookies(w http.ResponseWriter, pair *services.TokenPair) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    pair.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.jwtExpiry.Seconds()),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		Path:     "/api/refresh",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.refreshExpiry.Seconds()),
	})
}

// clearTokenCookies removes both authentication cookies by expiring them.
func (h *AuthHandler) clearTokenCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/refresh",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
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

	user, pair, err := h.svc.Register(req.Username, req.Password, h.jwtSecret, h.jwtExpiry, h.refreshExpiry)
	if err != nil {
		if errors.Is(err, services.ErrUsernameTaken) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.issueTokenCookies(w, pair)

	response.JSON(w, http.StatusCreated, map[string]any{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Login authenticates a user and provides session cookies.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, pair, err := h.svc.Login(req.Username, req.Password, h.jwtSecret, h.jwtExpiry, h.refreshExpiry)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.issueTokenCookies(w, pair)

	response.JSON(w, http.StatusOK, map[string]any{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Refresh validates the refresh token cookie, rotates it, and issues a new token pair.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	pair, err := h.svc.Refresh(cookie.Value, h.jwtSecret, h.jwtExpiry, h.refreshExpiry)
	if err != nil {
		h.clearTokenCookies(w)
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.issueTokenCookies(w, pair)
	response.JSON(w, http.StatusOK, map[string]string{"message": "refreshed"})
}

// Logout clears the user session cookies and invalidates the refresh token.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.svc.Logout(cookie.Value)
	}
	h.clearTokenCookies(w)
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

	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.svc.Logout(cookie.Value)
	}

	if err := h.svc.DeleteAccount(userID); err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.clearTokenCookies(w)
	response.JSON(w, http.StatusOK, map[string]string{"message": "account deleted successfully"})
}
