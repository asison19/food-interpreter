package main

import (
	"encoding/json"
	"food-interpreter/lexer"
	"net/http"
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

	http.ListenAndServe(":8080", mux)
}
