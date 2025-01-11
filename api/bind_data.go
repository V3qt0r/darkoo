package api

import (
	"log"
	"net/http"
	"strings"

	"darkoo/apperrors"

	"github.com/gin-gonic/gin"
)


type Request interface {
	Validate() error
}

// bindData is helper function, returns false if data is not bound
func BindData(c *gin.Context, req Request) bool {
	// Bind incoming json to struct and check for validation errors

	if err := c.ShouldBindJSON(req); err != nil {
		log.Printf("Error deserializing json data: %v\n", err)
		c.JSON(http.StatusBadRequest, NewResponse(http.StatusBadRequest, err.Error(), nil))
		return false
	}

	log.Printf("Successfully deserialized json data: %v\n", req)

	if err := req.Validate(); err != nil {
		errors := strings.Split(err.Error(), ";")
		fErrors := make([]apperrors.FieldError, 0)

		for _, e := range errors {
			split := strings.Split(e, ":")
			er := apperrors.FieldError{
				Field:		strings.TrimSpace(split[0]),
				Message:	strings.TrimSpace(split[1]),
			}
			fErrors = append(fErrors, er)
		}

		c.JSON(apperrors.Status(err), gin.H{
			"code":			http.StatusBadRequest,
			"message":      "Bad request",
			"details":      fErrors,
		})
		return false
	}
	return true
}