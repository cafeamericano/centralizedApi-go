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

    "github.com/gorilla/mux" // go get -u github.com/gorilla/mux

    "github.com/gorilla/handlers" //go get github.com/gorilla/handlers

    "go.mongodb.org/mongo-driver/bson" // go get go.mongodb.org/mongo-driver
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	
)

// Define custom types

type ObjectID string

type KeywordInCollection struct {
	ID      primitive.ObjectID  `json:"_id" bson:"_id,omitempty"`
    Name    string              `json:"name"`
    Type    string              `json:"type"`
}

type KeywordFromRequest struct {
	ID    	string  			`json:"_id"`
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
}

// Define route-specific functions

// LANDING PAGE

func sendHome(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, "<h2>Centralized API - Go</h2><p>Version 1.0</p><p>The server is listening for requests.</p>")
    return
}

// GET ALL KEYWORDS

func getAllKeywords(w http.ResponseWriter, req *http.Request) {
	formatResponseHeader(w)

    findOptions := options.Find()
    // findOptions.SetLimit(100)

    cursor, err := (keywordsCollection).Find(context.TODO(), bson.D{{}}, findOptions)
    logErrorIfOccurs(err)

    var results []*KeywordInCollection
    for cursor.Next(context.TODO()) { // Iterate over cursor and decode each document
        var item KeywordInCollection
        err := cursor.Decode(&item)
        logErrorIfOccurs(err)
        results = append(results, &item)
    }
    cursor.Close(context.TODO())
    
    json.NewEncoder(w).Encode(results)
}

// READ KEYWORD DATA

func getKeyword(w http.ResponseWriter, req *http.Request) {
	formatResponseHeader(w)

	queryStringId := req.URL.Query().Get("id")
	findOptions := options.Find()
	// findOptions.SetLimit(100)
	
	id, _ := primitive.ObjectIDFromHex(queryStringId)

    cursor, err := (keywordsCollection).Find(context.TODO(), bson.M{"_id": id}, findOptions)
    logErrorIfOccurs(err)
	
    var results []*KeywordInCollection
    for cursor.Next(context.TODO()) { // Iterate over cursor and decode each document
        var item KeywordInCollection
        err := cursor.Decode(&item)
        logErrorIfOccurs(err)
        results = append(results, &item)
    }
	cursor.Close(context.TODO())
	
    json.NewEncoder(w).Encode(results)
	
}

// CREATE A NEW KEYWORD

func addKeyword(w http.ResponseWriter, req *http.Request) {

    formatResponseHeader(w)

    decoder := json.NewDecoder(req.Body)
    var item KeywordInCollection
    err := decoder.Decode(&item)
    if err != nil {
        panic(err)
    }
    
    result, err := keywordsCollection.InsertOne(context.TODO(), item)
	fmt.Printf("Added document \n", result)
}

// UPDATE A KEYWORD

func updateKeyword(w http.ResponseWriter, req *http.Request) {
    formatResponseHeader(w)

	decoder := json.NewDecoder(req.Body)
    var item KeywordFromRequest
    err := decoder.Decode(&item)
    if err != nil {
        panic(err)
	}

	ctx := context.Background()
	id, _ := primitive.ObjectIDFromHex(item.ID)
	result, err := keywordsCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"name", &item.Name},
				{"type", &item.Type},
			}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v document(s)\n", result.ModifiedCount)
}

// DELTETE A KEYWORD

func deleteKeyword(w http.ResponseWriter, req *http.Request) {

    formatResponseHeader(w)

	decoder := json.NewDecoder(req.Body)
    var item KeywordFromRequest
    err := decoder.Decode(&item)
    logErrorIfOccurs(err)

	ctx := context.Background()
	id, _ := primitive.ObjectIDFromHex(item.ID)
	result, err := keywordsCollection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)
	logErrorIfOccurs(err)

	fmt.Printf("Deleted %v document(s)\n", result.DeletedCount)

}

// Define main function

func main() {

    router := mux.NewRouter()
	
	router.HandleFunc("/", sendHome).Methods("GET")

	router.HandleFunc("/KeywordFactory/api/allkeywords", getAllKeywords).Methods("GET")

	router.HandleFunc("/KeywordFactory/api/keyword", getKeyword).Methods("GET")
    router.HandleFunc("/KeywordFactory/api/keyword", addKeyword).Methods("POST")
    router.HandleFunc("/KeywordFactory/api/keyword", updateKeyword).Methods("PUT")
    router.HandleFunc("/KeywordFactory/api/keyword", deleteKeyword).Methods("DELETE")
    
    http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))

}