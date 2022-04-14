package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {
	mux := httprouter.New()

	ytApiKey := os.Getenv("YOUTUBE_API_KEY")

	if ytApiKey == "" {
		log.Fatal("youtube API key not provided")
	}

	mux.GET("/youtube/channel/stats", getChannelStats(ytApiKey))

	return mux
}

// func getChannelStats() httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 		w.Write([]byte("response!"))
// 	}
// }

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	
	srv := &http.Server{
		Addr: ":10101",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("service interrupt received")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}

		log.Println("shutdown complete")

		close(idleConnsClosed)
	}()

	log.Printf("Starting server on port 10101")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("fatal http server failed to start: %v", err)
		}
	}

	<-idleConnsClosed // this channel prevents program from exiting until this channel is closed
	log.Println("Service Stop")
}