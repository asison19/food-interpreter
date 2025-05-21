package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
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
		topic := client.Topic(topic_id)

		// Publish a message
		result := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(p.Diary),
		})

		// TODO 2025/05/20 21:32:34 Error getting message ID: ; rpc error: code = InvalidArgument desc = Invalid resource name given (name=projects/food-interpreter/topics/). Refer to https://cloud.google.com/pubsub/docs/pubsub-basics#resource_names for more information.
		// Get the message ID
		id, err := result.Get(ctx)
		if err != nil {
			log.Printf("Error getting message ID: %s; %s", id, err)
			return
		}
		io.WriteString(w, "Published message ID "+id+".\n")

		//l := lexer.LexString(p.Diary)

		//tokenBytes, err2 := json.Marshal(l.Tokens)
		//if err2 != nil {
		//	http.Error(w, err.Error(), http.StatusBadRequest)
		//	return
		//}
		//w.Write(tokenBytes)
	})
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
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
