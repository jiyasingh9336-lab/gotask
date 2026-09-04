package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sudarshanmg/gotask/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req RegisterRequest) (*User, error)
	Login(req LoginRequest) (string, string, error)
	Refresh(refreshToken string) (string, error)
	Logout(refreshToken string) error
	LogoutAll(userID int64) error
}

type authService struct {
	repo      AuthRepository
	jwtSecret string
	validator *validator.Validate
}

func NewService(repo AuthRepository, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: jwtSecret,
		validator: validator.New(),
	}
}

func generateSecureToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *authService) Register(req RegisterRequest) (*User, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, validation.FormatValidationError(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := "employee"
	id, err := s.repo.CreateUser(req.Username, string(hash), role)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Username: req.Username,
		Role:     role,
	}, nil
}

func (s *authService) Login(req LoginRequest) (string, string, error) {
	if err := s.validator.Struct(req); err != nil {
		return "", "", validation.FormatValidationError(err)
	}

	user, err := s.repo.FindByUsername(req.Username)
	if err != nil || user == nil {
		return "", "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	claims := Claims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken := generateSecureToken()
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)

	err = s.repo.SaveRefreshToken(user.ID, refreshToken, refreshExpiry)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Refresh(refreshToken string) (string, error) {
	rt, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil || rt == nil {
		return "", errors.New("invalid refresh token")
	}
	if rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
		return "", errors.New("refresh token expired or revoked")
	}

	// Create new access token
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(rt.UserID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (s *authService) Logout(refreshToken string) error {
	return s.repo.RevokeRefreshToken(refreshToken)
}

func (s *authService) LogoutAll(userID int64) error {
	return s.repo.RevokeAllRefreshTokens(userID)
}
