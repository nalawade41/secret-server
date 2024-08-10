package responses

import (
	"github.com/labstack/echo/v4"
)

// Error represents the error for UI
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response transforms data for the UI with data
func Response(c echo.Context, statusCode int, data interface{}) error {
	// If needed, we can set the headers here
	acceptHeader := c.Request().Header.Get(echo.HeaderAccept)

	switch {
	case acceptHeader == "application/xml":
		// // c.XML automatically sets the Content-Type to application/json
		return c.XML(statusCode, data)

	case acceptHeader == "application/json", acceptHeader == "":
		fallthrough
	default:
		// c.JSON automatically sets the Content-Type to application/json
		return c.JSON(statusCode, data)
	}
}

// ErrorResponseWithMessage transforms data for the UI with error message
func ErrorResponseWithMessage(c echo.Context, statusCode int, message string) error {
	return Response(c, statusCode, Error{
		Code:    statusCode,
		Message: message,
	})
}
