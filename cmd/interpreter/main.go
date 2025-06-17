package main

import (
	"context"
	"encoding/json"
	"flag"
	//"fmt"
	"food-interpreter/interpreter"
	pb "food-interpreter/interpreter/proto"
	"food-interpreter/lexer"
	"log"
	//"net"
	"io"
	"net/http"
	"os"
	"strings"
	//"crypto/tls"
	//"crypto/x509"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
)

var (
	grpc_port = flag.Int("port", 50051, "The gRPC server port")
)

//type LexerPost struct {
//	Diary string `json:"diary,string,omitempty"`
//}

// WrappedMessage is the payload of a Pub/Sub event.
//
// For more information about receiving messages from a Pub/Sub event
// see: https://cloud.google.com/pubsub/docs/push#receive_push
type WrappedMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

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

		log.Printf("Diary: " + p.Diary)
		l := lexer.LexString(p.Diary)

		tokenBytes, err2 := json.Marshal(l.Tokens)
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(tokenBytes)
		log.Printf("Tokens: " + string(tokenBytes))
	})
}

func pubsubInterpretHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m WrappedMessage
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Printf("io.ReadAll: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		// byte slice unmarshalling handles base64 decoding.
		if err := json.Unmarshal(body, &m); err != nil {
			log.Printf("json.Unmarshal: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		diary := string(m.Message.Data)
		log.Printf("Diary: %s!", diary)

		var p struct {
			Diary string `json:"diary"`
		}
		err = json.NewDecoder(strings.NewReader(diary)).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Diary: " + p.Diary)

		l := lexer.LexString(p.Diary)
		tokenBytes, err2 := json.Marshal(l.Tokens)
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(tokenBytes)
		log.Printf("Tokens: " + string(tokenBytes))
	})
}

//func normalizeDiary(body *io.ReadCloser) string {
//	var p struct {
//		Diary string `json:"diary"`
//	}
//	err := json.NewDecoder(body).Decode(&p)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	log.Printf("Diary: " + p.Diary)
//	return p.Diary
//}

func (s *server) Interpret(ctx context.Context, in *pb.DiaryRequest) (*pb.DiaryReply, error) {
	p := interpreter.Interpret(in.GetDiary())

	return &pb.DiaryReply{Tokens: lexer.GetTokensAsString(p.Tokens)}, nil
}

func main() {

	image_version := os.Getenv("IMAGE_VERSION")
	log.Printf("Interpreter running IMAGE_VERSION: %s", image_version)

	mux := http.NewServeMux()
	ih := interpretHandler()
	psih := pubsubInterpretHandler()

	mux.Handle("/interpret", ih)
	mux.Handle("/pubsub-interpret", psih)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}

	// grpc
	//flag.Parse()
	//lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpc_port))
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//s := grpc.NewServer()
	////pb.RegisterGreeterServer(s, &server{})
	//pb.RegisterInterpreterServerServer(s, &server{})
	//log.Printf("server listening at %v", lis.Addr())
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
}
