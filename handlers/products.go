package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"MicroService-Go/data"
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

func (p *Products) getProducts(rw http.ResponseWriter, req *http.Request) {
	productList := data.GetProducts()
	err := productList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, req *http.Request) {
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
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (p *Products) updateProduct(rw http.ResponseWriter, req *http.Request, id int) {
	requestBody := req.Body
	defer requestBody.Close()

	prod := &data.Product{}
	err := prod.FromJSON(requestBody)
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

/* ServeHTTP handles the different request operations and calls its respective CRUD Operation
 * This ServeHTTP is what makes our Products Struct a handler.
 */
func (p *Products) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p.logUtil.Printf("Handling %v request for %v\n", req.Method, req.RequestURI)

	if req.Method == http.MethodGet {
		p.getProducts(rw, req)
		return
	}

	if req.Method == http.MethodPost {
		p.addProduct(rw, req)
		return
	}

	if req.Method == http.MethodPut {
		r := regexp.MustCompile(`/products/([0-9]+)$`)
		matches := r.FindAllStringSubmatch(req.URL.Path, -1)

		if len(matches) != 1 {
			p.logUtil.Println("Error: Invalid URI too many matches", req.URL.Path)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(matches[0]) != 2 {
			p.logUtil.Println("Error: Invalid URI", req.URL.Path)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := matches[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.logUtil.Println("Error: Invalid URI", req.URL.Path)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(rw, req, id)
		p.logUtil.Printf("Received request for id: %v", id)
		return
	}

	p.logUtil.Printf("Error %v: Method %v not allowed", http.StatusMethodNotAllowed, req.Method)
	rw.WriteHeader(http.StatusMethodNotAllowed)
}
