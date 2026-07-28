// internal/core/domain/task.go
package domain

import "time"

// Task representa una tarea en Google Classroom
type Task struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	MaxPoints   float64   `json:"max_points"`
	State       string    `json:"state"`     // DRAFT, PUBLISHED, DELETED
	WorkType    string    `json:"work_type"` // ✅ NUEVO: ASSIGNMENT, SHORT_ANSWER, MULTIPLE_CHOICE
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskSubmission representa la entrega de un estudiante
type TaskSubmission struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	TaskID      string    `json:"task_id"`
	StudentID   string    `json:"student_id"`
	StudentName string    `json:"student_name"`
	State       string    `json:"state"` // ASSIGNED, TURNED_IN, RETURNED
	SubmittedAt time.Time `json:"submitted_at"`
	Grade       float64   `json:"grade"`
}

// WorkType constants
const (
	WorkTypeAssignment     = "ASSIGNMENT"
	WorkTypeShortAnswer    = "SHORT_ANSWER"
	WorkTypeMultipleChoice = "MULTIPLE_CHOICE"
)
