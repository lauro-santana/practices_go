package main

import (
	"log"
	"net/http"

	"movie.com/rating/internal/controller/rating"
	httphandler "movie.com/rating/internal/handler/http"
	"movie.com/rating/internal/repository/memory"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
