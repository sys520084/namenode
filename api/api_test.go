package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Pong", w.Body.String())
}

func TestAddFile(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	render := strings.NewReader("name=train/test/1.jpg&size=5&")
	render1 := strings.NewReader("name=train/test/2.jpg&size=5&")
	req, _ := http.NewRequest("POST", "/namenode/minist/", render)
	req1, _ := http.NewRequest("POST", "/namenode/minist/", render1)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(w, req)
	router.ServeHTTP(w, req1)

}

//func TestGetDirchildren(t *testing.T) {
//	router := SetupRouter()
//	w := httptest.NewRecorder()
//	render := strings.NewReader("dir=train/0&")
//	req, _ := http.NewRequest("POST", "/namenode/minist/files", render)
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	router.ServeHTTP(w, req)
//
//}
