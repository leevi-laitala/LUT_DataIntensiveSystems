package main

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"bufio"
	"context"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	mongoUri1 string
	mongoUri2 string
	mongoUri3 string

	servers map[string]*server

	currentSelectedServer *server

	defaultDatabaseName string = "test"
)

// Container for server name and client
type server struct {
	name   string
	client *mongo.Client
}

// Get provided mongo URI env files
// Can be set with ".env" file
func getMongoEnv() error {
	_ = godotenv.Load()

	mongoUri1 = os.Getenv("MONGO_URI_1")
	mongoUri2 = os.Getenv("MONGO_URI_2")
	mongoUri3 = os.Getenv("MONGO_URI_3")

	if mongoUri1 == "" || mongoUri2 == "" || mongoUri3 == "" {
		return fmt.Errorf("env variables MONGO_URI_1, MONGO_URI_2 and MONGO_URI_3 must be set")
	}

	return nil
}

// Return pointer to connected server client
func connectToServer(name string, uri string) (*server, error) {
	api := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(api)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	s := server{
		name:   name,
		client: client,
	}

	return &s, nil
}

// Make connections to the three servers and store them in "servers" map
func initMongoConnections() error {
	err := getMongoEnv()
	if err != nil {
		log.Fatalf("%v", err)
	}

	servers = make(map[string]*server, 3)

	server1, err := connectToServer("server1", mongoUri1)
	if err != nil {
		return fmt.Errorf("failed to connect to server 1: %v", err)
	}
	servers["server1"] = server1

	server2, err := connectToServer("server2", mongoUri2)
	if err != nil {
		return fmt.Errorf("failed to connect to server 2: %v", err)
	}
	servers["server2"] = server2

	server3, err := connectToServer("server3", mongoUri3)
	if err != nil {
		return fmt.Errorf("failed to connect to server 3: %v", err)
	}
	servers["server3"] = server3

	return nil
}

// Close server client connections as a cleanup
func closeMongoConnections() error {
	var err error

	if err = servers["server1"].client.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("failed to connect to server 1: %v", err)
	}

	if err = servers["server2"].client.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("failed to connect to server 2: %v", err)
	}

	if err = servers["server3"].client.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("failed to connect to server 3: %v", err)
	}

	return nil
}

// List collections from current server
func listCollections() {
	if currentSelectedServer == nil {
		fmt.Println("Must select server first")
		return
	}

	db := currentSelectedServer.client.Database(defaultDatabaseName)
	collections, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		fmt.Printf("Could not list collections: %v\n", err)
		return
	}

	fmt.Println("Collections:")
	for _, collection := range collections {
		fmt.Printf("- %s\n", collection)
	}
}

// List contents of collection from current server
func listCollectionContents(collectionName string, id int) {
	if currentSelectedServer == nil {
		fmt.Println("Must select server first")
		return
	}

	db := currentSelectedServer.client.Database(defaultDatabaseName)
	coll := db.Collection(collectionName)

	if coll == nil {
		fmt.Printf("No such collection")
		return
	}

	filter := bson.D{}

	if id > 0 {
		filter = bson.D{{"_id", id}}
	}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		fmt.Printf("Could not list collections: %v", err)
		return
	}

	var results []bson.M
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		fmt.Printf("Could not list collections: %v", err)
		return
	}


	// Pretty print if only single record is returned
	if len(results) == 1 {
		var out bytes.Buffer
		d := []byte(fmt.Sprintf("%s", results[0]))
		err = json.Indent(&out, d, "", "\t")
		if err != nil {
			fmt.Printf("JSON indent error: %v\n", err)
			return
		}
		fmt.Printf("%s\n", out.String())

		return
	}

	// "Ugly" print multiple records
	for _, result := range results {
		fmt.Println(result)
	}
}

func readJsonFromStdin() ([]byte, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	fmt.Println("Enter JSON. Finish with empty line:")

	for {
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		}

		lines = append(lines, line)
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	return []byte(strings.Join(lines, "\n")), nil
}

// Unmarshal bytearray JSON data to a map
func parseBytearrayToJson(data []byte) (map[string]any, error) {
	var doc map[string]any
	err := json.Unmarshal(data, &doc)
	return doc, err
}

// CLI command for data insertion
func cmdInsertData(collectionName string) {
	if currentSelectedServer == nil {
		fmt.Println("Must select server first")
		return
	}

	if collectionName == "" {
		fmt.Println("Must provide collection")
		fmt.Println("Usage 'insert <collection>'")
		return
	}

	jsonBytes, err := readJsonFromStdin()
	if err != nil {
		fmt.Printf("Failed to read JSON: %v\n", err)
		return
	}

	//doc, err := parseBytearrayToJson(jsonBytes)
	//if err != nil {
	//	fmt.Printf("Invalid JSON: %v\n", err)
	//	return
	//}

	db := currentSelectedServer.client.Database(defaultDatabaseName)
	coll := db.Collection(collectionName)

	if coll == nil {
		fmt.Printf("No such collection")
		return
	}

	doc := make(map[string]interface{})
	dec := json.NewDecoder(bytes.NewReader(jsonBytes))
	dec.UseNumber()

	err = dec.Decode(&doc)
	if err != nil {
		fmt.Printf("Invalid JSON: %v\n", err)
	}

	// Hack to get $numberInt IDs
	if num, ok := doc["_id"].(json.Number); ok {
		i, _ := num.Int64()
		doc["_id"] = int32(i)
	}

	res, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("Could not list collections: %v", err)
	}

	fmt.Printf("Inserted with ID: %v\n", res.InsertedID)

	return
}

// CLI command for data deletion
func cmdDeleteData(collectionName string, idStr string) {
	if currentSelectedServer == nil {
		fmt.Println("Must select server first")
		return
	}

	if collectionName == "" || idStr == "" {
		fmt.Println("Must provide collection and ID")
		fmt.Println("Usage 'delete <collection> <id>'")
		return
	}

	db := currentSelectedServer.client.Database(defaultDatabaseName)
	coll := db.Collection(collectionName)

	if coll == nil {
		fmt.Printf("No such collection")
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Must provide integer ID")
		return
	}

	filter := bson.D{{"_id", idInt}}

	res, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Printf("Could not delete entry: %v", err)
	}

	fmt.Printf("Deleted %d entries\n", res.DeletedCount)

	return
}
