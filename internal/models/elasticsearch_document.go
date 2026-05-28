package models

type LocalizedField struct {
	Locale string `json:"locale"`
	Data   string `json:"data"`
}

type BrandDocument struct {
	Code  string            `json:"code"`
	Label map[string]string `json:"label"`
}

type AttributeDocument struct {
	Code  string            `json:"code"`
	Label map[string]string `json:"label"`
}

type ProductDocument struct {
	UUID        string                       `json:"uuid"`
	SKU         string                       `json:"sku"`
	PartNumber  string                       `json:"part_number"`
	Brand       BrandDocument                `json:"brand"`
	ProductName []LocalizedField             `json:"productname"`
	Attributes  map[string]string            `json:"attributes"`
	OilGrade    map[string]interface{}       `json:"oil_grade,omitempty"`
}