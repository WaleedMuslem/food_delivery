package service

import (
	"database/sql"
	"errors"
	"food_delivery/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

type TokenService struct {
	cfg *config.Config
	db  *sql.DB
}

func NewTokenService(cfg *config.Config, db *sql.DB) *TokenService {
	return &TokenService{
		cfg: cfg,
		db:  db,
	}
}

func (s *TokenService) GenerateAccessToken(userId uint) (token string, err error) {
	claims := &JwtCustomClaims{
		ID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(s.cfg.AccessLifetimeminutes))),
		},
	}

	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := tokenStruct.SignedString([]byte(s.cfg.AccessSecret))

	// Store the token in the database
	_, err = s.db.Exec(`
	  INSERT INTO jwt_tokens (user_id, token_string, revoked) 
	  VALUES ($1, $2, $3)
  `, userId, tokenString, false)

	if err != nil {
		return "", err
	}

	return tokenStruct.SignedString([]byte(s.cfg.AccessSecret))
}

func (s *TokenService) GenerateRefreshToken(userId uint) (token string, err error) {
	claims := &JwtCustomClaims{
		ID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(s.cfg.RefreshLifetimeminutes))),
		},
	}

	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenStruct.SignedString([]byte(s.cfg.RefreshSecret))
}

func (s TokenService) ValidateAccessToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		&JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.AccessSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("falid to parse token claims")
	}

	if s.isTokenRevoked(tokenString) {
		return nil, errors.New("you dont have access")
	}

	// log.Print(claims)

	return claims, nil
}

func (s TokenService) ValidateRefreshToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		&JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.RefreshSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("falid to parse token claims")
	}

	return claims, nil
}

func GetTokenFromBearerString(bearerString string) string {
	if bearerString == "" {
		return ""
	}

	parts := strings.Split(bearerString, " ")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) == 0 {
		return ""
	}

	return token

}

func (s *TokenService) isTokenRevoked(tokenstring string) bool {
	var revoked bool
	_ = s.db.QueryRow("SELECT revoked FROM jwt_tokens WHERE token_string = $1", tokenstring).Scan(&revoked)

	return revoked
}
