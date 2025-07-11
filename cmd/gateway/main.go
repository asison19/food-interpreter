package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	pb "food-interpreter/interpreter/proto"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"

	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure" // TODO secure
	grpcMetadata "google.golang.org/grpc/metadata"
)

var (
	addr     = flag.String("addr", os.Getenv("INTERPRETER_GRPC_CLOUD_RUN_URI"), "The gRPC server address to connect to")
	isSecure = flag.Bool("secure", true, "Whether to use secure authenticated gRPC requests.")
)

//type LexerPost struct {
//	Diary string `json:"diary,string,omitempty"`
//}

func removeScheme(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	u.Scheme = ""
	r, _ := regexp.Compile("[^/].*")
	return r.FindString(u.String())
}

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

// NewConn creates a new gRPC connection.
// host should be of the form domain:port, e.g., example.com:443
func NewConn(host string, isSecure bool) (*grpc.ClientConn, error) {
	log.Println("Setting up new gRPC connection to: " + host)
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	if !isSecure {
		fmt.Println(isSecure)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	return grpc.NewClient(host, opts...)
}

// pingRequest sends a new gRPC ping request to the server configured in the connection.
func interpretRequest(conn *grpc.ClientConn, p *pb.DiaryRequest) (*pb.DiaryReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c := pb.NewInterpreterServerClient(conn)
	return c.Interpret(ctx, p)
}

// pingRequestWithAuth mints a new Identity Token for each request.
// This token has a 1 hour expiry and should be reused.
// audience must be the auto-assigned URL of a Cloud Run service or HTTP Cloud Function without port number.
func interpretRequestWithAuth(conn *grpc.ClientConn, p *pb.DiaryRequest, audience string) (*pb.DiaryReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create an identity token.
	// With a global TokenSource tokens would be reused and auto-refreshed at need.
	// A given TokenSource is specific to the audience.
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("idtoken.NewTokenSource: %w", err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("TokenSource.Token: %w", err)
	}

	// Add token to gRPC Request.
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)

	// Send the request.

	c := pb.NewInterpreterServerClient(conn)
	return c.Interpret(ctx, p)
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
		log.Println("Address of the gRPC server: " + *addr)

		// Check if address contains a scheme, if so, remove it.
		// TODO check if this is still necessary.
		host := *addr
		u, err := url.Parse(*addr)
		if err != nil {
			panic(err)
		}
		if u.Scheme == "http" || u.Scheme == "https" {
			host = removeScheme(host)
		}
		// Set up a connection to the server.
		conn, err := NewConn(host, *isSecure)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		reply := &pb.DiaryReply{}
		if !*isSecure {
			fmt.Println(*isSecure)
			reply, err = interpretRequest(conn, &pb.DiaryRequest{Diary: p.Diary})
		} else {
			reply, err = interpretRequestWithAuth(conn, &pb.DiaryRequest{Diary: p.Diary}, *addr)
		}

		if err != nil {
			log.Fatalf("Could not interpret: %v", err)
		}
		log.Printf("Diary Output: %s", reply.GetTokens())
		io.WriteString(w, "Diary Output: "+reply.GetTokens())
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
