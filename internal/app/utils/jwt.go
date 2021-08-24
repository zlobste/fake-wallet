package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	TokenExpiration = time.Hour * 24

	UserIDTag = "userID"
)

var (
	ErrInvalidJWT = errors.New("Invalid JWT")
)

type Auth interface {
	CreateJWT(userId int64) (string, error)
	TokenValid(r *http.Request) error
	ExtractTokenMetadata(r *http.Request) (*AccessDetails, error)
}

type auth struct {
	secret []byte
}

func NewAuth(secret string) Auth {
	return &auth{
		secret: []byte(secret),
	}
}

func (t *auth) CreateJWT(userId int64) (string, error) {
	claims := jwt.MapClaims{
		UserIDTag:   userId,
		"IssuedAt":  time.Now().Unix(),
		"ExpiresAt": time.Now().Add(TokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString(t.secret)
}

func (t *auth) TokenValid(r *http.Request) error {
	token, err := t.verifyToken(r)
	if err != nil {
		return errors.Wrap(err, "failed to verify token")
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return ErrInvalidJWT
	}
	return nil
}

func (t *auth) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := t.verifyToken(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims[UserIDTag]), 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse user id")
		}
		return &AccessDetails{
			UserID: userID,
		}, nil
	}
	return nil, err
}

func (t *auth) verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := t.extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.secret, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse JWT")
	}
	return token, nil
}

func (t *auth) extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

type AccessDetails struct {
	UserID int64
}
