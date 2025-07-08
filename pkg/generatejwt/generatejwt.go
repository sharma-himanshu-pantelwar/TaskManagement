package generatejwt

import (
	"taskmgmtsystem/internal/core/session"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Uid int
	jwt.StandardClaims
}

func GenerateJWT(uid int, jwtKey string) (string, time.Time, error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, expirationTime, nil
}

func GenerateSession(userId int) (session.Session, error) {
	tokenId := uuid.New()
	expiresAt := time.Now().Add(12 * time.Hour)
	issuedAt := time.Now()

	hashToken, err := bcrypt.GenerateFromPassword([]byte(tokenId.String()), bcrypt.DefaultCost)

	if err != nil {
		return session.Session{}, err
	}

	session := session.Session{
		Id:        tokenId,
		UserId:    userId,
		TokenHash: string(hashToken),
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}
	return session, nil
}

func ValidateJWT(tokenString string, jwtKey string) (*Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return &claims, nil
}
