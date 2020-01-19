package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

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
	router.HandleFunc("/v2/discovery:endpoints", aH.DiscoveryEndpointsV2).Methods(http.MethodPost)
	router.HandleFunc("/v1/registration/{serviceName}", aH.RegistrationServiceV1).Methods(http.MethodGet)
}

// DiscoveryEndpointsV2 process envoy EDS V2 api
func (aH *Handler) DiscoveryEndpointsV2(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Failed to read endpoint v2 body from buff").Error(), http.StatusUnprocessableEntity)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &body); err != nil {
		http.Error(w, errors.Wrap(err, "Failed to decode the eds endpoint v2 request body").Error(), http.StatusUnprocessableEntity)
	}
	resources := body["resource_names"].([]interface{})
	serviceName := resources[0].(string)
	hosts, err := aH.Store.GetService(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	var eps []lbEndpointV2
	for _, url := range hosts {
		splitedStrings := strings.Split(url, ":")
		host := splitedStrings[0]
		port, err := strconv.Atoi(splitedStrings[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		eps = append(eps, lbEndpointV2{
			Endpoint: endpointV2{
				Address: addressV2{
					SocketAddress: socketAddressV2{
						Address:   host,
						PortValue: port,
					},
				},
			},
		})
	}
	resp := edsV2Resp{
		VersionInfo: "v1",
		Resources: []resourceV2{
			{
				Type:        "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
				ClusterName: serviceName,
				Endpoints: []resourceEndpointV2{
					{
						LBEndpoints: eps,
					},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(resp)
}

type edsV2Resp struct {
	VersionInfo string       `json:"version_info"`
	Resources   []resourceV2 `json:"resources"`
}

type resourceV2 struct {
	Type        string               `json:"@type"`
	ClusterName string               `json:"cluster_name"`
	Endpoints   []resourceEndpointV2 `json:"endpoints"`
}

type resourceEndpointV2 struct {
	LBEndpoints []lbEndpointV2 `json:"lb_endpoints"`
}

type lbEndpointV2 struct {
	Endpoint endpointV2 `json:"endpoint"`
}

type endpointV2 struct {
	Address addressV2 `json:"address"`
}

type addressV2 struct {
	SocketAddress socketAddressV2 `json:"socket_address"`
}

type socketAddressV2 struct {
	Address   string `json:"address"`
	PortValue int    `json:"port_value"`
}

// RegistrationServiceV1 process envoy EDS V1 api
func (aH *Handler) RegistrationServiceV1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["serviceName"]
	hosts, err := aH.Store.GetService(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	var hostsV1 []hostV1
	for _, h := range hosts {
		splitStrings := strings.Split(h, ":")
		port, err := strconv.Atoi(splitStrings[1])
		if err != nil {
			http.Error(w, "Failed to parse port from host info", http.StatusUnprocessableEntity)
		}
		hostsV1 = append(hostsV1, hostV1{
			IPAddress: splitStrings[0],
			Port:      port,
			Tags: tagsV1{
				AZ:                  "default", //TODO support availability zone
				Canary:              false,
				LoadBalancingWeight: 1, // TODO check how this load balancing weight works
			},
		})
	}
	resp := edsV1Resp{
		Hosts: hostsV1,
	}
	json.NewEncoder(w).Encode(resp)
}

type edsV1Resp struct {
	Hosts []hostV1 `json:"version_info"`
}

type hostV1 struct {
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	Tags      tagsV1 `json:"tags"`
}

type tagsV1 struct {
	AZ                  string `json:"az"`
	Canary              bool   `json:"canary"`
	LoadBalancingWeight int    `json:"load_balancing_weight"`
}

// HealthService process healthcheck GET request
func (aH *Handler) HealthService(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
	ServiceName string `json:"serviceName"`
	Operation   string `json:"operation"`
	Host        string `json:"host"`
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
