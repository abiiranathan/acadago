package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type validationError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func RespondWithError(c *gin.Context, err error) {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make([]validationError, len(ve))
		for i, fe := range ve {
			out[i] = validationError{fe.Field(), msgForTag(fe)}
		}

		c.JSON(http.StatusBadRequest, gin.H{"errors": out})
		c.Abort()
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	c.Abort()

}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "dir":
		return "Not a valid directory"
	case "file":
		return "Not a valid file path"
	case "max":
		return "Too long"
	case "min":
		return "Too short"
	case "unique":
		return "This field must be unique"
	case "uuid4":
		return "Invalid uuid4"
	case "uuid":
		return "Invalid uuid"
	case "jwt":
		return "Invalid jwt"
	case "json":
		return "Invalid json"
	case "datetime":
		return "Invalid datetime format"
	case "uppercase":
		return "Must be uppercase"
	case "lowercase":
		return "Must be lowercase"
	case "boolean":
		return "Must be a boolean"
	case "alphanum":
		return "Must be a alpha-numeric"
	case "alpha":
		return "Must be alpha characters only"
	case "ascii":
		return "Must be alpha ascii only"
	case "numeric":
		return "Must be a number"
	case "iscolor":
	case "base64":
		return "Must be a valid base64 string"
	case "url":
		return "Must be a valid url"
	case "ip":
		return "Must be a valid ip"
	case "mac":
		return "Must be a valid mac address"
	default:
		return "Must be a a valid 'hexcolor|rgb|rgba|hsl|hsla'"
	}

	return fe.Error()
}
