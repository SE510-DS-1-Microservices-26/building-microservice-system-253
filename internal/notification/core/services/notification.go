package services

import (
	"cafeteria-delivery/internal/notification/core/domain"
	"cafeteria-delivery/internal/notification/repositories"
	"context"
)

type NotificationService struct {
	repo *repositories.NotificationRepository
}

func NewNotificationService(repository *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repository}
}

// Store - store a new notification record
func (s *NotificationService) Store(ctx context.Context, n *domain.Notification) error {
	return s.repo.Store(ctx, n)
}
