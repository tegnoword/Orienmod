package domain

import "time"

type ClassroomEvent struct {
	ResourceID string    `json:"resource_id"`
	EventType  string    `json:"event_type"`
	ReceivedAt time.Time `json:"received_at"`
}
