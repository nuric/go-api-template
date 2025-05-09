package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nuric/go-api-template/handlers"
	"github.com/stretchr/testify/require"
)

func makeRequest(t *testing.T, hf http.HandlerFunc, method string, endpoint string, body any, resp any) int {
	t.Helper()
	// ---------------------------
	var bodyReader io.Reader
	if body != nil {
		var bodyAsString string
		if bodyString, ok := body.(string); ok {
			bodyAsString = bodyString
		} else {
			jsonBody, err := json.Marshal(body)
			require.NoError(t, err)
			bodyAsString = string(jsonBody)
		}
		bodyReader = bytes.NewReader([]byte(bodyAsString))
	}
	req, err := http.NewRequest(method, endpoint, bodyReader)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	// ---------------------------
	recorder := httptest.NewRecorder()
	hf(recorder, req)
	// ---------------------------
	if resp != nil {
		err = json.Unmarshal(recorder.Body.Bytes(), resp)
		require.NoError(t, err)
	}
	t.Log(recorder.Body.String())
	// ---------------------------
	return recorder.Code
}

func Test_Greetings(t *testing.T) {
	// ---------------------------
	reqBody := handlers.GreetingRequest{
		Name: "Gandalf",
	}
	var respBody handlers.GreetingResponse
	resp := makeRequest(t, handlers.GreetingHandler, http.MethodPost, "/greetings", reqBody, &respBody)
	require.Equal(t, http.StatusOK, resp)
	require.Equal(t, respBody.Greeting, "Hello, Gandalf!")
	// ---------------------------
	// Empty name should return 400
	reqBody.Name = ""
	resp = makeRequest(t, handlers.GreetingHandler, http.MethodPost, "/greetings", reqBody, nil)
	require.Equal(t, http.StatusBadRequest, resp)
	// ---------------------------
}
