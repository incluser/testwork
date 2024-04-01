package main

import (
	"context"
	"encoding/json"
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

	refreshToken, err := generateToken(userID, refresh_secret_key, time.Hour*24*7)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedToken, err := generateSHA256Hash(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	coll := dbClient.Database("testwork").Collection("users")

	_, err = coll.InsertOne(context.TODO(), bson.M{"user_id": userID, "refresh_hash": hashedToken})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenPair := TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)

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

	refreshTokenClaims, err := verifyToken(refreshTokenRequest.RefreshToken, refresh_secret_key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, ok := refreshTokenClaims["user_id"].(string)
	if !ok {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := generateToken(userID, refresh_secret_key, time.Hour*24*7)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedToken, err := generateSHA256Hash(newRefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	collection := dbClient.Database("testwork").Collection("users")

	_, err = collection.UpdateOne(context.TODO(), bson.M{"user_id": userID}, bson.M{"$set": bson.M{"refresh_hash": hashedToken}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, err := generateToken(userID, access_secret_key, time.Minute*15)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenPair := TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}
