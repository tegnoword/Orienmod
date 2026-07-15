package services

import (
	"context"
	"fmt"
	"time"

	"github.com/tegnoword/orienmod/internal/core/domain"
	"github.com/tegnoword/orienmod/internal/core/ports"
)

type NotificationService struct {
	repo ports.NotificationRepository
}

func NewNotificationService(repo ports.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) HandleClassroomNotification(ctx context.Context, event domain.ClassroomEvent) error {
	alertMessage := fmt.Sprintf("Se detectó una nueva actividad en Classroom para el recurso %s", event.ResourceID)

	alert := domain.Notification{
		ID:        generateUniqueID(),
		TeacherID: "maestro_actual",
		ClassID:   "clase_detectada",
		Message:   alertMessage,
		Type:      domain.TypeNewSubmission,
		CreatedAt: time.Now(),
		IsRead:    false,
	}

	err := s.repo.Save(ctx, alert)
	if err != nil {
		return fmt.Errorf("error al guardar la alerta en el core: %w", err)
	}

	return nil
}

func generateUniqueID() string {
	return fmt.Sprintf("notif_%d", time.Now().UnixNano())
}
