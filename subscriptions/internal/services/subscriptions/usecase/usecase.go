package usecase

import (
	"context"
	"errors"
	"github.com/ekkserapopova/subscriptions/internal/models"
	"github.com/ekkserapopova/subscriptions/internal/services/subscriptions/repo"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log/slog"
	"time"
)

type Params struct {
	fx.In

	Logger *slog.Logger
	Repo   *repo.Repository
}

type UseCase struct {
	log  *slog.Logger
	repo *repo.Repository
}

func NewUseCase(params Params) *UseCase {
	return &UseCase{
		log:  params.Logger,
		repo: params.Repo,
	}
}

func (u *UseCase) CreateSubscription(ctx context.Context, subscriptionData *models.Subscription) (*models.Subscription, error) {
	if subscriptionData.ID == uuid.Nil {
		u.log.Warn("create subscription: id is nil")
		subscriptionData.ID = uuid.New()
	}

	nilTime := time.Time{}

	if subscriptionData.StartDate == models.MonthYear(nilTime) {
		u.log.Warn("create subscription: start date is nil")
		return nil, errors.New("start date is nil")
	}

	createdSubscription, err := u.repo.CreateSubscription(ctx, subscriptionData)
	if err != nil {
		return nil, err
	}

	return createdSubscription, nil
}

func (u *UseCase) UpdateSubscription(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Subscription, error) {
	if id == uuid.Nil {
		u.log.Warn("update subscription: id is nil")
		return nil, errors.New("id is required")
	}

	if len(updates) == 0 {
		u.log.Warn("update subscription: no fields to update")
		return nil, errors.New("no fields to update")
	}

	if startDate, ok := updates["start_date"]; ok {
		if t, ok := startDate.(time.Time); ok && t.IsZero() {
			u.log.Warn("update subscription: start date is nil")
			return nil, errors.New("start date is nil")
		}
	}

	updatedSubscription, err := u.repo.UpdateSubscription(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	return updatedSubscription, nil
}

func (u *UseCase) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}
	return u.repo.GetSubscriptionByID(ctx, id)
}

func (u *UseCase) GetAllSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	return u.repo.GetAllSubscriptions(ctx)
}

func (u *UseCase) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}
	return u.repo.DeleteSubscription(ctx, id)
}

func (u *UseCase) GetSumSubscriptions(ctx context.Context, startDate, endDate, name, usersIds string) (int, error) {
	return u.repo.GetSumSubscriptions(ctx, startDate, endDate, name, usersIds)
}
