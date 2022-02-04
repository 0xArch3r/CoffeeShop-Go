package handlers

import (
	"log"
	"net/http"
	"strconv"

	"MicroService-Go/data"

	"github.com/gorilla/mux"
)

//This section creates a struct of a Products with the logger
type Products struct {
	logUtil *log.Logger
}

// Function to create the Products Struct with a logger instance
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

/* This section is for utility functions for CRUD operations (Methods of the Products Struct/Class */

func (p *Products) GetProducts(rw http.ResponseWriter, req *http.Request) {
	p.logUtil.Printf("Handling %v request for %v", req.Method, req.URL.Path)
	productList := data.GetProducts()
	err := productList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, req *http.Request) {
	p.logUtil.Printf("Handling %v request for %v", req.Method, req.URL.Path)
	requestBody := req.Body
	defer requestBody.Close()

	prod := &data.Product{}
	err := prod.FromJSON(requestBody)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json payload", http.StatusInternalServerError)
		return
	}

	data.AddProduct(prod)
	p.logUtil.Printf("Prod: %#v", prod)

	if err := prod.ToJSON(rw); err != nil {
		p.logUtil.Printf("Unable to marshal json")
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, req *http.Request) {
	p.logUtil.Printf("Handling %v request for %v", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.logUtil.Println("Unable to convert id")
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	requestBody := req.Body
	defer requestBody.Close()

	prod := &data.Product{}
	err = prod.FromJSON(requestBody)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json payload", http.StatusBadRequest)
		return
	}
	prod.ID = id

	err = data.UpdateProduct(id, prod)
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
