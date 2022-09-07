package app

import (
	"bytes"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"main/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_ApiRecommend(t *testing.T) {
	a := NewApp(config.Config{Db: struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Database string `toml:"database"`
	}{Host: "localhost", Port: "5432", Username: "Artem", Password: "1234", Database: "postgres"}})

	tests := []struct {
		about    string
		reqBody  []byte
		respCode int
	}{
		{
			about:    "TEST 1: Valid. Geo data not specified",
			reqBody:  []byte(`{}`),
			respCode: 200,
		},
		{
			about:    "TEST 2: Valid",
			reqBody:  []byte(`{"lon":37.666, "lat":55.674}`),
			respCode: 200,
		},
		{
			about:    "Test 3: Bad JSON: No JSON",
			reqBody:  nil,
			respCode: 400,
		},
		{
			about:    "Test 4: Bad JSON: Wrong type: string",
			reqBody:  []byte(`{"lon":"37.666", "lat":"55.674"}`),
			respCode: 400,
		},
		{
			about:    "Test 5: Bad JSON: Empty values",
			reqBody:  []byte(`{"lon":, "lat":}`),
			respCode: 400,
		},
	}
	for i, val := range tests {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/recommend", bytes.NewBuffer(val.reqBody))
		router := mux.NewRouter()
		router.HandleFunc("/api/recommend", a.ApiRecommend)
		router.ServeHTTP(res, req)

		if res.Code != val.respCode {
			t.Error("Test", i+1, "FAIL", "\nExpected:", val.respCode, "Got:", res.Code, "\nBody: ", res.Body)
		} else {
			log.Infoln("Test", i+1, "OK", "\nExpected:", val.respCode, "Got:", res.Code, "\nBody: ", res.Body)
		}
	}
}
