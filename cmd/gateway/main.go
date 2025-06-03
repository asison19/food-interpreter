package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"

	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"


	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
)

//type LexerPost struct {
//	Diary string `json:"diary,string,omitempty"`
//}

func NewConn(host string, insecure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// Note: On the Windows platform, use of x509.SystemCertPool() requires
		// Go version 1.18 or higher.
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	return grpc.Dial(host, opts...)
}

func pingRequest(conn *grpc.ClientConn, p *pb.Request) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := pb.NewPingServiceClient(conn)
	return client.Send(ctx, p)
}

func enqueueDiaryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			Data: []byte(p.Diary),
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

func main() {

	image_version := os.Getenv("IMAGE_VERSION")
	log.Printf("Running IMAGE_VERSION: %s", image_version)

	mux := http.NewServeMux()
	edh := enqueueDiaryHandler()

	mux.Handle("/enqueue-diary", edh)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
