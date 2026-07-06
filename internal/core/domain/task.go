package domain

import "time"

type Task struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Points       int       `json:"points"`
	DeliveryDate time.Time `json:"delivery_date"`
	CreationDate time.Time `json:"creation_date"`
}
