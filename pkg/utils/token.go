package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// * TOKEN CREATE
func CreateToken(userID int64, username string, duration time.Duration, secretKey string) (string, error) {
	payload := &Payload{
		ID:        uuid.New(),
		UserID:    userID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         payload.ID,
		"user_id":    payload.UserID,
		"username":   payload.Username,
		"issued_at":  payload.IssuedAt.Unix(),
		"expired_at": payload.ExpiredAt.Unix(),
	})

	return jwtToken.SignedString([]byte(secretKey))
}

// * TOKEN DOĞRULAMA
func VerifyToken(token string, secretKey string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secretKey), nil
	}

	jwtToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	payload := &Payload{
		ID:       uuid.MustParse(claims["id"].(string)),
		UserID:   int64(claims["user_id"].(float64)),
		Username: claims["username"].(string),
	}

	return payload, nil
}
