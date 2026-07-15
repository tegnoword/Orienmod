package router

import (
	"net/http"
)

type Router struct {
	classroomAdapter interface{}
}

func NewRouter(classroom interface{}) *Router {
	return &Router{classroomAdapter: classroom}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, QUERY, OPTIONS") // Soportamos QUERY explícitamente

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.URL.Path {
	case "/api/v1/classroom/webhook":
		w.Write([]byte(`{"status": "webhook listo"}`))

	case "/api/v1/students/search":
		if r.Method == "QUERY" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "Endpoint QUERY configurado bajo RFC 10008"}`))
		} else {
			http.Error(w, "Método no permitido. Use QUERY para búsquedas.", http.StatusMethodNotAllowed)
		}

	default:
		http.NotFound(w, r)
	}
}
