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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lethalisumbra/go-sns-api/models/mercadolibre"
)

var ApiUrl string = "https://api.mercadolibre.com/"
var ClientId int = 572501217597774
var ClientSecret string = "Yk7f2sJjsNpHv1Ev9XVQUHYbdxLXzDaQ"
var httpClient = &http.Client{}

// Cache para token
var MercadoToken mercadolibre.MercadoToken

func GetLastToken() (mercadolibre.MercadoToken, error) {
	// Si el token está en caché y aún no ha expirado, lo devolvemos
	if MercadoToken.AccessToken != "" {
		if MercadoToken.CreatedAt.Add(time.Duration(MercadoToken.ExpiresIn)).After(time.Now()) {
			fmt.Println("El token en cache no está expirado")
			return MercadoToken, nil
		}
		fmt.Println("El token en cache sí está expirado")
	}

	var token mercadolibre.MercadoToken
	row := DB.QueryRow("SELECT * FROM token ORDER BY created_at DESC LIMIT 1;")
	err := row.Scan(&token.ID, &token.AccessToken, &token.TokenType, &token.ExpiresIn, &token.Scope, &token.UserID, &token.RefreshToken, &token.CreatedAt)
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

	// Almacenar el token en cache
	MercadoToken = token

	return MercadoToken, nil
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

func StoreMercadoCallback(callback mercadolibre.MercadoCallback) error {
	stmt, err := DB.Prepare("INSERT INTO callbacks (_id, resource, user_id, topic, application_id, attempts, sent, received) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	sent := callback.Sent.Format("2006-01-02 15:04:05")
	received := callback.Received.Format("2006-01-02 15:04:05")

	_, err = stmt.Exec(callback.MercadoID, callback.Resource, callback.UserID, callback.Topic, strconv.Itoa(int(callback.ApplicationID)), callback.Attempts, sent, received)
	if err != nil {
		return err
	}

	return nil
}

func GetMercadoOrder(resource string) (mercadolibre.Order, string, error) {
	// Obtener Auth Token
	token, err := GetLastToken()
	if err != nil {
		return mercadolibre.Order{}, "Error consiguiendo el último token", err
	}

	// Crear la solicitud HTTP GET
	req, err := http.NewRequest("GET", ApiUrl+resource, nil)
	if err != nil {
		return mercadolibre.Order{}, "No se pudo crear la solicitud HTTP", err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Realizar la petición HTTP utilizando el cliente HTTP global
	resp, err := httpClient.Do(req)
	if err != nil {
		return mercadolibre.Order{}, "No se pudo realizar la solicitud HTTP", err
	}
	defer resp.Body.Close()

	// leer la respuesta del servidor
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mercadolibre.Order{}, "No se pudo leer la respuesta del servidor", err
	}

	// Valida que el estado sea OK
	if resp.StatusCode >= 400 {
		return mercadolibre.Order{}, "MercadoLibre ha retornado un error", err
	}

	// Decodifica el JSON en la estructura de la Order
	var order mercadolibre.Order
	if err = json.Unmarshal(body, &order); err != nil {
		return mercadolibre.Order{}, "No se pudo decodificar la respuesta del servidor", err
	}

	return order, "", nil
}

func CreateOrder(c *gin.Context, resource string) {
	order, msg, err := GetMercadoOrder(resource)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err, msg)
		return
	}

	// Almacenar orders
	query := `INSERT INTO orders (date_created, last_updated, expiration_date, date_closed, buying_mode, total_amount, paid_amount, coupon_id, currency_id, _id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	if _, err := DB.Exec(query, order.DateCreated, order.LastUpdated, order.ExpirationDate, order.DateClosed, order.BuyingMode, order.TotalAmount, order.PaidAmount, order.Coupon.ID, order.CurrencyID, order.ID); err != nil {
		HandleError(c, http.StatusInternalServerError, err, "No se pudo almacenar la orden en la base de datos")
		return
	}

	// Identificar producto(s)
	for _, orderItem := range order.OrderItems {
		go updateStock(orderItem.Item.SellerSKU, orderItem.Quantity)
		go StoreOrderItem(orderItem, order.ID)
	}

	c.Status(http.StatusOK)
}

func updateStock(sku string, quantity int) {
	query := fmt.Sprintf("UPDATE products SET stock = stock - %d WHERE sku = '%s';", quantity, sku)

	res, err := DB.Exec(query)
	// Validar errores
	if err != nil {
		StoreError(http.StatusInternalServerError, err, "No se pudo actualizar el stock del producto, query: "+query)
		return
	}

	if rows, err := res.RowsAffected(); err != nil || rows == 0 {
		StoreError(http.StatusInternalServerError, err, "No se pudo actualizar el stock del producto")
		return
	}
}

func StoreOrderItem(orderItem mercadolibre.OrderItem, order_id uint64) error {
	// Almacenar orders
	query := `"INSERT INTO order_items (order_id, mercado_id, title, category_id, variation_id, seller_custom_field, item_condition, seller_sku, global_price, net_weight, quantity, requested_quantity_value, requested_quantity_measure, picked_quantity, unit_price, full_unit_price, currency_id, manufacturing_days, sale_fee, listing_type_id, base_exchange_rate, base_currency_id, element_id, discounts, bundle) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"`

	if _, err := DB.Exec(query, order_id, orderItem.Item.ID, orderItem.Item.Title, orderItem.Item.CategoryID, orderItem.Item.VariationID, orderItem.Item.SellerCustomField, orderItem.Item.Condition, orderItem.Item.SellerSKU, orderItem.Item.GlobalPrice, orderItem.Item.NetWeight, orderItem.Quantity, orderItem.RequestedQuantity.Value, orderItem.RequestedQuantity.Measure, orderItem.PickedQuantity, orderItem.UnitPrice, orderItem.FullUnitPrice, orderItem.CurrencyID, orderItem.ManufacturingDays, orderItem.SaleFee, orderItem.ListingTypeID, orderItem.BaseExchangeRate, orderItem.CurrencyID, orderItem.ElementID, orderItem.Discounts, orderItem.Bundle); err != nil {
		go StoreError(http.StatusInternalServerError, err, "Error almacenando el Order Item")
		return err
	}

	return nil
}
