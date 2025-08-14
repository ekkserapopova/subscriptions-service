package subscriptions

import (
	"context"
	"github.com/ekkserapopova/subscriptions/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateSubscription(ctx context.Context, subscriptionData *models.Subscription) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetAllSubscriptions(ctx context.Context) ([]*models.Subscription, error)
	GetSumSubscriptions(ctx context.Context, startDate, endDate, name, usersIds string) (int, error)
}

type Repository interface {
	CreateSubscription(ctx context.Context, subscriptionData *models.Subscription) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetAllSubscriptions(ctx context.Context) ([]*models.Subscription, error)
	GetSumSubscriptions(ctx context.Context, startDate, endDate, name, usersIds string) (int, error)
}
