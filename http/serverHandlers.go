package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Router struct{}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/admin":
		fmt.Fprintln(w, "Admin Page")
	case "/user":
		fmt.Fprintln(w, "User Page")
	default:
		http.Error(w, "404 Not Found", 404)
	}
}

func EchoParam(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprintln(w, "Invalid Param!!")
		return
	}
	name = strings.ReplaceAll(name, "<", "")
	name = strings.ReplaceAll(name, ">", "")
	fmt.Fprintf(w, "Hello %s", name)
}

func EchoUserParam(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	if user == "" {
		fmt.Fprintln(w, "Invalid Param!!")
		return
	}
	fmt.Fprintf(w, "Hello %s!!\n", user)
}
