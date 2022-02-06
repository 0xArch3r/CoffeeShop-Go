package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"MicroService-Go/data"
	"MicroService-Go/logutil"

	"github.com/gorilla/mux"
)

//This section creates a struct of a Products with the logger
type Products struct {
	logUtil *logutil.MyLogger
}

// Function to create the Products Struct with a logger instance
func NewProducts(l *logutil.MyLogger) *Products {
	return &Products{l}
}

/* This section is for utility functions for CRUD operations (Methods of the Products Struct/Class */

func (p *Products) GetProducts(rw http.ResponseWriter, req *http.Request) {
	p.logUtil.WriteLog(fmt.Sprintf("Handling %v request for %v", req.Method, req.URL.Path), 10)
	productList := data.GetProducts()
	err := productList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, req *http.Request) {
	p.logUtil.WriteLog(fmt.Sprintf("Handling %v request for %v", req.Method, req.URL.Path), 10)

	prod := req.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
	p.logUtil.WriteLog(fmt.Sprintf("Prod: %#v", prod), 10)

	if err := prod.ToJSON(rw); err != nil {
		p.logUtil.WriteLog(fmt.Sprintf("Unable to marshal json"), 10)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, req *http.Request) {
	p.logUtil.WriteLog(fmt.Sprintf("Handling %v request for %v", req.Method, req.URL.Path), 10)
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.logUtil.WriteLog("Unable to convert id", 10)
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	prod := req.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err != nil {
		if err == data.ErrProductNotFound {
			http.Error(rw, "Product not found", http.StatusNotFound)
			return
		}
		errorMsg := "Unable to update product with id: " + strconv.Itoa(id)
		http.Error(rw, errorMsg, http.StatusInternalServerError)
		return
	}

	if err := prod.ToJSON(rw); err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}
}

type KeyProduct struct {
}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.logUtil.WriteLog(fmt.Sprintf("[ERROR] Unable to deserialize product", err), 10)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
