package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"main/models"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var ApiUrl string = "https://api.mercadolibre.com/"
var ClientId int = 572501217597774
var ClientSecret string = "Yk7f2sJjsNpHv1Ev9XVQUHYbdxLXzDaQ"

func GetLastToken() (models.MercadoToken, error) {
	var token models.MercadoToken
	var err error

	row := DB.QueryRow("SELECT * FROM token ORDER BY created_at DESC LIMIT 1;")

	err = row.Scan(&token.ID, &token.AccessToken, &token.TokenType, &token.ExpiresIn, &token.Scope, &token.UserID, &token.RefreshToken, &token.CreatedAt)
	if err != nil {
		return models.MercadoToken{}, err
	}

	// Validar expiración de Token
	if token.CreatedAt.Add(time.Duration(token.ExpiresIn)).Before(time.Now()) {
		// Refrescar Token
		token, err = RefreshToken(token.RefreshToken)
		if err != nil {
			return models.MercadoToken{}, err
		}
	}

	return token, nil
}

func RefreshToken(refreshToken string) (models.MercadoToken, error) {
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
		return models.MercadoToken{}, err
	}

	// Establecer el header Content-Type en application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Realizar la petición HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.MercadoToken{}, err
	}
	defer resp.Body.Close()

	// lee la respuesta del servidor
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.MercadoToken{}, err
	}

	// Valida que el estado sea OK
	if resp.StatusCode != http.StatusOK {
		var mercadoError interface{}
		_ = json.Unmarshal(body, &mercadoError)
		fmt.Println(mercadoError)
		return models.MercadoToken{}, errors.New("No se ha podido refrescar el token")
	}

	// Decodifica el JSON en la estructura de Token
	token := models.MercadoToken{CreatedAt: time.Now()}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return models.MercadoToken{}, err
	}

	// Almacenar Token en la Base de Datos
	result, err := DB.Exec("INSERT INTO token (access_token, token_type, expires_in, scope, user_id, refresh_token, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)", token.AccessToken, token.TokenType, token.ExpiresIn, token.Scope, token.UserID, token.RefreshToken, token.CreatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return models.MercadoToken{}, err
	}
	id, _ := result.LastInsertId()

	token.ID = int(id)
	return token, nil
}
