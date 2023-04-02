package utils

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, statusCode int, err error, message string) {
	// Almacenar error en Base de Datos
	er := StoreError(statusCode, err, message)
	if er != nil {
		log.Printf("%s: %s", message, er.Error())
	}

	log.Printf("%s: %s", message, err.Error())
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}

func StoreError(statusCode int, err error, message string) error {
	_, er := DB.Exec("INSERT INTO errors (status, error, message, date) VALUES (?, ?, ?, ?)", statusCode, err.Error(), message, time.Now().Format("2006-01-02 15:04:05"))
	return er
}
