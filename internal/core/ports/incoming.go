package ports

import (
	"context"

	"github.com/tegnoword/orienmod/internal/core/domain"
)

type ClassroomEventUseCase interface {
	HandleClassroomNotification(ctx context.Context, event domain.ClassroomEvent) error
}
