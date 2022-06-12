package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlmjohnson/gateway"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)




func main() {
	
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	return
	// }
	
	listener := gateway.ListenAndServe
	
	portStr := ""
	if os.Getenv("ENV") == "development" {
		fmt.Println("local server")
		portStr = fmt.Sprintf(":%d", 1000)
		listener = http.ListenAndServe
		http.Handle("/", http.FileServer(http.Dir("./public")))
	}
	
	http.Handle("/api/products", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_DB_URI")))
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
		coll := client.Database("digital-store").Collection("products")
	
		
		cursor, err := coll.Find(context.TODO(), bson.M{})
		if err != nil {
			panic(err)
		}
		
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			log.Fatal(err)
		}
		
		var allData []interface{}
		
		for _, result := range results {
			allData = append(allData, result)
		}
		
		json.NewEncoder(writer).Encode(&allData)
		
	}))
	
	log.Fatal(listener(portStr, nil))
}

