package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	handler := HealthCheck()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)

	data, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "ok", data["status"])
}

func TestResponseHelpers(t *testing.T) {
	t.Run("Success returns 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		Success(w, map[string]string{"test": "value"})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

		var resp Response
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("Created returns 201", func(t *testing.T) {
		w := httptest.NewRecorder()
		Created(w, map[string]string{"id": "123"})

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp Response
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("BadRequest returns 400", func(t *testing.T) {
		w := httptest.NewRecorder()
		BadRequest(w, "invalid input")

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp Response
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		assert.False(t, resp.Success)
		assert.Equal(t, "invalid input", resp.Error)
	})

	t.Run("NotFound returns 404", func(t *testing.T) {
		w := httptest.NewRecorder()
		NotFound(w, "resource not found")

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp Response
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		assert.False(t, resp.Success)
		assert.Equal(t, "resource not found", resp.Error)
	})

	t.Run("Conflict returns 409", func(t *testing.T) {
		w := httptest.NewRecorder()
		Conflict(w, "resource already exists")

		assert.Equal(t, http.StatusConflict, w.Code)

		var resp Response
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		assert.False(t, resp.Success)
		assert.Equal(t, "resource already exists", resp.Error)
	})

	t.Run("InternalError returns 500", func(t *testing.T) {
		w := httptest.NewRecorder()
		InternalError(w, "something went wrong")

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp Response
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		assert.False(t, resp.Success)
		assert.Equal(t, "something went wrong", resp.Error)
	})

	t.Run("NoContent returns 204", func(t *testing.T) {
		w := httptest.NewRecorder()
		NoContent(w)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
