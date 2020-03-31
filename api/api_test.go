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
	render := strings.NewReader("name=train&isdir=true&size=5&")
	req, _ := http.NewRequest("POST", "/namenode/minist/", render)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

}
