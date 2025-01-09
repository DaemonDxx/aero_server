package middleware

import (
	"encoding/json"
	"fmt"
	httpErr "github.com/daemondxx/lks_back/internal/server/http_errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/http"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		e := c.Errors.Last()
		var errHttp *httpErr.HttpError
		if errors.As(e, &errHttp) {
			c.JSON(errHttp.Status, gin.H{
				"message": errHttp.Message,
			})
			return
		} else {
			var ve validator.ValidationErrors
			var jErr *json.UnmarshalTypeError

			if errors.As(e, &ve) {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": createValidationErrMessage(ve),
				})
			} else if errors.As(e, &jErr) {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("parse json body error: field [%s]", jErr.Field),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal server error",
				})
			}
			return
		}
	}
}

func createValidationErrMessage(arr validator.ValidationErrors) string {
	message := "validation body error: \n\r"
	for _, err := range arr {
		message += fmt.Sprintf("	- [field: %v, value: %v] %e;\n\r", err.Field(), err.Param(), err)
	}
	return message
}
