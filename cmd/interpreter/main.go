package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	pb "food-interpreter/interpreter/proto"
	"food-interpreter/lexer"
	//"crypto/tls"
	//"crypto/x509"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
)

var (
	port = flag.String("port", os.Getenv("PORT"), "The port to listen to for HTTP requests.")
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
		log.Println("interpretHandler")
		diary, err := normalizeDiaryJSON(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		l := lexer.LexString(diary)

		tokenBytes, err := json.Marshal(l.Tokens)
		if err != nil {
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

		diaryJSON := string(m.Message.Data)

		diary, err := normalizeDiaryJSON(strings.NewReader(diaryJSON))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		l := lexer.LexString(diary)
		tokenBytes, err := json.Marshal(l.Tokens)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(tokenBytes)
		log.Printf("Tokens: " + string(tokenBytes))
	})
}

func normalizeDiaryJSON(body io.Reader) (string, error) {
	var p struct {
		Diary string `json:"diary"`
	}
	err := json.NewDecoder(body).Decode(&p)
	if err != nil {
		return "", err
	}
	log.Printf("Diary: " + p.Diary)
	return p.Diary, nil
}

func main() {
	flag.Parse()

	image_version := os.Getenv("IMAGE_VERSION")
	log.Printf("Interpreter running IMAGE_VERSION: %s", image_version)

	mux := http.NewServeMux()
	ih := interpretHandler()
	psih := pubsubInterpretHandler()

	mux.Handle("/interpret", ih)
	mux.Handle("/pubsub-interpret", psih)

	// http
	log.Printf("Listening on port %s", *port)

	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatal(err)
	}
}
