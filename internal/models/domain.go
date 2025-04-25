package models

import (
	"time"

	"gorm.io/gorm"
)

// Domain represents main table.
type DomainInfo struct {
	Domain      string `gorm:"primaryKey;size:255"`
	Created     *time.Time
	Registrar   string       `gorm:"size:255"`
	AbuseEmail  string       `gorm:"size:255"`
	AbusePhone  string       `gorm:"size:50"`
	Nameservers []Nameserver `gorm:"foreignKey:Domain;constraint:OnDelete:CASCADE;"`
	Source      string       `gorm:"size:5"`
}

// Nameserver is the child table with one row per NS.
type Nameserver struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Domain     string `gorm:"index;size:255"` // FK back to Domain.Domain
	Nameserver string `gorm:"size:255;index"` // the actual NS string
	CreatedAt  time.Time
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&DomainInfo{}, &Nameserver{})
}
