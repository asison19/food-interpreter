package main

import (
	"encoding/json"
	"food-interpreter/lexer"
	"log"
	"net/http"
	"os"

	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//type LexerPost struct {
//	Diary string `json:"diary,string,omitempty"`
//}

func interpretHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p struct {
			Diary string `json:"diary"`
		}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		l := lexer.LexString(p.Diary)

		tokenBytes, err2 := json.Marshal(l.Tokens)
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(tokenBytes)
	})
}

// r type should be the diary?
func Interpret(r *http.Request) {
	var p struct {
		Diary string `json:"diary"`
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l := lexer.LexString(p.Diary)

	tokenBytes, err2 := json.Marshal(l.Tokens)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return tokenBytes
}

func main() {

	image_version := os.Getenv("IMAGE_VERSION")
	log.Printf("Running IMAGE_VERSION: %s", image_version)

	log.Printf("grpc-ping: starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPingServiceServer(grpcServer, &pingService{})
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
