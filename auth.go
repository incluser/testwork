package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func get_token(w http.ResponseWriter, r *http.Request) {

	var userIDRequest UserIDRequest
	err := json.NewDecoder(r.Body).Decode(&userIDRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := userIDRequest.UserID
	if userID == "" {
		http.Error(w, "missing field user_id", http.StatusBadRequest)
		return
	}

	accessToken, err := generateToken(userID, access_secret_key, time.Minute*15)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken := getRandomString()

	hashedToken, err := generateSHA256Hash(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	coll := dbClient.Database("testwork").Collection("users")

	_, err = coll.InsertOne(context.TODO(), bson.M{"user_id": userID, "refresh_hash": hashedToken, "exp": time.Now().Add(time.Hour * 24 * 7).Unix()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenPair := TokenPair{
		AccessToken:  accessToken,
		RefreshToken: hashedToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)

}

func getRandomString() string {
	rand.Seed(time.Now().UnixNano())

	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	length := 32
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

type User struct {
	UserID     string `bson:"user_id"`
	Expiration int64  `bson:"exp"`
}

func refresh_token(w http.ResponseWriter, r *http.Request) {

	var refreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&refreshTokenRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User

	userCollection := dbClient.Database("testwork").Collection("users")

	findByRefreshHash := bson.M{"refresh_hash": refreshTokenRequest.RefreshToken}

	err = userCollection.FindOne(context.TODO(), findByRefreshHash).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	if time.Now().Unix() > user.Expiration {
		http.Error(w, "token expired", http.StatusUnauthorized)
		return
	}

	refreshToken := getRandomString()

	hashedToken, err := generateSHA256Hash(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	collection := dbClient.Database("testwork").Collection("users")

	_, err = collection.UpdateOne(context.TODO(), bson.M{"user_id": user.UserID}, bson.M{"$set": bson.M{"refresh_hash": hashedToken, "exp": time.Now().Add(time.Hour * 24 * 7).Unix()}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, err := generateToken(user.UserID, access_secret_key, time.Minute*15)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenPair := TokenPair{
		AccessToken:  accessToken,
		RefreshToken: hashedToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}
