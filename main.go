// Define package

package main

// Define imports

import (
    "encoding/json"
    "net/http"
    "context"
    "log"
    "os"
    "fmt"

    "github.com/gorilla/mux"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Define custom types

type ObjectID string

type Keyword struct {
    ID      string            	`json:"_id"`
    Name    string              `json:"name"`
    Type    string              `json:"type"`
}

// Initialize collection variables

var keywordsCollection *mongo.Collection

// Initial setup

func init() {
    
    var dbUrl = os.Getenv("DB_URL")
	var clientOptions = options.Client().ApplyURI(dbUrl)
    var client, err = mongo.Connect(context.TODO(), clientOptions)
    logErrorIfOccurs(err)

    keywordsCollection = client.Database("AppGalleryLite").Collection("keywords")
}

// Define helper functions

func logErrorIfOccurs(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func formatResponseHeader(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
}

// Define data gathering functions

func getKeyword(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
    var t Keyword
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
	}


	findOptions := options.Find()
    // findOptions.SetLimit(100)
	id, _ := primitive.ObjectIDFromHex(t.ID)

	cursor, err := (keywordsCollection).Find(context.TODO(), bson.M{"_id": id}, findOptions)
	
    var results []*Keyword
    for cursor.Next(context.TODO()) { // Iterate over cursor and decode each document
        var elem Keyword
        err := cursor.Decode(&elem)
        logErrorIfOccurs(err)
        results = append(results, &elem)
    }
	cursor.Close(context.TODO())
	
    
    formatResponseHeader(w)
    json.NewEncoder(w).Encode(results)
}

func getKeywords(w http.ResponseWriter, r *http.Request) {
    findOptions := options.Find()
    // findOptions.SetLimit(100)

    cursor, err := (keywordsCollection).Find(context.TODO(), bson.D{{}}, findOptions)
    logErrorIfOccurs(err)

    var results []*Keyword
    for cursor.Next(context.TODO()) { // Iterate over cursor and decode each document
        var elem Keyword
        err := cursor.Decode(&elem)
        logErrorIfOccurs(err)
        results = append(results, &elem)
    }
    cursor.Close(context.TODO())
    
    formatResponseHeader(w)
    json.NewEncoder(w).Encode(results)
}

    
func addKeyword(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    decoder := json.NewDecoder(req.Body)
    var t Keyword
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
    }
    log.Println(t.Name)
    
    insertResult, err := keywordsCollection.InsertOne(context.TODO(), t)
    fmt.Print(insertResult)
    // json.NewEncoder(w).Encode(newRecord)
}

func updateKeyword(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
    var t Keyword
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
	}

	ctx := context.Background()
	id, _ := primitive.ObjectIDFromHex(t.ID)
	result, err := keywordsCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"name", &t.Name},
				{"type", &t.Type},
			}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
}

func deleteKeyword(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
    var t Keyword
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
	}

	ctx := context.Background()
	id, _ := primitive.ObjectIDFromHex(t.ID)
	result, err := keywordsCollection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v Documents!\n", result)
}

// Define main function

func main() {

    router := mux.NewRouter()
	
	router.HandleFunc("/AppGalleryLite/api/keyword", getKeyword).Methods("POST")
    router.HandleFunc("/AppGalleryLite/api/keywords", getKeywords).Methods("GET")
    router.HandleFunc("/AppGalleryLite/api/keywords", addKeyword).Methods("POST")
    router.HandleFunc("/AppGalleryLite/api/keywords", updateKeyword).Methods("PUT")
    router.HandleFunc("/AppGalleryLite/api/keywords", deleteKeyword).Methods("DELETE")
    
    http.ListenAndServe(":8080", router)

}