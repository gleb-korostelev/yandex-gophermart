package router

import (
	"github.com/gleb-korostelev/gophermart.git/internal/middleware"
	"github.com/gleb-korostelev/gophermart.git/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RouterInit(svc service.APIServiceI, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.GzipCompressMiddleware)
	router.Use(middleware.GzipDecompressMiddleware)
	router.Use(middleware.LoggingMiddleware(logger))
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", svc.Register)
		r.Post("/login", svc.Login)

		r.Route("/", func(r chi.Router) {
			r.Use(middleware.EnsureUserCookie)
			r.Post("/orders", svc.Orders)
			r.Post("/balance/withdraw", svc.Withdraw)
			r.Get("/orders", svc.GetOrders)
			r.Get("/balance", svc.GetBalance)
			r.Get("/withdrawals", svc.GetWithdrawals)
		})
	})
	return router
}
