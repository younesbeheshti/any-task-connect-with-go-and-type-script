package infra

import (
	"time"

	"github.com/google/uuid"
)

// CityModel is the GORM persistence model for cities.
type CityModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string    `gorm:"size:100;not null"`
	Province  string    `gorm:"size:100;not null;default:''"`
	CreatedAt time.Time
}

func (CityModel) TableName() string { return "cities" }
