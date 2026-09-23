package entity

import "time"

// Config represents a key-value configuration item stored in the database
type Config struct {
	ID        int64     `db:"id" json:"id"`
	Key       string    `db:"key" json:"key"`
	Value     string    `db:"value" json:"value"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
