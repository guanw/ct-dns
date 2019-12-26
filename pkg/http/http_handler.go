package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/clicktherapeutics/ct-dns/pkg/etcd"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

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
	router.HandleFunc("/api/service", aH.PostService).Methods(http.MethodPost)
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
			http.Error(w, errors.Wrap(err, "Service Name not found").Error(), http.StatusInternalServerError)
			return
		}
		hosts, err := unmarshalStrToHosts(res)
		if err != nil {
			http.Error(w, errors.Wrap(err, "UnmarshalStrToHosts failed").Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(hosts)
	default:
		http.Error(w, "Unsupported Request Operation", http.StatusInternalServerError)
	}
}

// PostService process POST request
func (aH *Handler) PostService(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b, err := decodeBody(r.Body)
		if err != nil {
			http.Error(w, errors.Wrap(err, "Failed to decode the Post request body").Error(), http.StatusInternalServerError)
			return
		}

		res, err := aH.Client.Get(b.ServiceName)
		if err != nil {
			if b.Operation == "add" {
				err := aH.Client.CreateOrSet(b.ServiceName, fmt.Sprintf(`["%s"]`, b.Host))
				if err != nil {
					http.Error(w, errors.Wrap(err, "Failed to create new service entry").Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Error(w, errors.Wrap(err, "Failed to delete service: Service name not found").Error(), http.StatusInternalServerError)
			return
		}
		hosts, err := unmarshalStrToHosts(res)
		if err != nil {
			http.Error(w, errors.Wrap(err, "UnmarshalStrToHosts failed").Error(), http.StatusInternalServerError)
			return
		}
		if b.Operation == "add" {
			include, _ := contains(hosts, b.Host)
			if !include {
				hosts = append(hosts, b.Host)
				string, err := marshalHostsToStr(hosts)
				if err != nil {
					http.Error(w, errors.Wrap(err, "Failed to marshal hosts to string").Error(), http.StatusInternalServerError)
					return
				}
				err = aH.Client.CreateOrSet(b.ServiceName, string)
				if err != nil {
					http.Error(w, errors.Wrap(err, "Failed to add host").Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				return
			}
		} else if b.Operation == "delete" {
			include, i := contains(hosts, b.Host)
			if include {
				hosts = append(hosts[:i], hosts[i+1:]...)
				string, err := marshalHostsToStr(hosts)
				if err != nil {
					http.Error(w, errors.Wrap(err, "Failed to marshal hosts to string").Error(), http.StatusInternalServerError)
					return
				}
				err = aH.Client.CreateOrSet(b.ServiceName, string)
				if err != nil {
					http.Error(w, errors.Wrap(err, "Failed to delete host").Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
	default:
		http.Error(w, "Unsupported Request Operation", http.StatusInternalServerError)
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

func contains(s []string, e string) (bool, int) {
	for i, a := range s {
		if a == e {
			return true, i
		}
	}
	return false, -1
}

func unmarshalStrToHosts(input string) ([]string, error) {
	var hosts []string
	if err := json.Unmarshal([]byte(input), &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

func marshalHostsToStr(input []string) (string, error) {
	res, err := json.Marshal(input)
	return string(res), err
}
