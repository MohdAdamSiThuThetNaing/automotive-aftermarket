package models

type LocalizedField struct {
	Locale string `json:"locale"`
	Data   string `json:"data"`
}

type BrandDocument struct {
	Code  string            `json:"code"`
	Label map[string]string `json:"label"`
}

type AttributeValueDocument struct {
	Code  string            `json:"code"`
	Label map[string]string `json:"label"`
}

type ProductDocument struct {
	UUID       string `json:"uuid"`
	SKU        string `json:"sku"`
	PartNumber string `json:"part_number"`
	Brand BrandDocument `json:"brand"`
	ProductName []LocalizedField `json:"productname"`
	Attributes map[string]string `json:"attributes"`
	AttributeDetails map[string]AttributeValueDocument `json:"attribute_details,omitempty"`
}