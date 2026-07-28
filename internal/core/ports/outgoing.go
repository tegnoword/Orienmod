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
	SaveEvent(ctx context.Context, clientEmail string, event domain.ClassroomEvent) error

	SaveCourse(ctx context.Context, clientEmail string, course domain.Course) error
	GetAllCourses(ctx context.Context, clientEmail string) ([]domain.Course, error)

	SaveStudent(ctx context.Context, clientEmail string, student domain.Student) error
	GetStudentsByCourse(ctx context.Context, clientEmail string, courseID string) ([]domain.Student, error)
	SaveTasks(ctx context.Context, clientEmail string, studentID string, tasks []domain.Task) error
	SaveGrade(ctx context.Context, clientEmail string, grade domain.Grade) error
}
