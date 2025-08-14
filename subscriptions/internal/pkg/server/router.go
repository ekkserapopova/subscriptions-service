package server

import (
	_ "github.com/ekkserapopova/subscriptions/docs"
	subscriptionHandler "github.com/ekkserapopova/subscriptions/internal/services/subscriptions/delivery/http"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
)

type RouterParams struct {
	fx.In

	Logger              *slog.Logger
	SubscriptionHandler *subscriptionHandler.Handler
}

type Router struct {
	handler *mux.Router
}

func NewRouter(p RouterParams) *Router {
	api := mux.NewRouter().PathPrefix("/api").Subrouter()
	v1 := api.PathPrefix("/v1").Subrouter()

	v1.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	v1.HandleFunc("/subscriptions", p.SubscriptionHandler.CreateSubscription).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/subscriptions/sum", p.SubscriptionHandler.GetSumSubscriptions).Methods(http.MethodGet)
	v1.HandleFunc("/subscriptions/{id}", p.SubscriptionHandler.UpdateSubscription).Methods(http.MethodPut, http.MethodOptions)
	v1.HandleFunc("/subscriptions", p.SubscriptionHandler.GetAllSubscriptions).Methods(http.MethodGet)
	v1.HandleFunc("/subscriptions/{id}", p.SubscriptionHandler.GetSubscriptionByID).Methods(http.MethodGet)
	v1.HandleFunc("/subscriptions/{id}", p.SubscriptionHandler.DeleteSubscription).Methods(http.MethodDelete)

	router := &Router{
		handler: api,
	}

	p.Logger.Info("registered router")

	return router
}
