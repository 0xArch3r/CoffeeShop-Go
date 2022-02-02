package handlers

import (
	"log"
	"net/http"

	"MicroService-Go/data"
)

type Products struct {
	logUtil *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p.logUtil.Printf("Handling %v request for %v\n", req.Method, req.RequestURI)
	if req.Method == http.MethodGet {
		p.getProducts(rw, req)
		return
	}
	p.logUtil.Printf("Error %v: Method %v not allowed", http.StatusMethodNotAllowed, req.Method)
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, req *http.Request) {
	productList := data.GetProducts()
	err := productList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}
