package domain

import "time"

type NotificationType string

const (
	TypeNewSubmission  NotificationType = "nueva_entrega"
	TypeLateSubmission NotificationType = "entrega_retrasada"
	TypeClassroomSync  NotificationType = "sincronizacion_classroom"
)

type Notification struct {
	ID        string           `json:"id"`
	TeacherID string           `json:"teacher_id"`
	ClassID   string           `json:"class_id"`
	Message   string           `json:"message"`
	Type      NotificationType `json:"type"`
	CreatedAt time.Time        `json:"created_at"`
	IsRead    bool             `json:"is_read"`
}
