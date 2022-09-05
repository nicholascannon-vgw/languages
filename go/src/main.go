package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"techtonic/src/middleware"
	"techtonic/src/req"
	"techtonic/src/res"

	"github.com/gorilla/mux"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type WordCountBody struct {
	Text string `json:"text"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		res.Status(w, 200) // Note: status must come before JSON
		res.JSON(w, MessageResponse{Message: "healthy"})
	})

	router.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		payload := WordCountBody{}
		req.ParseJSON(r, &payload)

		wordCounts := map[string]uint{}

		words := strings.Split(payload.Text, " ")
		for _, word := range words {
			word = strings.ReplaceAll(word, "!", "")
			word = strings.ReplaceAll(word, ",", "")
			word = strings.ReplaceAll(word, "?", "")
			word = strings.ReplaceAll(word, "\n", "")
			word = strings.ReplaceAll(word, "\t", "")
			word = strings.ReplaceAll(word, "\"", "")
			word = strings.ReplaceAll(word, "\\", "")
			word = strings.ReplaceAll(word, "#", "")

			if strings.TrimSpace(word) == "" {
				continue
			}

			wordCounts[word] += 1
		}

		res.Status(w, 200)
		res.JSON(w, wordCounts)
	})

	router.Use(middleware.Cors) // Allow the frontend to call this service
	router.Use(middleware.DisableKeepAlive)

	fmt.Println("Listening on port 8000...")
	err := http.ListenAndServe(":8000", router)
	log.Fatal(err)
}
