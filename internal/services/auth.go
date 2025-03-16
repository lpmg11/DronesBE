package services

import (
	"drones-be/internal/config"
	"drones-be/internal/models"
	"drones-be/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthServices struct{
	pg *storage.PostgresClient
	cfg *config.Config
}

func NewAuthService(pg *storage.PostgresClient, cfg *config.Config) *AuthServices{
	return &AuthServices{
		pg: pg,
		cfg: cfg,
	}
}


func (s *AuthServices) RegisterUser(username, pasword, role string)(*models.User, error){
	

	hashedPasword, err := bcrypt.GenerateFromPassword([]byte(pasword), bcrypt.DefaultCost)
	if err != nil{
		return nil, err
	}

	var existingUser models.User
	err = s.pg.DB.Where("username = ?", username).First(&existingUser).Error
	if err == nil{
		return nil, err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPasword),
		Role: role,
	}

	err = s.pg.DB.Save(user).Error
	if err != nil{
		return nil, err
	}

	return user, nil

}

func (s *AuthServices) LoginUser(username, password string)(*models.User, error){
	var user models.User
	err := s.pg.DB.Where("username = ?", username).First(&user).Error
	if err != nil{
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil{
		return nil, err
	}

	return &user, nil
}

