package models

import "time"

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	SKU         string  `json:"sku"`
	MercadoID   string  `json:"mercado_id"`
}

// Define la estructura de un producto de MercadoLibre
type MercadoProduct struct {
	ID                        string        `json:"id"`
	SiteID                    string        `json:"site_id"`
	Title                     string        `json:"title" binding:"required"`
	Subtitle                  string        `json:"subtitle"`
	SellerID                  int64         `json:"seller_id"`
	CategoryID                string        `json:"category_id" binding:"required"`
	OfficialStoreID           int           `json:"official_store_id"`
	Price                     float64       `json:"price" binding:"required"`
	BasePrice                 float64       `json:"base_price"`
	OriginalPrice             float64       `json:"original_price"`
	InventoryId               string        `json:"inventory_id"`
	CurrencyID                string        `json:"currency_id" binding:"required"`
	InitialQuantity           int64         `json:"initial_quantity"`
	AvailableQuantity         int64         `json:"available_quantity"  binding:"required"`
	SoldQuantity              int64         `json:"sold_quantity"`
	SaleTerms                 []SaleTerm    `json:"sale_terms"`
	BuyingMode                string        `json:"buying_mode" binding:"required"`
	ListingTypeID             string        `json:"listing_type_id" binding:"required"`
	StartTime                 time.Time     `json:"start_time"`
	StopTime                  time.Time     `json:"stop_time"`
	EndTime                   time.Time     `json:"end_time"`
	ExpirationTime            time.Time     `json:"expiration_time"`
	Condition                 string        `json:"condition" binding:"required"`
	Permalink                 string        `json:"permalink"`
	ThumbnailID               string        `json:"thumbnail_id"`
	Thumbnail                 string        `json:"thumbnail"`
	SecureThumbnail           string        `json:"secure_thumbnail"`
	Pictures                  []Picture     `json:"pictures" binding:"required"`
	VideoID                   string        `json:"video_id"`
	Descriptions              []string      `json:"descriptions"`
	AcceptsMP                 bool          `json:"accepts_mercadopago"`
	PaymentMeth               []string      `json:"non_mercado_pago_payment_methods"`
	Shipping                  Shipping      `json:"shipping"`
	InternationalDeliveryMode string        `json:"international_delivery_mode"`
	SellerAddress             SellerAddress `json:"seller_address"`
	SellerContact             interface{}   `json:"seller_contact"`
	Location                  struct{}      `json:"location"`
	Geolocation               interface{}   `json:"geolocation"`
	CoverageAreas             []string      `json:"coverage_areas"`
	Attributes                []Attribute   `json:"attributes" binding:"required"`
	Warnings                  []struct{}    `json:"warnings"`
	ListingSource             string        `json:"listing_source"`
	Variations                []Variation   `json:"variations"`
	Status                    string        `json:"status"`
	SubStatus                 []interface{} `json:"sub_status"`
	Tags                      []string      `json:"tags"`
	Warranty                  struct{}      `json:"warranty"`
	CatalogProductID          interface{}   `json:"catalog_product_id"`
	DomainID                  string        `json:"domain_id"`
	SellerCustomField         interface{}   `json:"seller_custom_field"`
	ParentItemID              interface{}   `json:"parent_item_id"`
	DifferentialPricing       interface{}   `json:"differential_pricing"`
	DealIDs                   []interface{} `json:"deal_ids"`
	AutomaticRelist           bool          `json:"automatic_relist"`
	DateCreated               time.Time     `json:"date_created"`
	LastUpdated               time.Time     `json:"last_updated"`
	Health                    float64       `json:"health"`
	CatalogListing            bool          `json:"catalog_listing"`
	Channels                  []string      `json:"channels"`
}

type PostMercadoProduct struct {
	Title             string      `json:"title" binding:"required"`
	CategoryID        string      `json:"category_id" binding:"required"`
	Price             float64     `json:"price" binding:"required"`
	CurrencyID        string      `json:"currency_id" binding:"required"`
	AvailableQuantity int64       `json:"available_quantity"  binding:"required"`
	BuyingMode        string      `json:"buying_mode" binding:"required"`
	ListingTypeID     string      `json:"listing_type_id" binding:"required"`
	Condition         string      `json:"condition" binding:"required"`
	Pictures          interface{} `json:"pictures" binding:"required"`
	Attributes        interface{} `json:"attributes"`
}
