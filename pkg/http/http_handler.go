package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guanw/ct-dns/pkg/store"
	"github.com/pkg/errors"
)

// Handler use etcd client to query
type Handler struct {
	Store store.Store
}

// NewHandler creates a new Handler
func NewHandler(store store.Store) *Handler {
	return &Handler{
		Store: store,
	}
}

// RegisterRoutes registers GetService with router
func (aH *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/service/{serviceName}", aH.GetService).Methods(http.MethodGet)
	router.HandleFunc("/api/service", aH.PostService).Methods(http.MethodPost)
	router.HandleFunc("/api/health", aH.HealthService).Methods(http.MethodGet)
}

// HealthService process healthcheck GET request
func (aH *Handler) HealthService(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

// GetService process GET service request
func (aH *Handler) GetService(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		serviceName := vars["serviceName"]
		hosts, err := aH.Store.GetService(serviceName)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(hosts)
	default:
		http.Error(w, "Unsupported Request Operation", http.StatusMethodNotAllowed)
	}
}

// PostService process POST service request
func (aH *Handler) PostService(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b, err := decodeBody(r.Body)
		if err != nil {
			http.Error(w, errors.Wrap(err, "Failed to decode the Post request body").Error(), http.StatusUnprocessableEntity)
			return
		}

		_ = aH.Store.UpdateService(b.ServiceName, b.Operation, b.Host)
		// TODO think of way to log logic error and not panic
		w.Header().Set("Content-Type", "application/json")
	default:
		http.Error(w, "Unsupported Request Operation", http.StatusMethodNotAllowed)
	}
}

type postBody struct {
	ServiceName string `json:'serviceName'`
	Operation   string `json:'operation'`
	Host        string `json:'host'`
}

func decodeBody(in io.Reader) (postBody, error) {
	var b postBody
	decoder := json.NewDecoder(in)
	err := decoder.Decode(&b)
	if err != nil {
		return postBody{}, err
	}
	return b, nil
}
