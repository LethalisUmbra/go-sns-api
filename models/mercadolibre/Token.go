package mercadolibre

import "time"

type MercadoToken struct {
	ID           int       `json:"id"`
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	Scope        string    `json:"scope"`
	UserID       int       `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
}
