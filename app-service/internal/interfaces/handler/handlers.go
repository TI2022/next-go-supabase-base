package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TI2022/next-go-supabase-base/app-service/internal/application/usecase"
	"github.com/TI2022/next-go-supabase-base/app-service/internal/domain/entity"
	"github.com/TI2022/next-go-supabase-base/app-service/internal/infrastructure/persistence"
	"github.com/golang-jwt/jwt/v5"
)

type HealthHandler struct{}

func NewHealthHandler() http.Handler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type LoginHandler struct {
	db *sql.DB
}

func NewLoginHandler(db *sql.DB) http.Handler {
	return &LoginHandler{db: db}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userRepo := persistence.NewUserRepositoryPostgres(h.db)
	uc := usecase.NewLoginUsecase(userRepo)

	out, err := uc.Execute(r.Context(), usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Printf("login error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	token, err := createSessionToken(out.User.ID)
	if err != nil {
		log.Printf("create token error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // 本番では有効にする
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
}

type MeHandler struct {
	db *sql.DB
}

func NewMeHandler(db *sql.DB) http.Handler {
	return &MeHandler{db: db}
}

func (h *MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := parseSessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userRepo := persistence.NewUserRepositoryPostgres(h.db)
	uc := usecase.NewGetCurrentUserUsecase(userRepo)

	out, err := uc.Execute(r.Context(), usecase.GetCurrentUserInput{
		UserID: userID,
	})
	if err != nil {
		log.Printf("get current user error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if out.User == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	type response struct {
		ID    string  `json:"id"`
		Email string  `json:"email"`
		Name  *string `json:"name,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{
		ID:    string(out.User.ID),
		Email: out.User.Email,
		Name:  out.User.Name,
	})
}

// session token helpers

func sessionSecret() []byte {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "dev-session-secret-change-me"
	}
	return []byte(secret)
}

func createSessionToken(userID entity.UserID) (string, error) {
	claims := jwt.MapClaims{
		"sub": string(userID),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sessionSecret())
}

func parseSessionToken(tokenStr string) (entity.UserID, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return sessionSecret(), nil
	})
	if err != nil || !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrSignatureInvalid
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", jwt.ErrSignatureInvalid
	}
	return entity.UserID(sub), nil
}

