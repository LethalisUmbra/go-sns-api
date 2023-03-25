package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lethalisumbra/go-sns-api/models"
	"github.com/lethalisumbra/go-sns-api/utils"

	"github.com/gin-gonic/gin"
)

var ApiUrl string = "https://api.mercadolibre.com/"
var ClientId int = 572501217597774
var ClientSecret string = "Yk7f2sJjsNpHv1Ev9XVQUHYbdxLXzDaQ"

// Obtener todos los productos
func GetMercadoCode(c *gin.Context) {
	// Datos del formulario a enviar
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", strconv.Itoa(ClientId))
	data.Set("client_secret", ClientSecret)
	data.Set("code", c.Query("code"))
	data.Set("redirect_uri", "http://localhost:8080/mercadolibre")

	// Crear el request HTTP POST
	req, err := http.NewRequest("POST", ApiUrl+"/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "create_request_error")
		return
	}

	// Establecer el header Content-Type en application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Crear un cliente HTTP personalizado
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Realizar la petición HTTP
	resp, err := client.Do(req)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "http_request_error")
		return
	}
	defer resp.Body.Close()

	// lee la respuesta del servidor
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.HandleError(c, resp.StatusCode, err, "read_body_error")
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
		utils.HandleError(c, resp.StatusCode, err, "json_unmarshal_error")
		return
	}

	// Almacenar Token en la Base de Datos
	result, err := utils.DB.Exec("INSERT INTO token (access_token, token_type, expires_in, scope, user_id, refresh_token, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)", token.AccessToken, token.TokenType, token.ExpiresIn, token.Scope, token.UserID, token.RefreshToken, token.CreatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "database_error")
		return
	}
	id, _ := result.LastInsertId()

	token.ID = int(id)
	c.JSON(http.StatusCreated, token)
}

// Obtener un producto por su ID
func GetMercadoProduct(c *gin.Context) {
	product := models.MercadoProduct{ID: c.Param("id")}

	// Obtener Auth Token
	token, err := utils.GetLastToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Crea un objeto request con el header
	req, err := http.NewRequest("GET", ApiUrl+"items/"+product.ID, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Agregar AccessToken en Header
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Crea un cliente HTTP y realiza el GET request con el objeto request creado
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Valida que el estado sea OK
	if resp.StatusCode != http.StatusOK {
		var mercadoError interface{}
		_ = json.NewDecoder(resp.Body).Decode(&mercadoError)
		c.AbortWithStatusJSON(resp.StatusCode, mercadoError)
		return
	}

	// Decodifica el JSON en la estructura de Product de MercadoLibre
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Crear un nuevo producto
func CreateMercadoProduct(c *gin.Context) {
	var bodyProduct models.PostMercadoProduct
	err := c.BindJSON(&bodyProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := utils.CreateMercadoProduct(bodyProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

func HandleMercadoCallback(c *gin.Context) {
	var callback models.MercadoCallback
	err := c.BindJSON(&callback)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := utils.DB.Prepare("INSERT INTO callbacks (_id, resource, user_id, topic, application_id, attempts, sent, received) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	sent := callback.Sent.Format("2006-01-02 15:04:05")
	received := callback.Received.Format("2006-01-02 15:04:05")

	go func() {
		result, err := stmt.Exec(callback.MercadoID, callback.Resource, callback.UserID, callback.Topic, callback.ApplicationID, callback.Attempts, sent, received)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		callback.ID = int(id)
	}()

	c.Status(http.StatusOK)
}

func CreateMercadoUser(c *gin.Context) {
	user, err := utils.CreateMercadoUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = utils.DB.Exec("INSERT INTO users (id, email, nickname, site_status, password) VALUES (?, ?, ?, ?, ?);", user.ID, user.Email, user.Nickname, user.SiteStatus, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
