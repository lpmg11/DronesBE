package services

import (
	"drones-be/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenServices struct{
	cfg *config.Config
}

func NewTokenService(cfg *config.Config) *TokenServices{
	return &TokenServices{
		cfg: cfg,
	}
}

func (s *TokenServices) GenerateToken(userID, role string) (string, error){

	accesToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"role": role,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte(s.cfg.JWTSecret))
	if err != nil{
		return "", err
	}


	return accesToken, nil

}

func (s *TokenServices) VerifyToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, userIDOk := claims["userID"].(string)
		role, roleOk := claims["role"].(string)
		if !userIDOk || !roleOk {
			return "", "", jwt.ErrInvalidKey
		}
		return userID, role, nil
	}

	return "", "", jwt.ErrInvalidKey
}