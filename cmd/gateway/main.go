package main

import (
	"context"
	"encoding/json"
	"flag"
	pb "food-interpreter/interpreter/proto"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"

	//"crypto/tls"
	//"crypto/x509"

	//"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // TODO secure
)

var (
	addr = flag.String("addr", os.Getenv("INTERPRETER_CLOUD_RUN_URI"), "The gRPC server address to connect to")
)

//type LexerPost struct {
//	Diary string `json:"diary,string,omitempty"`
//}

func enqueueDiaryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode the diary
		// TODO turn to function
		var p struct {
			Diary string `json:"diary"`
		}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create a client
		ctx := context.Background()
		gcp_project_id := os.Getenv("GCP_PROJECT_ID")
		client, err := pubsub.NewClient(ctx, gcp_project_id)
		if err != nil {
			log.Printf("Error using GCP Project ID: %s; Error: %s", gcp_project_id, err)
			return
		}

		// Define the topic
		topic_id := os.Getenv("TOPIC_ID")
		topic_id_split := strings.Split(topic_id, "/")
		topic := client.Topic(topic_id_split[len(topic_id_split)-1])

		// Publish a message
		result := topic.Publish(ctx, &pubsub.Message{
			Data: []byte("{\"diary\": \"" + p.Diary + "\"}"),
		})

		// Get the message ID
		id, err := result.Get(ctx)
		if err != nil {
			log.Printf("Error getting message ID: %s; %s", id, err)
			return
		}
		io.WriteString(w, "Published message ID "+id+".\n")
	})
}

func interpretHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flag.Parse()

		// Decode the diary
		// TODO turn to function
		var p struct {
			Diary string `json:"diary"`
		}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set up a connection to the server.
		conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewInterpreterServerClient(conn)

		log.Println("Address of the gRPC server: " + *addr)
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		reply, err := c.Interpret(ctx, &pb.DiaryRequest{Diary: p.Diary})
		if err != nil {
			log.Fatalf("Could not interpret: %v", err)
		}
		log.Printf("Diary Output: %s", reply.GetTokens())
	})
}

func main() {
	image_version := os.Getenv("IMAGE_VERSION")
	log.Printf("Gateway running IMAGE_VERSION: %s", image_version)

	mux := http.NewServeMux()
	edh := enqueueDiaryHandler()
	ih := interpretHandler()

	mux.Handle("/enqueue-diary", edh)
	mux.Handle("/interpret", ih)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
