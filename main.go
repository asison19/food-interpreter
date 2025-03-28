package main

import (
	"encoding/json"
	"food-interpreter/lexer"
	"log"
	"net/http"
	"os"
)

//type LexerPost struct {
//	Diary string `json:"diary,string,omitempty"`
//}

func lexerHandler() http.Handler {
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

func main() {
	mux := http.NewServeMux()
	lh := lexerHandler()

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
