package domain

type Grade struct {
	StudentID string  `json:"student_id"`
	TaskID    string  `json:"task_id"`
	Score     float64 `json:"score"`
	CourseID  string  `json:"course_id"`
}
