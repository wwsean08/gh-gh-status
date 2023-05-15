package status

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClientReturnsClient(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	require.Nil(t, client.etag)
	require.Equal(t, http.DefaultClient, client.client)
}

func TestClient_GetDataDoesNotPassInEtag(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NotContains(t, r.Header, "If-None-Match")
		w.WriteHeader(200)
	}))
	defer svr.Close()
	client := NewClient()
	client.apiURL = svr.URL

	resp, err := client.getData()
	require.NotNil(t, resp)
	require.NoError(t, err)
}

func TestClient_GetDataDoesPassesInEtag(t *testing.T) {
	expected := "foo"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, expected, r.Header.Get("If-None-Match"))
		w.WriteHeader(200)
	}))
	defer svr.Close()
	client := NewClient()
	client.apiURL = svr.URL
	client.etag = &expected

	resp, err := client.getData()
	require.NotNil(t, resp)
	require.NoError(t, err)
}

func TestClient_GetDataShouldPassInUserAgent(t *testing.T) {
	expected := "gh-status/dev"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, expected, r.Header.Get("User-Agent"))
		w.WriteHeader(200)
	}))
	defer svr.Close()
	client := NewClient()
	client.apiURL = svr.URL
	client.etag = &expected

	resp, err := client.getData()
	require.NotNil(t, resp)
	require.NoError(t, err)
}

func TestClient_PollShouldReturnNilNil(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(304)
	}))
	defer svr.Close()
	client := NewClient()
	client.apiURL = svr.URL

	status, err := client.Poll()
	require.Nil(t, status)
	require.NoError(t, err)
}

func TestClient_PollShouldReturnNilErr(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(418)
	}))
	defer svr.Close()
	client := NewClient()
	client.apiURL = svr.URL

	status, err := client.Poll()
	require.Nil(t, status)
	require.Error(t, err)
}

func TestClient_PollShouldReturnStatusNil(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Etag", "foo")
		w.WriteHeader(200)
		_, err := w.Write([]byte(testJsonResponse))
		require.NoError(t, err)
	}))
	defer svr.Close()
	client := NewClient()
	client.apiURL = svr.URL

	status, err := client.Poll()
	require.NotNil(t, status)
	require.NoError(t, err)
	require.Equal(t, "foo", *client.etag)
}
