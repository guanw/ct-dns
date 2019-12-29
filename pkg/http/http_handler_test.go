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

	"github.com/clicktherapeutics/ct-dns/pkg/store/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var httpClient = &http.Client{Timeout: 2 * time.Second}

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
	handler := NewHandler(store)
	handler.RegisterRoutes(r)
	return httptest.NewServer(r)
}

func makePostReq(t *testing.T, body interface{}, server *httptest.Server) (io.ReadCloser, int) {
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/service", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	res, err := httpClient.Do(req)
	assert.NoError(t, err)
	return res.Body, res.StatusCode
}

func makeGetReq(t *testing.T, server *httptest.Server, serviceName string) (io.ReadCloser, int) {
	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/service/"+serviceName, bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)
	res, err := httpClient.Do(req)
	assert.NoError(t, err)
	return res.Body, res.StatusCode
}

func Test_GetRequest(t *testing.T) {
	mockClient := &mocks.Store{}
	mockClient.On("GetService", "valid-service").Return([]string{"192.0.0.1"}, nil)
	mockClient.On("GetService", "error-service").Return(nil, errors.New("new error"))
	server := initializeTestServer(mockClient)
	defer server.Close()

	getRes, statusCode := makeGetReq(t, server, "valid-service")
	defer getRes.Close()
	assert.Equal(t, 200, statusCode)

	getRes, statusCode = makeGetReq(t, server, "error-service")
	defer getRes.Close()
	res, err := ioutil.ReadAll(getRes)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(res), "new error"))
	assert.Equal(t, 500, statusCode)
}

func Test_PostRequest(t *testing.T) {
	mockClient := &mocks.Store{}
	mockClient.On("UpdateService", "valid-service", "add", "192.0.0.1").Return(nil)
	mockClient.On("UpdateService", "error-service", mock.Anything, mock.Anything).Return(errors.New("new error"))
	server := initializeTestServer(mockClient)
	defer server.Close()

	postRes, statusCode := makePostReq(t, postBody{
		ServiceName: "valid-service",
		Operation:   "add",
		Host:        "192.0.0.1",
	}, server)
	defer postRes.Close()
	assert.Equal(t, 200, statusCode)

	postRes, statusCode = makePostReq(t, postBody{
		ServiceName: "error-service",
	}, server)
	defer postRes.Close()
	assert.Equal(t, 500, statusCode)

	postRes, statusCode = makePostReq(t, "", server)
	defer postRes.Close()
	res, err := ioutil.ReadAll(postRes)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(res), "Failed to decode the Post request body"))
	assert.Equal(t, 500, statusCode)
}
