package models

import "time"

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
