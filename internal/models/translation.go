package models

import "time"

type Translation struct {
	EntityType string
	EntityID   string
	Locale     string
	FieldName  string
	FieldValue string
	UpdatedAt  time.Time
}