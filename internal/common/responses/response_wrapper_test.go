package responses

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestResponseJSON tests the Response function for JSON responses
func TestResponseJSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := map[string]interface{}{
		"key": "value",
	}

	err := Response(c, http.StatusOK, data)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, echo.MIMEApplicationJSON, rec.Header().Get(echo.HeaderContentType))

	var responseData map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, data, responseData)
}

// TestResponseXML tests the Response function for XML responses
func TestResponseXML(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationXML)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := struct {
		XMLName xml.Name `xml:"response"`
		Key     string   `xml:"key"`
	}{
		Key: "value",
	}

	err := Response(c, http.StatusOK, data)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, echo.MIMEApplicationXMLCharsetUTF8, rec.Header().Get(echo.HeaderContentType))

	var responseData struct {
		Key string `xml:"key"`
	}
	err = xml.Unmarshal(rec.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, data.Key, responseData.Key)
}

// TestErrorResponseWithMessage tests the ErrorResponseWithMessage function
func TestErrorResponseWithMessage(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	statusCode := http.StatusBadRequest
	message := "An error occurred"

	err := ErrorResponseWithMessage(c, statusCode, message)
	assert.NoError(t, err)
	assert.Equal(t, statusCode, rec.Code)
	assert.Equal(t, echo.MIMEApplicationJSON, rec.Header().Get(echo.HeaderContentType))

	var responseData Error
	err = json.Unmarshal(rec.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, statusCode, responseData.Code)
	assert.Equal(t, message, responseData.Message)
}
