package mercadolibre

type User struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Nickname   string `json:"nickname"`
	SiteStatus string `json:"site_status"`
	Password   string `json:"password"`
}
