package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lethalisumbra/go-sns-api/models/mercadolibre"
)

var ApiUrl string = "https://api.mercadolibre.com/"
var ClientId int = 572501217597774
var ClientSecret string = "Yk7f2sJjsNpHv1Ev9XVQUHYbdxLXzDaQ"
var httpClient = &http.Client{}

func GetLastToken() (mercadolibre.MercadoToken, error) {
	var token mercadolibre.MercadoToken
	var err error

	row := DB.QueryRow("SELECT * FROM token ORDER BY created_at DESC LIMIT 1;")

	err = row.Scan(&token.ID, &token.AccessToken, &token.TokenType, &token.ExpiresIn, &token.Scope, &token.UserID, &token.RefreshToken, &token.CreatedAt)
	if err != nil {
		return mercadolibre.MercadoToken{}, err
	}

	// Validar expiración de Token
	if token.CreatedAt.Add(time.Duration(token.ExpiresIn)).Before(time.Now()) {
		// Refrescar Token
		token, err = RefreshToken(token.RefreshToken)
		if err != nil {
			return mercadolibre.MercadoToken{}, err
		}
	}

	return token, nil
}

func RefreshToken(refreshToken string) (mercadolibre.MercadoToken, error) {
	// Datos del formulario a enviar
	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("client_id", strconv.Itoa(ClientId))
	values.Set("client_secret", ClientSecret)
	values.Set("refresh_token", refreshToken)

	// Codificar los datos del formulario en la URL
	data := values.Encode()

	// Crear un buffer de bytes con los datos del formulario
	buffer := bytes.NewBufferString(data)

	// Crear el request HTTP POST
	req, err := http.NewRequest("POST", ApiUrl+"/oauth/token", buffer)
	if err != nil {
		return mercadolibre.MercadoToken{}, err
	}

	// Establecer el header Content-Type en application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Realizar la petición HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return mercadolibre.MercadoToken{}, err
	}
	defer resp.Body.Close()

	// lee la respuesta del servidor
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mercadolibre.MercadoToken{}, err
	}

	// Valida que el estado sea OK
	if resp.StatusCode != http.StatusOK {
		var mercadoError interface{}
		_ = json.Unmarshal(body, &mercadoError)
		fmt.Println(mercadoError)
		return mercadolibre.MercadoToken{}, errors.New("no se ha podido refrescar el token")
	}

	// Decodifica el JSON en la estructura de Token
	token := mercadolibre.MercadoToken{CreatedAt: time.Now()}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return mercadolibre.MercadoToken{}, err
	}

	// Almacenar Token en la Base de Datos
	result, err := DB.Exec("INSERT INTO token (access_token, token_type, expires_in, scope, user_id, refresh_token, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)", token.AccessToken, token.TokenType, token.ExpiresIn, token.Scope, token.UserID, token.RefreshToken, token.CreatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return mercadolibre.MercadoToken{}, err
	}
	id, _ := result.LastInsertId()

	token.ID = int(id)
	return token, nil
}

func CreateMercadoProduct(formProduct mercadolibre.PostMercadoProduct) (mercadolibre.MercadoProduct, error) {
	// Convertir product struct en JSON utilizando json.Marshal:
	jsonBody, err := json.Marshal(formProduct)
	if err != nil {
		return mercadolibre.MercadoProduct{}, fmt.Errorf("no se pudo codificar el cuerpo de la solicitud como JSON: %w", err)
	}

	// Obtener Auth Token
	token, err := GetLastToken()
	if err != nil {
		return mercadolibre.MercadoProduct{}, fmt.Errorf("no se pudo obtener el token: %w", err)
	}

	// Crear la solicitud HTTP POST
	req, err := http.NewRequest("POST", ApiUrl+"items", bytes.NewBuffer(jsonBody))
	if err != nil {
		return mercadolibre.MercadoProduct{}, fmt.Errorf("no se pudo crear la solicitud HTTP: %w", err)
	}

	// Establecer el header Authorization
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Realizar la petición HTTP utilizando el cliente HTTP global
	resp, err := httpClient.Do(req)
	if err != nil {
		return mercadolibre.MercadoProduct{}, fmt.Errorf("no se pudo realizar la solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	// leer la respuesta del servidor
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mercadolibre.MercadoProduct{}, fmt.Errorf("no se pudo leer la respuesta del servidor: %w", err)
	}

	// Valida que el estado sea OK
	if resp.StatusCode >= 400 {
		var mercadoError interface{}
		_ = json.Unmarshal(body, &mercadoError)
		return mercadolibre.MercadoProduct{}, fmt.Errorf("el servidor devolvió un estado inesperado (%d): %v", resp.StatusCode, mercadoError)
	}

	// Decodifica el JSON en la estructura del Producto
	var product mercadolibre.MercadoProduct
	err = json.Unmarshal(body, &product)
	if err != nil {
		return mercadolibre.MercadoProduct{}, fmt.Errorf("no se pudo decodificar la respuesta del servidor: %w", err)
	}

	return product, nil
}

func CreateMercadoUser() (mercadolibre.User, error) {
	// Crear el objeto JSON
	data := map[string]string{
		"site_id": "MLC",
	}
	// Convertir product struct en JSON utilizando json.Marshal:
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return mercadolibre.User{}, fmt.Errorf("no se pudo codificar el cuerpo de la solicitud como JSON: %w", err)
	}

	// Obtener Auth Token
	token, err := GetLastToken()
	if err != nil {
		return mercadolibre.User{}, fmt.Errorf("no se pudo obtener el token: %w", err)
	}

	// Crear la solicitud HTTP POST
	req, err := http.NewRequest("POST", ApiUrl+"users/test_user", bytes.NewBuffer(jsonBody))
	if err != nil {
		return mercadolibre.User{}, fmt.Errorf("no se pudo crear la solicitud HTTP: %w", err)
	}

	// Establecer el header Authorization
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-type", "application/json")

	// Realizar la petición HTTP utilizando el cliente HTTP global
	resp, err := httpClient.Do(req)
	if err != nil {
		return mercadolibre.User{}, fmt.Errorf("no se pudo realizar la solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	// leer la respuesta del servidor
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mercadolibre.User{}, fmt.Errorf("no se pudo leer la respuesta del servidor: %w", err)
	}

	// Valida que el estado sea OK
	if resp.StatusCode >= 400 {
		var mercadoError interface{}
		_ = json.Unmarshal(body, &mercadoError)
		return mercadolibre.User{}, fmt.Errorf("el servidor devolvió un estado inesperado (%d): %v", resp.StatusCode, mercadoError)
	}

	// Decodifica el JSON en la estructura del Producto
	var user mercadolibre.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return mercadolibre.User{}, fmt.Errorf("no se pudo decodificar la respuesta del servidor: %w", err)
	}

	return user, nil
}

func StoreMLCallback(c *gin.Context) {
	var callback mercadolibre.MercadoCallback
	err := c.BindJSON(&callback)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err, "Error binding callback")
		return
	}

	stmt, err := DB.Prepare("INSERT INTO callbacks (_id, resource, user_id, topic, application_id, attempts, sent, received) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err, "Error preparing query to insert callback")
		return
	}
	defer stmt.Close()

	sent := callback.Sent.Format("2006-01-02 15:04:05")
	received := callback.Received.Format("2006-01-02 15:04:05")

	_, err = stmt.Exec(callback.MercadoID, callback.Resource, callback.UserID, callback.Topic, strconv.Itoa(int(callback.ApplicationID)), callback.Attempts, sent, received)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err, fmt.Sprintf("Error storing callback in database: application_id (int): %d", callback.ApplicationID))
		return
	}

	// Identificar 'orders' en callback
	if callback.Topic != "orders_v2" {
		return
	}

	// Hacer query con el resource a mercadolibre

	// Obtener Auth Token
	token, err := GetLastToken()
	if err != nil {
		StoreError(http.StatusInternalServerError, err, "Error consiguiendo el último token")
		return
	}

	// Crear la solicitud HTTP GET
	req, err := http.NewRequest("GET", ApiUrl+callback.Resource, nil)
	if err != nil {
		StoreError(http.StatusInternalServerError, err, "No se pudo crear la solicitud HTTP")
		return
	}

	// Establecer el header Authorization
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Realizar la petición HTTP utilizando el cliente HTTP global
	resp, err := httpClient.Do(req)
	if err != nil {
		StoreError(http.StatusInternalServerError, err, "No se pudo realizar la solicitud HTTP")
		return
	}
	defer resp.Body.Close()

	// leer la respuesta del servidor
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		StoreError(http.StatusInternalServerError, err, "No se pudo leer la respuesta del servidor")
		return
	}

	// Valida que el estado sea OK
	if resp.StatusCode >= 400 {
		var mercadoError interface{}
		_ = json.Unmarshal(body, &mercadoError)
		c.JSON(resp.StatusCode, mercadoError)
		return
	}

	// Decodifica el JSON en la estructura de la Order
	var order mercadolibre.Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		StoreError(http.StatusInternalServerError, err, "No se pudo decodificar la respuesta del servidor")
		return
	}

	// Almacenar orders
	fields := []string{
		"date_created",
		"last_updated",
		"expiration_date",
		"date_closed",
		"buying_mode",
		"total_amount",
		"paid_amount",
		"coupon_id",
		"currency_id",
		"_id",
	}

	values := []string{
		"'" + order.DateCreated.Format("2006-01-02 15:04:05") + "'",
		"'" + order.LastUpdated.Format("2006-01-02 15:04:05") + "'",
		"'" + order.ExpirationDate.Format("2006-01-02 15:04:05") + "'",
		"'" + order.DateClosed.Format("2006-01-02 15:04:05") + "'",
		"'" + order.BuyingMode + "'",
		fmt.Sprintf("%.2f", order.TotalAmount),
		fmt.Sprintf("%.2f", order.PaidAmount),
		"'" + fmt.Sprintf("%v", order.Coupon.ID) + "'",
		"'" + order.CurrencyID + "'",
		"'" + fmt.Sprintf("%v", order.ID) + "'",
	}

	query := fmt.Sprintf("INSERT INTO orders (%s) VALUES (%s);", strings.Join(fields, ","), strings.Join(values, ","))

	// Ejecutar query
	_, err = DB.Exec(query)

	// Validar errores
	if err != nil {
		StoreError(http.StatusInternalServerError, err, "No se pudo almacenar la orden en la base de datos")
		return
	}

	// Identificar producto(s)
	// For loop en order.OrderItems
	for _, orderItem := range order.OrderItems {
		// 1 - Por cada item hacer un update a la base de datos para actualizar el stock
		// 1.1 - UPDATE products SET stock = stock - quantity WHERE sku/mercado_id = order_item.seller_sku/order_item.item_id;
		go updateStock(orderItem.Item.SellerSKU, orderItem.Quantity)
	}

}

func updateStock(sku string, quantity int) {
	query := fmt.Sprintf("UPDATE products SET stock = stock - %d WHERE sku = %s;", quantity, sku)

	res, err := DB.Exec(query)
	// Validar errores
	if err != nil {
		StoreError(http.StatusInternalServerError, err, "No se pudo actualizar el stock del producto")
		return
	}

	if rows, err := res.RowsAffected(); err != nil || rows == 0 {
		StoreError(http.StatusInternalServerError, err, "No se pudo actualizar el stock del producto")
		return
	}
}
