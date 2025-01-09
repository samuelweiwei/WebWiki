package server

import "net/http"

type RouteMap struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
}
