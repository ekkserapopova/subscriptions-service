package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/ekkserapopova/subscriptions/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"log/slog"
	"strings"
	"time"
)

type Params struct {
	fx.In

	Logger  *slog.Logger
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

type Repository struct {
	pool    *pgxpool.Pool
	log     *slog.Logger
	builder squirrel.StatementBuilderType
}

func NewRepository(params Params) *Repository {
	return &Repository{
		pool:    params.Pool,
		log:     params.Logger,
		builder: params.Builder,
	}
}

func (repo *Repository) CreateSubscription(ctx context.Context, subscriptionData *models.Subscription) (*models.Subscription, error) {
	var endDate *time.Time
	if subscriptionData.EndDate != nil {
		endDate = subscriptionData.EndDate.PtrTime()
	}

	query, args, err := repo.builder.
		Insert("subscriptions").
		Columns("id", "service_name", "price", "user_id", "start_date", "end_date").
		Values(
			subscriptionData.ID,
			subscriptionData.ServiceName,
			subscriptionData.Price,
			subscriptionData.UserID,
			subscriptionData.StartDate.Time(),
			endDate,
		).
		Suffix("RETURNING id, service_name, price, user_id, start_date, end_date").
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	createdSubscription := &models.Subscription{}
	if err := repo.pool.QueryRow(ctx, query, args...).Scan(
		&createdSubscription.ID,
		&createdSubscription.ServiceName,
		&createdSubscription.Price,
		&createdSubscription.UserID,
		&createdSubscription.StartDate,
		&createdSubscription.EndDate,
	); err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			repo.log.Warn("subscription with this id already exists")
			return nil, errors.New("subscription with this id already exists")
		}
		repo.log.Error("failed to create subscription: " + err.Error())
		return nil, err
	}

	return createdSubscription, nil
}

func (repo *Repository) UpdateSubscription(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Subscription, error) {
	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	builder := repo.builder.Update("subscriptions")
	for field, value := range updates {
		builder = builder.Set(field, value)
	}

	query, args, err := builder.
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, service_name, price, user_id, start_date, end_date").
		ToSql()

	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	updatedSubscription := &models.Subscription{}

	if err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&updatedSubscription.ID,
		&updatedSubscription.ServiceName,
		&updatedSubscription.Price,
		&updatedSubscription.UserID,
		&updatedSubscription.StartDate,
		&updatedSubscription.EndDate,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			repo.log.Warn("subscription not found for update")
			return nil, errors.New("subscription not found")
		}
		repo.log.Error("failed to update subscription: " + err.Error())
		return nil, err
	}

	return updatedSubscription, nil
}

func (repo *Repository) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query, args, err := repo.builder.
		Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From("subscriptions").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	sub := &models.Subscription{}
	if err := repo.pool.QueryRow(ctx, query, args...).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("subscription not found")
		}
		return nil, err
	}

	return sub, nil
}

func (repo *Repository) GetAllSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	query, args, err := repo.builder.
		Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From("subscriptions").
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (repo *Repository) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	query, args, err := repo.builder.
		Delete("subscriptions").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}

	cmdTag, err := repo.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

func (repo *Repository) GetSumSubscriptions(ctx context.Context, startDate, endDate, name, usersIds string) (int, error) {
	builder := repo.builder.Select("SUM(price)").From("subscriptions")

	if startDate != "" {
		t, err := time.Parse("01-2006", startDate)
		if err != nil {
			repo.log.Warn("invalid start_date format, expected MM-YYYY: " + err.Error())
		} else {
			startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
			builder = builder.Where(squirrel.GtOrEq{"start_date": startTime})
		}
	}

	if endDate != "" {
		t, err := time.Parse("01-2006", endDate)
		if err != nil {
			repo.log.Warn("invalid end_date format, expected MM-YYYY: " + err.Error())
		} else {
			endTime := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
			builder = builder.Where(squirrel.LtOrEq{"end_date": endTime})
		}
	}

	if name != "" {
		builder = builder.Where(squirrel.Eq{"service_name": name})
	}

	if usersIds != "" {
		var validIDs []uuid.UUID
		for _, idStr := range strings.Split(usersIds, ",") {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := uuid.Parse(idStr)
			if err != nil {
				continue
			}
			validIDs = append(validIDs, id)
		}

		if len(validIDs) > 0 {
			builder = builder.Where(squirrel.Eq{"user_id": validIDs})
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return 0, err
	}

	var sum sql.NullInt64

	if err := repo.pool.QueryRow(ctx, query, args...).Scan(&sum); err != nil {
		repo.log.Error("failed to fetch sum subscription: " + err.Error())
		return 0, err
	}

	if !sum.Valid {
		return 0, nil
	}

	return int(sum.Int64), nil
}
