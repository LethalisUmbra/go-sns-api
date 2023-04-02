package routes

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/lethalisumbra/go-sns-api/models"
	"github.com/lethalisumbra/go-sns-api/utils"

	"github.com/gin-gonic/gin"
)

// Obtener todos los productos
func GetProducts(c *gin.Context) {
	var product models.Product
	rows, err := utils.DB.Query("SELECT * FROM products")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.SKU, &product.MercadoID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// Obtener un producto por su ID
func GetProduct(c *gin.Context) {
	var product models.Product
	var err error

	product.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := utils.DB.QueryRow("SELECT name, description, price, stock, sku, mercado_id FROM products WHERE id = ?", product.ID)

	err = row.Scan(&product.Name, &product.Description, &product.Price, &product.Stock, &product.SKU, &product.MercadoID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, product)
}

// Crear un nuevo producto
func CreateProduct(c *gin.Context) {
	var product models.Product
	var err error
	err = c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := utils.DB.Exec("INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)", product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := result.LastInsertId()

	product.ID = int(id)
	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	var product models.Product
	var err error

	product.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// obtener el tipo y valor de la estructura de producto
	t := reflect.TypeOf(product)
	v := reflect.ValueOf(product)

	// construir la query de actualización
	var query strings.Builder
	query.WriteString("UPDATE products SET ")
	var args []interface{}
	var hasFields bool // variable de bandera para controlar la coma final
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if value != nil && !reflect.DeepEqual(value, reflect.Zero(field.Type).Interface()) { // manejar diferentes tipos de valores nulos
			if hasFields {
				query.WriteString(", ")
			}
			query.WriteString(fmt.Sprintf("%s = ?", field.Tag.Get("json")))
			args = append(args, value)
			hasFields = true
		}
	}
	query.WriteString(" WHERE id = ?;")
	args = append(args, product.ID)

	fmt.Println(args...)

	// ejecutar la query
	_, err = utils.DB.Exec(query.String(), args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	GetProduct(c)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := utils.DB.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
