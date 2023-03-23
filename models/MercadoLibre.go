package models

import "time"

// Define la estructura de un producto de MercadoLibre
type MercadoProduct struct {
	ID                        string        `json:"id"`
	SiteID                    string        `json:"site_id"`
	Title                     string        `json:"title" binding:"reuired"`
	Subtitle                  string        `json:"subtitle"`
	SellerID                  int64         `json:"seller_id"`
	CategoryID                string        `json:"category_id" binding:"required"`
	OfficialStoreID           string        `json:"official_store_id"`
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

type Variation struct {
	Id                    int         `json:"id"`
	Price                 float64     `json:"price"`
	AttributeCombinations []Attribute `json:"attribute_combinations"`
	AvailableQuantity     int64       `json:"available_quantity"`
	SoldQuantity          int64       `json:"sold_quantity"`
	SaleTerms             []SaleTerm  `json:"sale_terms"`
	PictureIds            []string    `json:"picture_ids"`
	CatalogProductId      interface{} `json:"catalog_product_id"`
}

type Attribute struct {
	ID                 string           `json:"id"`
	Name               string           `json:"name"`
	ValueID            string           `json:"value_id"`
	ValueName          string           `json:"value_name"`
	ValueStruct        interface{}      `json:"value_struct"`
	Values             []AttributeValue `json:"values"`
	AttributeGroupID   string           `json:"attribute_group_id"`
	AttributeGroupName string           `json:"attribute_group_name"`
	ValueType          string           `json:"value_type"`
}

type AttributeValue struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Struct interface{} `json:"struct"`
}

type SaleTerm struct {
}

type Picture struct {
	ID        string `json:"id"`
	Url       string `json:"url"`
	SecureUrl string `json:"secure_url"`
	Size      string `json:"size"`
	MaxSize   string `json:"max_size"`
	Quality   string `json:"quality"`
}

type Shipping struct {
	Mode         string      `json:"mode"`
	Methods      []string    `json:"methods"`
	Tags         []string    `json:"tags"`
	Dimensions   interface{} `json:"dimensions"`
	LocalPickup  bool        `json:"local_pick_up"`
	FreeShipping bool        `json:"free_shipping"`
	LogisticType string      `json:"logistic_type"`
	StorePickup  bool        `json:"store_pick_up"`
}

type SellerAddress struct {
	ID             int         `json:"id"`
	Comment        string      `json:"comment"`
	AddressLine    string      `json:"address_line"`
	City           interface{} `json:"city"`
	State          interface{} `json:"state"`
	Country        interface{} `json:"country"`
	SearchLocation interface{} `json:"search_location"`
	Latitude       float64     `json:"latitude"`
	Longitude      float64     `json:"longitude"`
}

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
