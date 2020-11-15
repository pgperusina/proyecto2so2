package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
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

func setupTracer() (trace.Tracer, *jaeger.Exporter, error) {
	// Register installs a new global tracer instance.
	tracer := sdk.Register()

	// Construct and register an export pipeline using the Jaeger
	// exporter and a span processor.
	exporter, err := jaeger.NewExporter(
		jaeger.Options{
			AgentEndpoint: "jaeger-agent.observability.svc.cluster.local:6831",
		},
	)
	if err != nil {
		return nil, nil, err
	}

	// A simple span processor calls through to the exporter
	// without buffering.
	ssp := sdk.NewSimpleSpanProcessor(exporter)
	sdk.RegisterSpanProcessor(ssp)

	// Use sdk.AlwaysSample sampler to send all spans.
	sdk.ApplyConfig(
		sdk.Config{
			DefaultSampler: sdk.AlwaysSample(),
		},
	)

	return tracer, exporter, nil
}

func homePage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tracer := trace.GlobalTracer()

	ctx, trace := tracer.Start(ctx, "go-ws-homepage")
	fmt.Fprintf(w, "Covid 19 Go API -- Nginx ingress")

	trace.End()
}
func enviarGrcp(caso string) {
	ctx := context.Background()
	tracer := trace.GlobalTracer()

	ctx, trace := tracer.Start(ctx, "go-ws-enviar-caso-via-grcp")

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
	trace.End()
}

func home(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tracer := trace.GlobalTracer()

	ctx, trace := tracer.Start(ctx, "go-ws-homepage")
	fmt.Fprintf(w, "Covid 19 Go API -- Nginx ingress")

	trace.End()
}

func crearCaso(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tracer := trace.GlobalTracer()

	ctx, trace := tracer.Start(ctx, "go-ws-recibe-caso-grcp")
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
	trace.End()
}

func handleRequests() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/caso", crearCaso).Methods("POST")
	log.Println("Servidor levantado en el puerto 80")
	log.Fatal(http.ListenAndServe(":80", r))

}

func main() {
	tracer, exporter, err := setupTracer()
	if err != nil {
		log.Fatal("Could not initialize tracing: ", err)
	}
	handleRequests()
}
