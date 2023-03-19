package routes

import (
	"main/models"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Obtener todos los productos
func GetProducts(c *gin.Context) {
	rows, err := utils.DB.Query("SELECT * FROM products")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// Obtener un producto por su ID
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	row := utils.DB.QueryRow("SELECT * FROM products WHERE id = ?", id)

	var p models.Product
	err := row.Scan(&p.Name, &p.Description, &p.Price)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, p)
}

// Crear un nuevo producto
func CreateProduct(c *gin.Context) {
	var p models.Product
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := utils.DB.Exec("INSERT INTO products (name, description, price) VALUES (?, ?, ?)", p.Name, p.Description, p.Price)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	id, _ := result.LastInsertId()

	p.ID = int(id)
	c.JSON(201, p)
}
