package service

import (
	"fmt"
	"os"
	"time"

	"github.com/SantiagoBedoya/otp-api/internal/dto"
	"github.com/SantiagoBedoya/otp-api/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) SignIn(data *dto.SignInDto) (string, error) {
	user, err := s.repo.GetByEmail(data.Email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return "", ErrInvalidPassword
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "otp-api",
		Subject:   fmt.Sprint(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *AuthService) GetUserByID(userID string) (*model.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) SetUser2FA(userID, secret string) error {
	return s.repo.SaveSecret(userID, secret)
}

func (s *AuthService) SetUser2FAValid(userID string) error {
	return s.repo.SetValidSecret(userID)
}

func (s *AuthService) SaveUser(data *dto.SignUpDto) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  string(hash),
	}
	err = s.repo.Save(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
