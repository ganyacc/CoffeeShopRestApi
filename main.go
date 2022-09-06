package main

import (
	"CoffeShopRestApi/handler"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

var bindadrress = env.String("Bind_Address", false, ":8090", "Bind Address for server")

func main() {
	env.Parse()
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	ph := handler.NewProduct(l)
	//ServeMux is an HTTP request multiplexer. It matches the URL of each incoming request against a list of registered patterns and calls the handler
	//for the pattern that most closely matches the URL.

	sm := mux.NewRouter()
	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareValidateProduct)
	//if we had created routing instance with NewServeMux method we would not able to use Regex with HandlFunc Method as NewServeMux doesnt support
	//wild cards/ regex
	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods("DELETE").Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteByID)
	s := &http.Server{
		Addr:         ":8090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {

		err := s.ListenAndServe()
		if err != nil {

			log.Fatal(err)
		}
	}()

	//unbuffered channel
	SigChan := make(chan os.Signal)
	signal.Notify(SigChan, os.Interrupt)
	signal.Notify(SigChan, os.Kill)

	sig := <-SigChan
	log.Println("Recieved Terminate, shutdown Gracefully", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
