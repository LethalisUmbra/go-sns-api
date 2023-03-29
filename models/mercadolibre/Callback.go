package mercadolibre

import "time"

type MercadoCallback struct {
	ID            int       `json:"id"`
	MercadoID     string    `json:"_id"`
	Resource      string    `json:"resource"`
	UserID        int       `json:"user_id"`
	Topic         string    `json:"topic"`
	ApplicationID int       `json:"application_id"`
	Attempts      int       `json:"attempts"`
	Sent          time.Time `json:"sent"`
	Received      time.Time `json:"received"`
}
