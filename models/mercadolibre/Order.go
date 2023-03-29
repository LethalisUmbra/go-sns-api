package mercadolibre

import "time"

type Coupon struct {
	Amount float64     `json:"amount"`
	ID     interface{} `json:"id"`
}

type RequestedQuantity struct {
	Value   int    `json:"value"`
	Measure string `json:"measure"`
}

type Item struct {
	ID                  string        `json:"id"`
	Title               string        `json:"title"`
	CategoryID          string        `json:"category_id"`
	VariationID         interface{}   `json:"variation_id"`
	SellerCustomField   interface{}   `json:"seller_custom_field"`
	VariationAttributes []interface{} `json:"variation_attributes"`
	Warranty            interface{}   `json:"warranty"`
	Condition           string        `json:"condition"`
	SellerSKU           string        `json:"seller_sku"`
	GlobalPrice         interface{}   `json:"global_price"`
	NetWeight           interface{}   `json:"net_weight"`
}

type PaymentCollector struct {
	ID int `json:"id"`
}

type ATMTransferReference struct {
	TransactionID string `json:"transaction_id"`
}

type Payment struct {
	ID                   int                  `json:"id"`
	OrderID              int                  `json:"order_id"`
	PayerID              int                  `json:"payer_id"`
	Collector            PaymentCollector     `json:"collector"`
	CardID               int                  `json:"card_id"`
	Reason               string               `json:"reason"`
	SiteID               string               `json:"site_id"`
	PaymentMethodID      string               `json:"payment_method_id"`
	CurrencyID           string               `json:"currency_id"`
	Installments         int                  `json:"installments"`
	IssuerID             string               `json:"issuer_id"`
	ATMTransferReference ATMTransferReference `json:"atm_transfer_reference"`
	DateApproved         time.Time            `json:"date_approved"`
	DateLastModified     time.Time            `json:"date_last_modified"`
	MoneyReleaseDate     time.Time            `json:"money_release_date"`
	AvailableActions     []string             `json:"available_actions"`
	Status               string               `json:"status"`
	StatusDetail         string               `json:"status_detail"`
	AuthorizationCode    interface{}          `json:"authorization_code"`
	TransactionAmount    float64              `json:"transaction_amount"`
	ShippingCost         float64              `json:"shipping_cost"`
	OverpaidAmount       float64              `json:"overpaid_amount"`
	TotalPaidAmount      float64              `json:"total_paid_amount"`
	MarketplaceFee       float64              `json:"marketplace_fee"`
	TransactionNetAmount float64              `json:"transaction_net_amount"`
}

type OrderItem struct {
	Item              Item              `json:"item"`
	Quantity          int               `json:"quantity"`
	RequestedQuantity RequestedQuantity `json:"requested_quantity"`
	PickedQuantity    interface{}       `json:"picked_quantity"`
	UnitPrice         float64           `json:"unit_price"`
	FullUnitPrice     float64           `json:"full_unit_price"`
	CurrencyID        string            `json:"currency_id"`
	ManufacturingDays interface{}       `json:"manufacturing_days"`
	SaleFee           float64           `json:"sale_fee"`
	ListingTypeID     string            `json:"listing_type_id"`
	BaseExchangeRate  interface{}       `json:"base_exchange_rate"`
	BaseCurrencyID    interface{}       `json:"base_currency_id"`
	ElementID         interface{}       `json:"element_id"`
	Discounts         interface{}       `json:"discounts"`
	Bundle            interface{}       `json:"bundle"`
}

type Order struct {
	ID              int           `json:"id"`
	DateCreated     time.Time     `json:"date_created"`
	LastUpdated     time.Time     `json:"last_updated"`
	ExpirationDate  time.Time     `json:"expiration_date"`
	DateClosed      time.Time     `json:"date_closed"`
	Comment         interface{}   `json:"comment"`
	PackID          interface{}   `json:"pack_id"`
	PickupID        interface{}   `json:"pickup_id"`
	Fulfilled       interface{}   `json:"fulfilled"`
	HiddenForSeller interface{}   `json:"hidden_for_seller"`
	BuyingMode      string        `json:"buying_mode"`
	ShippingCost    interface{}   `json:"shipping_cost"`
	ApplicationID   interface{}   `json:"application_id"`
	Mediations      []interface{} `json:"mediations"`
	TotalAmount     float64       `json:"total_amount"`
	PaidAmount      float64       `json:"paid_amount"`
	Coupon          Coupon        `json:"coupon"`
	OrderItems      []OrderItem   `json:"order_items"`
	CurrencyID      string        `json:"currency_id"`
	Payments        []Payment     `json:"payments"`
}
