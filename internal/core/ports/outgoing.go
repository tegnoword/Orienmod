package ports

import (
	"context"

	"github.com/tegnoword/orienmod/internal/core/domain"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification domain.Notification) error
	GetByTeacher(ctx context.Context, teacherID string) ([]domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
}

type ClassroomRepository interface {
	SaveEvent(ctx context.Context, event domain.ClassroomEvent) error
	SaveCourse(ctx context.Context, course domain.Course) error
	GetAllCourses(ctx context.Context) ([]domain.Course, error)
	SaveStudent(ctx context.Context, student domain.Studen) error
	GetStudentsByCourse(ctx context.Context, courseID string) ([]domain.Studen, error)
	SaveTasks(ctx context.Context, studentID string, tasks []domain.Task) error
	SaveGrade(ctx context.Context, grade domain.Grande) error
}
