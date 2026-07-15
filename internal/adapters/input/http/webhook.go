package http

import (
	"context"
	"net/http"
	"time"

	"github.com/tegnoword/orienmod/internal/core/domain"
	"github.com/tegnoword/orienmod/internal/core/ports"
)

type WebhookHandler struct {
	useCase ports.ClassroomEventUseCase
}

func NewWebhookHandler(uc ports.ClassroomEventUseCase) *WebhookHandler {
	return &WebhookHandler{useCase: uc}
}

func (h *WebhookHandler) HandleGoogleClassroomWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	resourceID := r.Header.Get("X-Goog-Resource-ID")
	resourceState := r.Header.Get("X-Goog-Resource-State")

	if resourceState == "sync" {
		w.WriteHeader(http.StatusOK)
		return
	}

	event := domain.ClassroomEvent{
		ResourceID: resourceID,
		EventType:  resourceState,
		ReceivedAt: time.Now(),
	}

	go func(evt domain.ClassroomEvent) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = h.useCase.HandleClassroomNotification(ctx, evt)
	}(event)

	w.WriteHeader(http.StatusNoContent)
}
