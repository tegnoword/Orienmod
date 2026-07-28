package validation

import (
	"fmt"
	"time"

	"github.com/tegnoword/orienmod/internal/core/domain"
)

func ValidateCourse(course *domain.Course) error {
	if course.Name == "" {
		return fmt.Errorf("name es requerido")
	}
	if len(course.Name) > 250 {
		return fmt.Errorf("name no puede exceder 250 caracteres")
	}
	return nil
}

func ValidateStudent(student *domain.Student) error {
	if student.ID == "" {
		return fmt.Errorf("student_id es requerido")
	}
	if student.Email == "" {
		return fmt.Errorf("email es requerido")
	}
	if !contains(student.Email, "@") {
		return fmt.Errorf("email inválido")
	}
	return nil
}

func ValidateTask(task *domain.Task) error {
	if task.CourseID == "" {
		return fmt.Errorf("course_id es requerido")
	}
	if task.Title == "" {
		return fmt.Errorf("title es requerido")
	}
	if len(task.Title) > 250 {
		return fmt.Errorf("title no puede exceder 250 caracteres")
	}
	if task.MaxPoints < 0 {
		return fmt.Errorf("max_points no puede ser negativo")
	}
	if !task.DueDate.IsZero() && task.DueDate.Before(time.Now()) {
		return fmt.Errorf("due_date no puede ser en el pasado")
	}
	// Validar WorkType
	validWorkTypes := map[string]bool{
		"ASSIGNMENT":               true,
		"SHORT_ANSWER":             true,
		"MULTIPLE_CHOICE_QUESTION": true,
	}
	if task.WorkType != "" && !validWorkTypes[task.WorkType] {
		return fmt.Errorf("work_type inválido: %s", task.WorkType)
	}
	return nil
}

func ValidateGrade(grade *domain.Grade) error {
	if grade.StudentID == "" {
		return fmt.Errorf("student_id es requerido")
	}
	if grade.TaskID == "" {
		return fmt.Errorf("task_id es requerido")
	}
	if grade.Score < 0 || grade.Score > 100 {
		return fmt.Errorf("score debe estar entre 0 y 100")
	}
	return nil
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
