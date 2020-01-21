package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/guanw/ct-dns/pkg/store/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	httpClient = &http.Client{Timeout: 2 * time.Second}
	metrics    = InitializeMetrics()
)

func Test_decodeBody(t *testing.T) {
	tests := []struct {
		body        io.Reader
		expectedErr bool
		description string
		expected    postBody
	}{
		{
			body:        bytes.NewReader([]byte(`{random:"random"}`)),
			description: "empty body should throw error",
			expectedErr: true,
		},
		{
			body: bytes.NewReader([]byte(`{"serviceName":"dummy-service", "host":"192.0.0.1", "operation":"add"}`)),
			expected: postBody{
				ServiceName: "dummy-service",
				Host:        "192.0.0.1",
				Operation:   "add",
			},
			expectedErr: false,
			description: "body should be parsed correctly",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			b, err := decodeBody(test.body)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, b)
			}
		})
	}
}

func initializeTestServer(store *mocks.Store) *httptest.Server {
	r := mux.NewRouter()
	handler := NewHandler(store, metrics)
	handler.RegisterRoutes(r)
	return httptest.NewServer(r)
}

func makePostReq(t *testing.T, server *httptest.Server, body string, path string) (io.ReadCloser, int) {
	var jsonStr = []byte(body)
	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	res, err := httpClient.Do(req)
	assert.NoError(t, err)
	return res.Body, res.StatusCode
}

func makeGetReq(t *testing.T, server *httptest.Server, path, serviceName string) (io.ReadCloser, int) {
	req, err := http.NewRequest(http.MethodGet, server.URL+path+serviceName, bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)
	res, err := httpClient.Do(req)
	assert.NoError(t, err)
	return res.Body, res.StatusCode
}

func Test_Healthcheck(t *testing.T) {
	server := initializeTestServer(&mocks.Store{})
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/health", bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)
	res, err := httpClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	res.Body.Close()
	assert.Equal(t, 1.0, testutil.ToFloat64(metrics.HealthcheckSuccess))
}

func Test_GetRequest(t *testing.T) {
	mockClient := &mocks.Store{}
	mockClient.On("GetService", "valid-service").Return([]string{"192.0.0.1"}, nil)
	mockClient.On("GetService", "error-service").Return(nil, errors.New("new error"))
	server := initializeTestServer(mockClient)
	defer server.Close()

	getRes, statusCode := makeGetReq(t, server, "/api/service/", "valid-service")
	defer getRes.Close()
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, 1.0, testutil.ToFloat64(metrics.GetServiceSuccess))

	getRes, statusCode = makeGetReq(t, server, "/api/service/", "error-service")
	defer getRes.Close()
	res, err := ioutil.ReadAll(getRes)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(res), "new error"))
	assert.Equal(t, 404, statusCode)
	assert.Equal(t, 1.0, testutil.ToFloat64(metrics.GetServiceFailure))
}

func Test_PostRequest(t *testing.T) {
	mockClient := &mocks.Store{}
	mockClient.On("UpdateService", "valid-service", "add", "192.0.0.1").Return(nil)
	mockClient.On("UpdateService", "error-service", mock.Anything, mock.Anything).Return(errors.New("new error"))
	server := initializeTestServer(mockClient)
	defer server.Close()

	t.Run("POST valid service", func(t *testing.T) {
		postRes, statusCode := makePostReq(t, server, `{"serviceName":"valid-service","operation":"add","host":"192.0.0.1"}`, "/api/service")
		defer postRes.Close()
		assert.Equal(t, 200, statusCode)
		assert.Equal(t, 1.0, testutil.ToFloat64(metrics.PostServiceSuccess))
	})

	t.Run("POST error service", func(t *testing.T) {
		postRes, statusCode := makePostReq(t, server, `{"serviceName":"error-service"}`, "/api/service")
		defer postRes.Close()
		assert.Equal(t, 502, statusCode)
		assert.Equal(t, 1.0, testutil.ToFloat64(metrics.PostServiceFailure))
	})

	t.Run("POST with invalid json", func(t *testing.T) {
		postRes, statusCode := makePostReq(t, server, ``, "/api/service")
		defer postRes.Close()
		res, err := ioutil.ReadAll(postRes)
		assert.NoError(t, err)
		assert.True(t, strings.Contains(string(res), "Failed to decode the Post request body"), "/api/service")
		assert.Equal(t, 422, statusCode)
		assert.Equal(t, 2.0, testutil.ToFloat64(metrics.PostServiceFailure))
	})
}

func Test_RegistrationServiceV1(t *testing.T) {
	mockClient := &mocks.Store{}
	mockClient.On("GetService", "valid-service").Return([]string{"192.0.0.1:8080"}, nil)
	mockClient.On("GetService", "error-service").Return(nil, errors.New("service not found"))
	mockClient.On("GetService", "service-without-port").Return([]string{"192.0.0.1"}, nil)
	mockClient.On("GetService", "service-with-invalid-port").Return([]string{"192.0.0.1:abc"}, nil)
	server := initializeTestServer(mockClient)
	defer server.Close()

	t.Run("get from valid service", func(t *testing.T) {
		validServiceResp, statusCode := makeGetReq(t, server, "/v1/registration/", "valid-service")
		defer validServiceResp.Close()
		assert.Equal(t, 200, statusCode)
		res, err := ioutil.ReadAll(validServiceResp)
		assert.NoError(t, err)
		var resp edsV1Resp
		err = json.Unmarshal(res, &resp)
		assert.NoError(t, err)
		assert.Equal(t, "192.0.0.1", resp.Hosts[0].IPAddress)
		assert.Equal(t, 8080, resp.Hosts[0].Port)
		assert.Equal(t, 1.0, testutil.ToFloat64(metrics.V1RegistrationSuccess))
	})

	t.Run("get from error service", func(t *testing.T) {
		errorServiceResp, statusCode := makeGetReq(t, server, "/v1/registration/", "error-service")
		defer errorServiceResp.Close()
		assert.Equal(t, 404, statusCode)
		assert.Equal(t, 1.0, testutil.ToFloat64(metrics.V1RegistrationFailure))
	})

	t.Run("get without port info", func(t *testing.T) {
		serviceWithoutPortResp, statusCode := makeGetReq(t, server, "/v1/registration/", "service-without-port")
		defer serviceWithoutPortResp.Close()
		assert.Equal(t, 502, statusCode)
		assert.Equal(t, 2.0, testutil.ToFloat64(metrics.V1RegistrationFailure))
	})

	t.Run("get with invalid port", func(t *testing.T) {
		serviceWithInvalidPort, statusCode := makeGetReq(t, server, "/v1/registration/", "service-with-invalid-port")
		defer serviceWithInvalidPort.Close()
		assert.Equal(t, 502, statusCode)
		assert.Equal(t, 3.0, testutil.ToFloat64(metrics.V1RegistrationFailure))
	})
}

func Test_DiscoveryEndpointsV2(t *testing.T) {
	mockClient := &mocks.Store{}
	mockClient.On("GetService", "valid-service").Return([]string{"192.0.0.1:8080"}, nil)
	mockClient.On("GetService", "error-service").Return(nil, errors.New("service not found"))
	mockClient.On("GetService", "service-without-port").Return([]string{"192.0.0.1"}, nil)
	mockClient.On("GetService", "service-with-invalid-port").Return([]string{"192.0.0.1:abc"}, nil)
	server := initializeTestServer(mockClient)
	defer server.Close()

	t.Run("get with invalid body", func(t *testing.T) {
		invalidServiceResp2, statusCode := makePostReq(t, server, `[]`, "/v2/discovery:endpoints")
		defer invalidServiceResp2.Close()
		assert.Equal(t, 422, statusCode)
		assert.Equal(t, 1.0, testutil.ToFloat64(metrics.V2DiscoveryFailure))
	})

	t.Run("get with empty resource names", func(t *testing.T) {
		invalidServiceResp1, statusCode := makePostReq(t, server, `{"resource_names": []}`, "/v2/discovery:endpoints")
		defer invalidServiceResp1.Close()
		assert.Equal(t, 200, statusCode)
		assert.Equal(t, 1.0, testutil.ToFloat64(metrics.V2DiscoverySuccess))
	})

	t.Run("get from valid service", func(t *testing.T) {
		validServiceResp, statusCode := makePostReq(t, server, `{"resource_names":["valid-service"]}`, "/v2/discovery:endpoints")
		defer validServiceResp.Close()
		assert.Equal(t, 200, statusCode)
		assert.Equal(t, 2.0, testutil.ToFloat64(metrics.V2DiscoverySuccess))
		res, err := ioutil.ReadAll(validServiceResp)
		assert.NoError(t, err)
		var resp edsV2Resp
		err = json.Unmarshal(res, &resp)
		assert.NoError(t, err)
		assert.Equal(t, "192.0.0.1", resp.Resources[0].Endpoints[0].LBEndpoints[0].Endpoint.Address.SocketAddress.Address)
		assert.Equal(t, 8080, resp.Resources[0].Endpoints[0].LBEndpoints[0].Endpoint.Address.SocketAddress.PortValue)
	})

	t.Run("get from error service", func(t *testing.T) {
		validServiceResp, statusCode := makePostReq(t, server, `{"resource_names":["error-service"]}`, "/v2/discovery:endpoints")
		defer validServiceResp.Close()
		assert.Equal(t, 404, statusCode)
		assert.Equal(t, 2.0, testutil.ToFloat64(metrics.V2DiscoveryFailure))
	})

	t.Run("get from service without port", func(t *testing.T) {
		validServiceResp, statusCode := makePostReq(t, server, `{"resource_names":["service-without-port"]}`, "/v2/discovery:endpoints")
		defer validServiceResp.Close()
		assert.Equal(t, 502, statusCode)
		assert.Equal(t, 3.0, testutil.ToFloat64(metrics.V2DiscoveryFailure))
	})

	t.Run("get from service with invalid port", func(t *testing.T) {
		validServiceResp, statusCode := makePostReq(t, server, `{"resource_names":["service-with-invalid-port"]}`, "/v2/discovery:endpoints")
		defer validServiceResp.Close()
		assert.Equal(t, 502, statusCode)
		assert.Equal(t, 4.0, testutil.ToFloat64(metrics.V2DiscoveryFailure))
	})
}
