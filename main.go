package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//mongo
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	coll := client.Database("sample_mflix").Collection("movies")
	title := "Back to the Future"

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}

func get_token(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет!")
}

func refresh_token(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет refresh")
}
