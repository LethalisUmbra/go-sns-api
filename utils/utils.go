package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, statusCode int, err error, message string) {
	log.Printf("%s: %s", message, err.Error())
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}
