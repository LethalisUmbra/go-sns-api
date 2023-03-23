package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/models"
	"main/utils"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ApiUrl string = "https://api.mercadolibre.com/"
var ClientId int = 572501217597774
var ClientSecret string = "Yk7f2sJjsNpHv1Ev9XVQUHYbdxLXzDaQ"

// Obtener todos los productos
func GetMercadoCode(c *gin.Context) {
	// Datos del formulario a enviar
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("client_id", strconv.Itoa(ClientId))
	values.Set("client_secret", ClientSecret)
	values.Set("code", c.Query("code"))
	values.Set("redirect_uri", "http://localhost:8080/mercadolibre")

	// Codificar los datos del formulario en la URL
	data := values.Encode()

	// Crear un buffer de bytes con los datos del formulario
	buffer := bytes.NewBufferString(data)

	// Crear el request HTTP POST
	req, err := http.NewRequest("POST", ApiUrl+"/oauth/token", buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"create_request_error": err.Error()})
		return
	}

	// Establecer el header Content-Type en application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Realizar la petición HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"http_request_error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// lee la respuesta del servidor
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"read_body_error": err.Error()})
		return
	}

	// Valida que el estado sea OK
	if resp.StatusCode != http.StatusOK {
		var mercadoError interface{}
		_ = json.Unmarshal(body, &mercadoError)
		c.JSON(resp.StatusCode, mercadoError)
		return
	}

	// Decodifica el JSON en la estructura de Token
	token := models.MercadoToken{CreatedAt: time.Now()}
	err = json.Unmarshal(body, &token)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"json_unmarshal_error": err.Error()})
		return
	}

	// Almacenar Token en la Base de Datos
	result, err := utils.DB.Exec("INSERT INTO token (access_token, token_type, expires_in, scope, user_id, refresh_token, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)", token.AccessToken, token.TokenType, token.ExpiresIn, token.Scope, token.UserID, token.RefreshToken, token.CreatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"database_error": err.Error()})
		return
	}
	id, _ := result.LastInsertId()

	token.ID = int(id)
	c.JSON(http.StatusCreated, token)
}

// Obtener un producto por su ID
func GetMercadoProduct(c *gin.Context) {
	var product models.MercadoProduct
	var err error

	product.ID = c.Param("id")

	// Obtener Auth Token
	token, err := utils.GetLastToken()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// crea un objeto request con el header
	req, err := http.NewRequest("GET", ApiUrl+"items/"+product.ID, nil)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Agregar AccessToken en Header
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// crea un cliente HTTP y realiza el GET request con el objeto request creado
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// lee la respuesta del servidor
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"error": err.Error()})
		return
	}

	// Valida que el estado sea OK
	if resp.StatusCode != http.StatusOK {
		var mercadoError interface{}
		_ = json.Unmarshal(body, &mercadoError)
		c.JSON(resp.StatusCode, mercadoError)
		return
	}

	// decodifica el JSON en la estructura de Product de MercadoLibre
	err = json.Unmarshal(body, &product)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, product)
}

// Crear un nuevo producto
func CreateMercadoProduct(c *gin.Context) {
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

func UpdateMercadoProduct(c *gin.Context) {
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

func DeleteMercadoProduct(c *gin.Context) {
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
