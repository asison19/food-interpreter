package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"food-interpreter/interpreter"
	pb "food-interpreter/interpreter/proto"
	"food-interpreter/lexer"
	"log"
	"net"
	"net/http"
	"os"

	//"crypto/tls"
	//"crypto/x509"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
)

var (
	pb_port = flag.Int("port", 50051, "The server port")
)

//type LexerPost struct {
//	Diary string `json:"diary,string,omitempty"`
//}

type server struct {
	pb.UnimplementedInterpreterServerServer // TODO name
}

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

		parser := interpreter.Interpret(p.Diary)

		tokenBytes, err2 := json.Marshal(parser.Tokens)
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(tokenBytes)
	})
}

func (s *server) Interpret(ctx context.Context, in *pb.DiaryRequest) (*pb.DiaryReply, error) {
	p := interpreter.Interpret(in.GetDiary())

	return &pb.DiaryReply{Tokens: lexer.GetTokensAsString(p.Tokens)}, nil
}

func main() {

	image_version := os.Getenv("IMAGE_VERSION")
	log.Printf("Running IMAGE_VERSION: %s", image_version)

	mux := http.NewServeMux()
	ih := interpretHandler()

	mux.Handle("/interpret", ih)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	go func() {
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatal(err)
		}
	}()

	// grpc
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *pb_port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &server{})
	pb.RegisterInterpreterServerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
