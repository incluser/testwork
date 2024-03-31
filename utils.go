package main

import (
	"fmt"
	"time"

	//mongo

	//jwt
	"github.com/dgrijalva/jwt-go"
)

func verifyToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем тип подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Проверяем, действителен ли токен
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Получаем утверждения токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func generateToken(userID string, secretKey []byte, expirationTime time.Duration) (string, error) {
	// Создаем новый JWT токен
	token := jwt.New(jwt.SigningMethodHS256)
	// Устанавливаем поля токена
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(expirationTime).Unix() // Устанавливаем время истечения токена
	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
