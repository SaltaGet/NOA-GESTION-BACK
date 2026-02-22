package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Invoice struct {
	ID int64 `gorm:"primaryKey" json:"id"`
    InvoiceData InvoiceData `gorm:"type:json" json:"invoice_data"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type InvoiceData map[string]interface{}

func (j *InvoiceData) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    return json.Unmarshal(bytes, j)
}

func (j InvoiceData) Value() (driver.Value, error) {
    return json.Marshal(j)
}
