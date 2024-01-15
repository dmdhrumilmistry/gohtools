package http

import (
	"fmt"
	"log"
	"net/http"
)

type LoggerMiddleware struct {
	Handler http.Handler
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// start block
	log.Println("start")
	l.Handler.ServeHTTP(w, r)
	// stop
	log.Println("stop")

	log.Println(r.RemoteAddr, r.RequestURI)
}

// Trivial Middleware example
type TrivialMiddleware struct{}

func (t *TrivialMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("Trivial Middleware!!")
	next(w, r)
}
