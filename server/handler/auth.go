package handler

import (
	"encoding/json"
	"food_delivery/repository"
	"food_delivery/request"
	"food_delivery/respond"
	"food_delivery/server/middlware"
	"food_delivery/service"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo         repository.UserRepository
	tokenService *service.TokenService
}

func NewAuthHandler(tokenService *service.TokenService, repo repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		tokenService: tokenService,
		repo:         repo,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	req := new(request.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessString, err := h.tokenService.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshString, err := h.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := respond.LoginRespond{
		AccessToken:  accessString,
		RefreshToken: refreshString,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	req := new(request.RegisterRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.repo.RegisterUser(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(200)

}

func (h *AuthHandler) ValidRefreshToken(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenService.ValidateRefreshToken(service.GetTokenFromBearerString(
		r.Header.Get("Authorization")),
	)

	if err != nil {
		http.Error(w, "invalid claims.credentials", http.StatusUnauthorized)
		return
	}

	if !claims.ExpiresAt.After(time.Now()) {
		http.Error(w, "the refersh token has expired", http.StatusUnauthorized)
		return
	}

	accessString, err := h.tokenService.GenerateAccessToken(claims.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshString, err := h.tokenService.GenerateRefreshToken(claims.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := respond.LoginRespond{
		AccessToken:  accessString,
		RefreshToken: refreshString,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	// // Type assertion to convert from `any` to `*service.JwtCustomClaims`
	claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	if !ok {
		// Handle the case where the type assertion fails
		log.Print("Failed to retrieve JWT claims from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// fmt.Println("from logout handler")
	err := h.repo.Logout(claims.ID)
	if err != nil {
		http.Error(w, "faild to invalidate the tokens", http.StatusInternalServerError)

	}

	w.WriteHeader(http.StatusOK)

}
