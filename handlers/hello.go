package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	logUtil *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.logUtil.Println("Handling request for /")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Something went wrong", http.StatusBadRequest)
		h.logUtil.Println("Something went wrong")
		return
	}
	defer r.Body.Close()

	fmt.Fprintf(rw, "Hola %s", data)
}
