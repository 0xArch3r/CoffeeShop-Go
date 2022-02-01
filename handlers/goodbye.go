package handlers

import (
	"log"
	"net/http"
)

type Goodbye struct {
	logUtil *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	g.logUtil.Println("Handling request for /goodbye")
	rw.Write([]byte("Adios mis amigo"))
}
