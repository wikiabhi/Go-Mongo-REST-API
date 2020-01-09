package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Student struct {
	Name 	string	`json: "name""`
	Age 	int		`json: "age"`
	City	string	`json: "city"`
}

const (
	DBNAME = "go_mongo"
	URI = "mongodb://127.0.0.1:27017/"
)

func getStudents(w http.ResponseWriter, r *http.Request){
	ctx := context.Background()
	//set client options
	clientOptions := options.Client().ApplyURI(URI)

	//connect to mongodb
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return
	}

	fmt.Println("DB Connected")

	db := client.Database(DBNAME)
	collection := db.Collection("student")

	student := Student{}
	studentList := []Student{}
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		fmt.Println(err)
		return
	}

	for cursor.Next(ctx) {
		cursor.Decode(&student)
		studentList = append(studentList, student)
	}

	json.NewEncoder(w).Encode(&studentList)

	//close connection
	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB disconnected")
}

func createStudent(w http.ResponseWriter, r *http.Request){
	ctx := context.Background()
	//set client options
	clientOptions := options.Client().ApplyURI(URI)

	//connect to mongodb
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return
	}

	fmt.Println("DB Connected")

	db := client.Database(DBNAME)
	collection := db.Collection("student")


	var student Student
	_ = json.NewDecoder(r.Body).Decode(&student)


	//new dummy data
	//student := Student{}
	//student.Name = "Abhishek"
	//student.Age = 24
	//student.City = "Pune"

	result, err := collection.InsertOne(ctx, student)

	//objectId := result.InsertedID.(primitive.ObjectID)

	json.NewEncoder(w).Encode(result)

	//close connection
	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB disconnected")
}

func updateStudent(w http.ResponseWriter, r *http.Request){
	ctx := context.Background()
	//set client options
	clientOptions := options.Client().ApplyURI(URI)

	//connect to mongodb
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return
	}

	fmt.Println("DB Connected")
	db := client.Database(DBNAME)
	collection := db.Collection("student")



	params := mux.Vars(r)

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		fmt.Println(err)
		return
	}

	var student Student
	_ = json.NewDecoder(r.Body).Decode(&student)

	resultUpdate, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			bson.M{
				"$set" : bson.M {
					"name" 	: student.Name,
					"age"	: student.Age,
					"city"	: student.City,
				},
			},
	)

	json.NewEncoder(w).Encode(resultUpdate.ModifiedCount)

	//close connection
	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB disconnected")
}

func deleteStudent(w http.ResponseWriter, r *http.Request)  {
	ctx := context.Background()
	//set client options
	clientOptions := options.Client().ApplyURI(URI)

	//connect to mongodb
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return
	}

	fmt.Println("DB Connected")
	db := client.Database(DBNAME)
	collection := db.Collection("student")


	params := mux.Vars(r)

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})

	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(result.DeletedCount)

	//close connection
	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB disconnected")
}

func main(){

	router := mux.NewRouter()

	//Routes
	router.HandleFunc("/api/students", getStudents).Methods("GET")
	router.HandleFunc("/api/students", createStudent).Methods("POST")
	router.HandleFunc("/api/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/api/students/{id}", deleteStudent).Methods("DELETE")


	fmt.Println("Server Starting")
	http.ListenAndServe(":8000", router)

}