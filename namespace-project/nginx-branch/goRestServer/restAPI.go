package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type Caso struct {
	Name         string
	Location     string
	Age          int
	InfectedType string
	State        string
}

const (
	address = "python-service-grpc:50051"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Covid 19 Go API -- Nginx ingress")

}
func enviarGrcp(caso string) {

	log.Println("Enviando caso: " + caso)
	// Envio de mensaje
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("Error conectando via grpc: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.SayHello(ctx, &pb.HelloRequest{Name: caso})
	if err != nil {
		log.Printf("Error al enviar el mensaje via grcp: %v", err)
	}
	log.Printf("Respuesta grcp: %s", res.GetMessage())
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Covid 19 Go API -- Nginx ingress")
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
		log.Println("Caso %+v no pudo ser insertado!", bodyString)
		fmt.Fprintf(w, "Caso no pudo ser insertado!", bodyString)
		return
	}
	enviarGrcp(bodyString)
	fmt.Fprintf(w, "Caso %+v insertado via grcp!", caso)
}

func handleRequests() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/caso", crearCaso).Methods("POST")
	log.Println("Servidor levantado en el puerto 80")
	log.Fatal(http.ListenAndServe(":80", r))

}

func main() {

	handleRequests()
}
