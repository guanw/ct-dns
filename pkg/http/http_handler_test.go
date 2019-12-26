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

	"github.com/stretchr/testify/assert"
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

func Test_unmarshalStrToHosts(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr bool
		description string
		expected    []string
	}{
		{
			input:       `["192.0.0.1","192.0.0.2"]`,
			expectedErr: false,
			expected:    []string{"192.0.0.1", "192.0.0.2"},
			description: `input: "[192.0.0.1, 192.0.0.2]"`,
		},
		{
			input:       `""`,
			expectedErr: true,
			description: "input invalid",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			out, err := unmarshalStrToHosts(test.input)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, out)
			}
		})
	}
}

func Test_marshalHostsToStr(t *testing.T) {
	tests := []struct {
		input       []string
		expectedErr bool
		description string
		expected    string
	}{
		{
			input:       []string{"192.0.0.1", "192.0.0.2"},
			expectedErr: false,
			expected:    `["192.0.0.1","192.0.0.2"]`,
			description: "input: [192.0.0.1, 192.0.0.2]",
		},
		{
			input:       nil,
			expectedErr: false,
			expected:    "null",
			description: "input invalid",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			out, err := marshalHostsToStr(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, out)
		})
	}
}

type mockETCDClient struct {
	m map[string]string
}

func newMockETCDClient() *mockETCDClient {
	return &mockETCDClient{
		m: make(map[string]string),
	}
}

func (c *mockETCDClient) CreateOrSet(key, value string) error {
	c.m[key] = value
	return nil
}

func (c *mockETCDClient) Get(key string) (string, error) {
	val, found := c.m[key]
	if !found {
		return "", errors.New("invalid service")
	}
	return val, nil
}

func initializeTestServer(etcdClient *mockETCDClient) *httptest.Server {
	r := mux.NewRouter()
	handler := NewHandler(etcdClient)
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

func Test_GetService(t *testing.T) {
	mockClient := newMockETCDClient()
	mockClient.CreateOrSet("dummy-service", `["192.0.0.1"]`)
	mockClient.CreateOrSet("empty-service", `[]`)
	mockClient.CreateOrSet("invalid-service", ``)
	server := initializeTestServer(mockClient)
	defer server.Close()

	tests := []struct {
		expectedErr      bool
		expectedResponse string
		serviceName      string
	}{
		{
			serviceName:      "dummy-service",
			expectedErr:      false,
			expectedResponse: `["192.0.0.1"]`,
		},
		{
			serviceName: "invalid-service",
			expectedErr: true,
		},
		{
			serviceName:      "empty-service",
			expectedErr:      false,
			expectedResponse: `[]`,
		},
		{
			serviceName: "non-exist",
			expectedErr: true,
		},
	}

	for _, test := range tests {
		res, statusCode := makeGetReq(t, server, test.serviceName)
		defer res.Close()
		if test.expectedErr {
			assert.Equal(t, 500, statusCode)
		} else {
			body, err := ioutil.ReadAll(res)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, strings.TrimSpace(string(body)))
		}
	}
}

func Test_ExistingServiceAddNewHost(t *testing.T) {
	mockClient := newMockETCDClient()
	server := initializeTestServer(mockClient)
	defer server.Close()

	postRes1, statusCode := makePostReq(t, &postBody{
		ServiceName: "dummy-service",
		Operation:   "add",
		Host:        "192.0.0.1",
	}, server)
	assert.Equal(t, 200, statusCode)
	defer postRes1.Close()
	postRes2, statusCode := makePostReq(t, &postBody{
		ServiceName: "dummy-service",
		Operation:   "add",
		Host:        "192.0.0.2",
	}, server)
	assert.Equal(t, 200, statusCode)
	defer postRes2.Close()

	res, statusCode := makeGetReq(t, server, "dummy-service")
	defer res.Close()
	assert.Equal(t, 200, statusCode)
	result, err := ioutil.ReadAll(res)
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1","192.0.0.2"]`, strings.TrimSpace(string(result)))
}

func Test_ExistingServiceAddExistingHost(t *testing.T) {
	mockClient := newMockETCDClient()
	server := initializeTestServer(mockClient)
	defer server.Close()

	postRes1, statusCode := makePostReq(t, &postBody{
		ServiceName: "dummy-service",
		Operation:   "add",
		Host:        "192.0.0.1",
	}, server)
	assert.Equal(t, 200, statusCode)
	defer postRes1.Close()
	postRes2, statusCode := makePostReq(t, &postBody{
		ServiceName: "dummy-service",
		Operation:   "add",
		Host:        "192.0.0.1",
	}, server)
	defer postRes2.Close()
	postRes2Result, err := ioutil.ReadAll(postRes2)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to add Host: Host already existed", strings.TrimSpace(string(postRes2Result)))
	assert.Equal(t, 500, statusCode)

	res, statusCode := makeGetReq(t, server, "dummy-service")
	defer res.Close()
	assert.Equal(t, 200, statusCode)
	result, err := ioutil.ReadAll(res)
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1"]`, strings.TrimSpace(string(result)))
}

func Test_NonExistingServiceDeleteHost(t *testing.T) {
	mockClient := newMockETCDClient()
	server := initializeTestServer(mockClient)
	defer server.Close()

	postRes, statusCode := makePostReq(t, &postBody{
		ServiceName: "non-exist-service",
		Operation:   "delete",
		Host:        "192.0.0.1",
	}, server)
	defer postRes.Close()
	assert.Equal(t, 500, statusCode)
	postResResult, err := ioutil.ReadAll(postRes)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(postResResult), "Failed to delete service: Service name not found"))
}

func Test_ExistingServiceDeleteNonExistingHost(t *testing.T) {
	mockClient := newMockETCDClient()
	server := initializeTestServer(mockClient)
	defer server.Close()

	postRes, statusCode := makePostReq(t, &postBody{
		ServiceName: "dummy-service",
		Operation:   "add",
		Host:        "192.0.0.1",
	}, server)
	defer postRes.Close()
	assert.Equal(t, 200, statusCode)

	postRes2, statusCode := makePostReq(t, &postBody{
		ServiceName: "dummy-service",
		Operation:   "delete",
		Host:        "192.0.0.2",
	}, server)
	defer postRes2.Close()
	assert.Equal(t, 500, statusCode)

	postResResult, err := ioutil.ReadAll(postRes2)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(postResResult), "Failed to delete Host: Host not found"))
}

func Test_ExistingServiceDeleteExistingHost(t *testing.T) {
	mockClient := newMockETCDClient()
	server := initializeTestServer(mockClient)
	defer server.Close()

	postRes, statusCode := makePostReq(t, &postBody{
		ServiceName: "dummy-service",
		Operation:   "add",
		Host:        "192.0.0.1",
	}, server)
	defer postRes.Close()
	assert.Equal(t, 200, statusCode)

	postRes, statusCode = makePostReq(t, &postBody{
		ServiceName: "dummy-service",
		Operation:   "add",
		Host:        "192.0.0.2",
	}, server)
	defer postRes.Close()
	assert.Equal(t, 200, statusCode)

	postRes2, statusCode := makePostReq(t, &postBody{
		ServiceName: "dummy-service",
		Operation:   "delete",
		Host:        "192.0.0.1",
	}, server)
	defer postRes2.Close()
	assert.Equal(t, 200, statusCode)

	res, statusCode := makeGetReq(t, server, "dummy-service")
	defer res.Close()
	assert.Equal(t, 200, statusCode)
	result, err := ioutil.ReadAll(res)
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.2"]`, strings.TrimSpace(string(result)))
}

func Test_invalidPostBody(t *testing.T) {
	mockClient := newMockETCDClient()
	server := initializeTestServer(mockClient)
	defer server.Close()

	postRes, statusCode := makePostReq(t, "", server)
	defer postRes.Close()
	assert.Equal(t, 500, statusCode)
	postResResult, err := ioutil.ReadAll(postRes)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(postResResult), "Failed to decode the Post request body"))
}
