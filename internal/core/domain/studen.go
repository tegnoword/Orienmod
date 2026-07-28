package domain

type Student struct {
	ID       string `json:"id"`
	Name     string `json:"name"` // No 'FullName'
	Email    string `json:"email"`
	CourseID string `json:"course_id"` // Si existe en el dominio
}
