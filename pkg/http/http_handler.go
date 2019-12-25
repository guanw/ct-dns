package http

import (
	"encoding/json"
	"net/http"

	"github.com/clicktherapeutics/ct-dns/pkg/etcd"
	"github.com/gorilla/mux"
)

// Result is used to parse hosts
type Result struct {
	Hosts []string `json:'hosts'`
}

// Handler use etcd client to query
type Handler struct {
	Client *etcd.Client
}

// NewHandler creates a new Handler
func NewHandler(client *etcd.Client) *Handler {
	return &Handler{
		Client: client,
	}
}

// RegisterRoutes registers GetService with router
func (aH *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/service/{serviceName}", aH.GetService).Methods(http.MethodGet)
}

// GetService process GET request
func (aH *Handler) GetService(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		serviceName := vars["serviceName"]
		res, err := aH.Client.Get(serviceName)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			http.Error(w, "Service Name not found", http.StatusNotFound)
			return
		}
		hosts := unmarshalStrToHosts(res)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(hosts)
	default:
		http.Error(w, "Unsupported Operation", http.StatusNotFound)
	}
}

// TODO
func marshalHostsToStr() {

}

func unmarshalStrToHosts(hosts string) *Result {
	return &Result{
		Hosts: []string{"192.0.0.1", "192.0.0.2"},
	}
}
