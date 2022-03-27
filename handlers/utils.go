package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParam(ctx *gin.Context, paramName string) (uint, error) {
	id := ctx.Param(paramName)
	idInt, err := strconv.Atoi(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid parameter: %s", paramName)})
		ctx.Abort()
		return 0, err

	}

	return uint(idInt), nil
}
