package utils

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/igostfost/avito_backend_trainee/pkg/repository"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
	"os"
	"time"
)

type AuthUtils struct {
	repo repository.Authorization
}

type CustomClaims struct {
	jwt.StandardClaims
	UserId  int  `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
}

func NewAuthService(repo repository.Authorization) *AuthUtils {
	return &AuthUtils{repo: repo}
}

func (u *AuthUtils) CreateUser(user types.User) (int, error) {
	user.Password = genPasswordHash(user.Password)
	return u.repo.CreateUser(user)
}

func (u *AuthUtils) CreateAdmin(user types.User) (int, bool, error) {
	user.Password = genPasswordHash(user.Password)
	return u.repo.CreateAdmin(user)
}

func (u *AuthUtils) GenerateToken(username, password string) (string, error) {
	user, err := u.repo.GetUser(username, genPasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 8).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.IsAdmin,
	})

	return token.SignedString([]byte(os.Getenv("SIGN_TOKEN_KEY")))
}

func genPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))
}

func (u *AuthUtils) ParseToken(accessToken string) (int, bool, error) {
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("SIGN_TOKEN_KEY")), nil
	})

	if err != nil {
		return 0, false, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, false, errors.New("invalid token")
	}

	return claims.UserId, claims.IsAdmin, nil
}
