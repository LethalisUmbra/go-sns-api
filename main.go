package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/lethalisumbra/go-sns-api/routes"
	"github.com/lethalisumbra/go-sns-api/utils"
)

func main() {
	// Leer .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Abrir Base de Datos
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
		mercado.POST("/users", routes.CreateMercadoUser)
		mercado.POST("/notifications", routes.HandleMercadoCallback)
	}

	router.Run(":8080")
}
