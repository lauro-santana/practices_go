package main

import (
	"log"
	"net/http"

	"movie.com/movie/internal/controller/movie"
	metadatagateway "movie.com/movie/internal/gateway/metadata/http"
	ratinggateway "movie.com/movie/internal/gateway/rating/http"
	httphandler "movie.com/movie/internal/handler/http"
)

func main() {
	log.Println("Starting the movie service")
	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
