package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMainHandlerEmptyCountParameter(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/?city=moscow", nil)

	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "count missing", responseRecorder.Body.String())
}

func TestMainHandlerUnsupportedCities(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/cafe?city=ivanovo&count=100500", nil)

	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "wrong city value")
}

func TestMainHandlerCounterNotIsNumber(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/?city=moscow&counter=абырвалг", nil)

	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "count missing", responseRecorder.Body.String())
}

func TestMainHandlerCorrectRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/cafe?&city=moscow&count=2", nil)

	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.NotEmpty(t, responseRecorder.Body.String())
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/cafe?city=moscow&count=10", nil)

	require.Nil(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	data, err := io.ReadAll(responseRecorder.Body)

	require.Equal(t, responseRecorder.Code, http.StatusOK)
	assert.Equal(t, len(strings.Split(string(data), ",")), len(cafeList["moscow"]))
}
