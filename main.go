package main

import (
	"context"
	"log"
	"net/http"

	//mongo

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//jwt
)

var (
	access_secret_key  = []byte("putin")
	refresh_secret_key = []byte("biden")
)

type UserIDRequest struct {
	UserID string `json:"user_id"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	go start_server()
	connect_mongo()
	select {}
}

func start_server() {
	http.HandleFunc("/get_token", get_token)
	http.HandleFunc("/refresh_token", refresh_token)
	err := http.ListenAndServe(":5656", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func connect_mongo() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://overcoder9:gm244922ssFqaLRR@default.pkdymz1.mongodb.net/?retryWrites=true&w=majority&appName=Default"))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

}
