package response

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestNewAPIResponse(t *testing.T) {
	t.Run("creates response with payload", func(t *testing.T) {
		data := &TestData{Name: "John", Age: 30}
		resp := NewAPIResponse(data)

		assert.NotNil(t, resp.Payload)
		assert.Equal(t, data, resp.Payload)
		assert.Empty(t, resp.Errors)
		assert.Nil(t, resp.Pagination)
	})

	t.Run("creates response with nil payload", func(t *testing.T) {
		var data *TestData = nil
		resp := NewAPIResponse(data)

		assert.Nil(t, resp.Payload)
		assert.Empty(t, resp.Errors)
	})
}

func TestNewAPIListResponse(t *testing.T) {
	t.Run("creates list response with pagination", func(t *testing.T) {
		data := []TestData{
			{Name: "John", Age: 30},
			{Name: "Jane", Age: 25},
		}
		resp := NewAPIListResponse(data, 1, 1, 10, 2)

		assert.Equal(t, data, resp.Payload)
		assert.Empty(t, resp.Errors)
		assert.NotNil(t, resp.Pagination)
		assert.Equal(t, 1, resp.Pagination.Page)
		assert.Equal(t, 1, resp.Pagination.CurrentPage)
		assert.Equal(t, 10, resp.Pagination.PerPage)
		assert.Equal(t, 2, resp.Pagination.LastPage)
	})

	t.Run("creates empty list response", func(t *testing.T) {
		data := []TestData{}
		resp := NewAPIListResponse(data, 1, 1, 10, 0)

		assert.Empty(t, resp.Payload)
		assert.NotNil(t, resp.Pagination)
		assert.Equal(t, 0, resp.Pagination.LastPage)
	})
}

func TestNewAPIErrorResponse(t *testing.T) {
	t.Run("creates error response", func(t *testing.T) {
		resp := NewAPIErrorResponse[TestData]("Validation Failed", "Email is required")

		assert.Nil(t, resp.Payload)
		assert.Len(t, resp.Errors, 1)
		assert.Equal(t, "Validation Failed", resp.Errors[0].Error)
		assert.Equal(t, "Email is required", resp.Errors[0].Message)
	})

	t.Run("creates error response with empty messages", func(t *testing.T) {
		resp := NewAPIErrorResponse[TestData]("", "")

		assert.Len(t, resp.Errors, 1)
		assert.Equal(t, "", resp.Errors[0].Error)
		assert.Equal(t, "", resp.Errors[0].Message)
	})
}

func TestAPIResponseJSONSerialization(t *testing.T) {
	t.Run("serializes APIResponse correctly", func(t *testing.T) {
		data := &TestData{Name: "John", Age: 30}
		resp := NewAPIResponse(data)

		jsonBytes, err := json.Marshal(resp)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err)

		assert.Contains(t, result, "payload")
		assert.Contains(t, result, "error")
	})

	t.Run("serializes APIListResponse correctly", func(t *testing.T) {
		data := []TestData{{Name: "John", Age: 30}}
		resp := NewAPIListResponse(data, 1, 1, 10, 1)

		jsonBytes, err := json.Marshal(resp)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err)

		assert.Contains(t, result, "payload")
		assert.Contains(t, result, "pagination")
	})

	t.Run("serializes APIErrorResponse correctly", func(t *testing.T) {
		resp := NewAPIErrorResponse[TestData]("Error", "Message")

		jsonBytes, err := json.Marshal(resp)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err)

		assert.Contains(t, result, "error")
		errors := result["error"].([]interface{})
		assert.Len(t, errors, 1)
	})
}

func TestPaginationOmitempty(t *testing.T) {
	t.Run("payload response omits pagination", func(t *testing.T) {
		data := &TestData{Name: "John", Age: 30}
		resp := NewAPIResponse(data)

		jsonBytes, err := json.Marshal(resp)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err)

		_, hasPagination := result["pagination"]
		assert.False(t, hasPagination)
	})

	t.Run("list response includes pagination", func(t *testing.T) {
		data := []TestData{{Name: "John", Age: 30}}
		resp := NewAPIListResponse(data, 1, 1, 10, 1)

		jsonBytes, err := json.Marshal(resp)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err)

		_, hasPagination := result["pagination"]
		assert.True(t, hasPagination)
	})
}