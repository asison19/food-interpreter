package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	"food-interpreter/interpreter"
	pb "food-interpreter/interpreter/proto"
	"food-interpreter/lexer"

	"google.golang.org/grpc"
	//"crypto/tls"
	//"crypto/x509"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
)

var (
	grpc_port = flag.String("port", os.Getenv("PORT"), "The port to listen to for gRPC requests.")
)

//type LexerPost struct {
//	Diary string `json:"diary,string,omitempty"`
//}

type server struct {
	pb.UnimplementedInterpreterServerServer // TODO name
}

func (s *server) Interpret(ctx context.Context, in *pb.DiaryRequest) (*pb.DiaryReply, error) {
	p := interpreter.Interpret(in.GetDiary())

	return &pb.DiaryReply{Tokens: lexer.GetTokensAsString(p.Tokens)}, nil
}

func main() {
	flag.Parse()

	image_version := os.Getenv("IMAGE_VERSION")
	log.Printf("Interpreter gRPC service running IMAGE_VERSION: %s", image_version)
	log.Printf("Interpreter gRPC service starting on port %s", *grpc_port)

	// grpc
	lis, err := net.Listen("tcp", ":"+*grpc_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInterpreterServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
