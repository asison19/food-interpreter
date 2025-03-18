package main

import (
	"context"
	"encoding/json"
	"food-interpreter/lexer"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
)

//type LexerPost struct {
//	Diary string `json:"diary,string,omitempty"`
//}

func lexerHandler(logger *log.Logger) http.Handler {
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
		log.Printf(string(tokenBytes))
		logger.Println("lexerHandler: " + string(tokenBytes))
	})
}

func setupLogging() *log.Logger {
	ctx := context.Background()

	client, err := logging.NewClient(ctx, "food-interpreter")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	logName := "lexer-log"

	logger := client.Logger(logName).StandardLogger(logging.Info)

	logger.Println("Lexer logging set up.")
	return logger
}

func main() {
	logger := setupLogging() // TODO don't run when not running on GCP.

	mux := http.NewServeMux()
	lh := lexerHandler(logger)

	mux.Handle("/lexer", lh)

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
