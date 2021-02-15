package lib

import (
	"net/http"
)

func NotAuthorized(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("\nNo token given\n"))
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func RenderApp(w http.ResponseWriter, r *http.Request, a []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(a)
}

func RenderAllApps(w http.ResponseWriter, r *http.Request, a []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(a)
}

func Empty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("\nNone\n"))
}

func Failed(w http.ResponseWriter, req *http.Request, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func NotFound(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("\nNot Found\n"))
}

func FailedTimeout(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(408)
	w.Write([]byte(`{"message": "Command Timeout"}`))
}

func NotSupported(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusUnsupportedMediaType)
	w.Write([]byte("Unsupported Media Type. Only JSON files are allowed"))
}

func InternalError(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Error"))
}

func ActionDone(w http.ResponseWriter, req *http.Request, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func Added(w http.ResponseWriter, req *http.Request, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
}
