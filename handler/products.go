package handler

import (
	"CoffeShopRestApi/data"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Encode", http.StatusInternalServerError)
		return
	}

}

// func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {

// 	p.l.Println("Handle Post Request")
// 	prod := r.Context().Value(KeyProduct{}).(data.Product)
// 	//data.AddProduct(&prod)
// 	//prod := &data.Product{}
// 	// err := prod.FromJSON(r.Body)
// 	// if err != nil {
// 	// 	http.Error(rw, "Unable to unmarshal", http.StatusBadRequest)
// 	// 	return
// 	// }
// 	data.AddProduct(&prod)

//		// p.l.Printf("Prod: %#v", prod)
//	}
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//id := vars["id"]

	ID, err := strconv.Atoi(vars["id"])
	if err != nil {

		http.Error(rw, "unable to convert id", http.StatusBadRequest)
	}

	p.l.Println("Handle Put request for ID ->", ID)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	//prod := &data.Product{}

	// err = prod.FromJSON(r.Body)
	// if err != nil {

	// 	http.Error(rw, "Unable to update product", http.StatusBadRequest)
	// 	return
	// }

	err = data.UpdateProd(ID, &prod)
	if err == data.ErrorNotFound {

		http.Error(rw, "Error not Found", http.StatusBadRequest)
		return
	}

}

func (p *Products) DeleteByID(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//id := vars["id"]
	ID, err := strconv.Atoi(vars["id"])
	prod := &data.Product{}
	if err != nil {

		http.Error(rw, "unable to Convert id ", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle Delete Request by Id and deleted Product of ID ->", ID)
	err = data.DeleteProd(ID, prod)
	if err != nil {
		http.Error(rw, "Unable to Delete Product", http.StatusBadRequest)
		return
	}

}

type KeyProduct struct{}

//The practice of setting up shared functionality that needs to run for many or all HTTP requests is called middleware.

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}
		err = prod.Validate()
		if err != nil {

			p.l.Println("[Error] Validating Product", err)
			http.Error(rw, fmt.Sprintf("Error validating product %s", err),
				http.StatusBadRequest)
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
