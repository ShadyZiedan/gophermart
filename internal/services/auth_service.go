package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"

	"github.com/ShadyZiedan/gophermart/internal/models"
	"github.com/ShadyZiedan/gophermart/internal/security"
)

type AuthService struct {
	secretKey      string
	userRepository userRepository
}

func NewAuthService(secretKey string, userRepository userRepository) *AuthService {
	return &AuthService{secretKey: secretKey, userRepository: userRepository}
}

type userRepository interface {
	SaveUser(context.Context, *models.User) error
	FindUserByUsername(context.Context, string) (*models.User, error)
	IsUserExist(context.Context, string) (bool, error)
}

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
)

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	exists, err := s.userRepository.IsUserExist(ctx, username)
	if err != nil {
		return err
	}
	if exists {
		return ErrUsernameAlreadyTaken
	}
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return err
	}
	newUser := &models.User{
		Username: username,
		Password: hashedPassword,
	}
	return s.userRepository.SaveUser(ctx, newUser)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username or password is empty")
	}
	user, err := s.userRepository.FindUserByUsername(ctx, username)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	hashPassword, err := security.HashPassword(password)
	if err != nil {
		return "", err
	}
	if !security.CheckPasswordHash(password, hashPassword) {
		return "", ErrInvalidCredentials
	}
	claims := security.CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gophermart",
			Subject:   user.Username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(s.secretKey))
	return signedString, err
}

func (s *AuthService) NewJWTVerifyMiddleware() func(http.Handler) http.Handler {
	return security.JwtVerify([]byte(s.secretKey))
}
