package app

import (
	"github.com/go-chi/chi"
	"github.com/zlobste/fake-wallet/internal/app/handlers"
	"github.com/zlobste/fake-wallet/internal/app/middlewares"
	"github.com/zlobste/fake-wallet/internal/app/utils"
	"github.com/zlobste/fake-wallet/internal/data/postgres"
)

func (a *app) router() chi.Router {
	router := chi.NewRouter()

	router.Use(
		middlewares.CorsMiddleware(),
		middlewares.LoggingMiddleware(a.log),
		middlewares.CtxMiddleware(
			utils.CtxLog(a.log),
			utils.CtxJWT(a.auth),
			utils.CtxUsers(postgres.NewUsersStorage(a.db)),
			utils.CtxWallets(postgres.NewWalletsStorage(a.db)),
			utils.CtxTransactions(postgres.NewTransactionsStorage(a.db)),
			utils.CtxAssets(postgres.NewAssetsStorage(a.db)),
		),
	)

	router.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlers.Register)
			r.Post("/login", handlers.Login)
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Route("/wallets", func(r chi.Router) {
			r.Get("/", handlers.GetUserWallets)
			r.Post("/transfer", handlers.TransferFunds)
		})

		r.Get("/transactions", handlers.GetUserTxs)
	})

	return router
}
