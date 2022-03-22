package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type CountryParityFactor struct {
	ID             uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ParityFactor   float64
	CountryFromISO string `gorm:"column:country_from_iso"`
	CountryToISO   string `gorm:"column:country_to_iso"`
	Year           int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Country struct {
	ID              uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CountryName     string
	ISO             string `gorm:"column:iso"`
	GDPPerCapitaPPP int    `gorm:"column:gdp_percapita_ppp"`
	Year            int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
