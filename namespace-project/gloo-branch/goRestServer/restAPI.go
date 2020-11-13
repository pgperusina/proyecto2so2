package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var redisClient *redis.Client
var ctx = context.TODO()

type Caso struct {
	Name         string
	Location     string
	Age          int
	InfectedType string
	State        string
}

var counter int = 0

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Covid 19 Go API")
}

func crearCaso(w http.ResponseWriter, r *http.Request) {
	log.Println("Creando nuevo caso")

	var caso Caso

	reqBody, _ := ioutil.ReadAll(r.Body)
	bodyString := string(reqBody)

	err := json.Unmarshal(reqBody, &caso)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Print("Error decoding request body to json")
		//log.Fatal(err)
		fmt.Fprintf(w, "Caso %+v no pudo ser insertado!", bodyString)
		return
	}

	fmt.Fprintf(w, "Caso %+v insertado!", caso)
	//insertarMongoDB(caso)
	//insertarRedis(bodyString)

	publishToRabbitMQ(bodyString)
}

func insertarMongoDB(document Caso) {
	collection := mongoClient.Database("covid19").Collection("casos")
	insertResult, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Println("Error posting to mongoDB: ", err)
		//log.Fatal(err)
		return
	}
	log.Println("Caso insertado en mongoDB:", insertResult.InsertedID)
}

func insertarRedis(document string) {
	counter += 1
	err := redisClient.Set(strconv.Itoa(counter), document, 0).Err()

	if err != nil {
		log.Println("Error posting to Redis")
		log.Println(err)
		return
	}
	log.Println("Caso insertado en Redis:", document)
}

func failOnError(err error, message string) {
	if err != nil {
		log.Println("%s: %s", message, err)
	}
}

func publishToRabbitMQ(caso string) {
	conn, err := amqp.Dial("amqp://guest:guest@34.72.226.148:5672/")
	failOnError(err, "Error conectando con RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Error abriendo el canal de RabitMQ")
	defer conn.Close()

	q, err := ch.QueueDeclare(
		"covid19",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Error declarando la cola en RabbitMQ.")

	//body := fmt.Sprintf("%s", caso)
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(caso),
		})
	log.Printf("Enviando... %+v", caso)
	failOnError(err, "Error al enviar el mensaje")
}

func handleRequests() {
	r := mux.NewRouter()

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/caso", crearCaso).Methods("POST")

	log.Fatal(http.ListenAndServe(":80", r))
}

func connectMongoDB() *mongo.Client {

	clientOptions := options.Client().ApplyURI("mongodb://sopes1:sopes1proyecto2@34.67.186.172:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error connecting to mongoDB: ", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Connected to MongoDB!")
	return client
}

func connectRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "34.66.203.76:6379",
		Password: "sopes1proyecto2",
		DB:       0, // use default DB
	})

	pong, err := redisClient.Ping().Result()
	log.Println(pong, err)

	log.Println("Connected to Redis!")
	return redisClient
}

func main() {
	//mongoClient = connectMongoDB()
	//redisClient = connectRedis()
	handleRequests()
}
