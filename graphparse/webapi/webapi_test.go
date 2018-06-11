package webapi

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	// "fmt"
	"encoding/base64"
)


type dict map[string]interface{}

func TestGraphGenerated(t *testing.T) {
	repoId := base64.StdEncoding.EncodeToString([]byte("github.com/liamzebedee/graphparse/graphparse"))
	req, _ := http.NewRequest("GET", "/graph/public/" + repoId, nil)
	rawRes := httptest.NewRecorder()
	webAPI().ServeHTTP(rawRes, req)

	assert.Equal(t, http.StatusOK, rawRes.Code, rawRes.Body.String())
	res := dict{}
	err := json.NewDecoder(rawRes.Body).Decode(&res)
	assert.NoError(t, err)
	assert.NotNil(t, res["nodes"])
}

func TestPackageNotFound(t *testing.T) {
	repoId := base64.StdEncoding.EncodeToString([]byte("DTA"))
	req, _ := http.NewRequest("GET", "/graph/public/" + repoId, nil)
	rawRes := httptest.NewRecorder()
	webAPI().ServeHTTP(rawRes, req)

	assert.Equal(t, http.StatusInternalServerError, rawRes.Code, rawRes.Body.String())
}