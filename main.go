package main

import (
	"log"
	"main/routes"
	"main/utils"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Crea una carpeta para los logs si no existe
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Crea un archivo de log en la carpeta "logs"
	file, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Abrir SQLITE
	err = utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	products := router.Group("/products")
	{
		products.GET("/", routes.GetProducts)
		products.GET("/:id", routes.GetProduct)
		products.POST("/", routes.CreateProduct)
		products.PATCH("/:id", routes.UpdateProduct)
		products.DELETE("/:id", routes.DeleteProduct)
	}

	mercado := router.Group("/mercadolibre")
	{
		mercado.GET("/", routes.GetMercadoCode)
		mercado.POST("/", routes.CreateMercadoProduct)
		mercado.GET("/:id", routes.GetMercadoProduct)
	}

	// Especifica el archivo de log en el middleware Gin Logger
	router.Use(gin.LoggerWithWriter(file))

	router.Run(":8080")
}
