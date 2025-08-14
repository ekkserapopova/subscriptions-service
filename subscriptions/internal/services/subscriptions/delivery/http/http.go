package http

import (
	"github.com/ekkserapopova/subscriptions/internal/models"
	"github.com/ekkserapopova/subscriptions/internal/services/subscriptions/usecase"
	"github.com/ekkserapopova/subscriptions/pkg/reader"
	"github.com/ekkserapopova/subscriptions/pkg/responser"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"time"
)

type Params struct {
	fx.In

	Logger  *slog.Logger
	UseCase *usecase.UseCase
}

type Handler struct {
	logger   *slog.Logger
	useacase *usecase.UseCase
}

func NewHandler(params Params) *Handler {
	return &Handler{
		logger:   params.Logger,
		useacase: params.UseCase,
	}
}

// CreateSubscription godoc
// @Summary Создать подписку
// @Description Создает новую подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Подписка"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionData := &models.Subscription{}
	subscriptionData.ID = uuid.New()
	if err := reader.ReadResponseData(r, subscriptionData); err != nil {
		h.logger.Error("create subscription request err: " + err.Error())
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if subscriptionData.ServiceName == "" {
		h.logger.Error("create subscription request err: service name is nil")
		responser.SendErr(w, http.StatusBadRequest, "create subscription request err: service name is nil")
		return
	}

	nilTime := time.Time{}

	if subscriptionData.StartDate == models.MonthYear(nilTime) {
		h.logger.Error("create subscription request err: start date is nil")
		responser.SendErr(w, http.StatusBadRequest, "create subscription request err: start date is nil")
		return
	}

	if subscriptionData.Price == nil {
		h.logger.Error("create subscription request err: price is nil")
		responser.SendErr(w, http.StatusBadRequest, "create subscription request err: price is nil")
		return
	}

	if subscriptionData.UserID == uuid.Nil {
		h.logger.Error("create subscription request err: user_id is nil")
		responser.SendErr(w, http.StatusBadRequest, "create subscription request err: user_id is nil")
		return
	}

	createdSubscription, err := h.useacase.CreateSubscription(r.Context(), subscriptionData)
	if err != nil {
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOK(w, http.StatusCreated, createdSubscription)
}

// @Summary Изменить подписку
// @Description Изменяет существующую подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param subscription body models.Subscription true "Подписка"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		h.logger.Error("update subscription request err: id is nil")
		responser.SendErr(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("update subscription request err: invalid id format")
		responser.SendErr(w, http.StatusBadRequest, "invalid id format")
		return
	}

	updates := make(map[string]interface{})
	if err := reader.ReadResponseData(r, &updates); err != nil {
		h.logger.Error("update subscription request err: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(updates) == 0 {
		h.logger.Error("update subscription request err: no fields to update")
		responser.SendErr(w, http.StatusBadRequest, "no fields to update")
		return
	}

	if startDate, ok := updates["start_date"]; ok {
		if t, ok := startDate.(string); ok {
			parsedTime, parseErr := time.Parse("2006-01", t)
			if parseErr != nil || parsedTime.IsZero() {
				h.logger.Error("update subscription request err: invalid start_date")
				responser.SendErr(w, http.StatusBadRequest, "invalid start_date")
				return
			}
			updates["start_date"] = parsedTime
		}
	}

	updatedSubscription, err := h.useacase.UpdateSubscription(r.Context(), id, updates)
	if err != nil {
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOK(w, http.StatusOK, updatedSubscription)
}

// @Summary Получить запись об одной подписке
// @Description Получить запись об одной подписке по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		responser.SendErr(w, http.StatusBadRequest, "invalid id format")
		return
	}

	sub, err := h.useacase.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		responser.SendErr(w, http.StatusNotFound, err.Error())
		return
	}

	responser.SendOK(w, http.StatusOK, sub)
}

// @Summary Получить все подписки
// @Description Получить записи всех подписок
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *Handler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.useacase.GetAllSubscriptions(r.Context())
	if err != nil {
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOK(w, http.StatusOK, subs)
}

// @Summary Удалить подписку
// @Description Удалить подписку по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		responser.SendErr(w, http.StatusBadRequest, "invalid id format")
		return
	}

	if err := h.useacase.DeleteSubscription(r.Context(), id); err != nil {
		if err.Error() == "subscription not found" {
			responser.SendErr(w, http.StatusNotFound, err.Error())
		} else {
			responser.SendErr(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	responser.SendOK(w, http.StatusNoContent, map[string]string{"msg": "subscription deleted"})
}

// @Summary Получить суммарную стоимость подписок
// @Description Получить суммарную стоимость подписок с фильтрацией по дате, названию сервиса и по пользователям
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param start_date query string false "Дата начала фильтрации в формате YYYY-MM"
// @Param end_date query string false "Дата окончания фильтрации в формате YYYY-MM"
// @Param name query string false "Название сервиса"
// @Param users_ids query string false "Список ID пользователей через запятую"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /subscriptions/sum [get]
func (h *Handler) GetSumSubscriptions(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	nameService := r.URL.Query().Get("name")
	usersIds := r.URL.Query().Get("users_ids")

	sum, err := h.useacase.GetSumSubscriptions(r.Context(), startDate, endDate, nameService, usersIds)
	if err != nil {
		h.logger.Error("get sum subscriptions err: " + err.Error())
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOK(w, http.StatusOK, map[string]interface{}{
		"sum": sum,
	})
}
