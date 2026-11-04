package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/farrelarea/budgeting-app/apps/api/internal/auth"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

// Login authenticates a user and returns a JWT.
func Login(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid body"})
			return
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Email == "" || req.Password == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "email and password required"})
			return
		}

		var id, name, hash string
		err := pool.QueryRow(r.Context(),
			"SELECT id, name, password_hash FROM users WHERE email = $1", req.Email).
			Scan(&id, &name, &hash)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid credentials"})
			return
		}
		if !auth.CheckPassword(hash, req.Password) {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid credentials"})
			return
		}

		token, err := auth.GenerateToken(id, req.Email, auth.SecretKey(), 0)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "token generation failed"})
			return
		}
		writeJSON(w, http.StatusOK, loginResponse{
			Token: token,
			User:  userResponse{ID: id, Email: req.Email, Name: name},
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
