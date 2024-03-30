package handlers

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	authService authService,
	orderService orderService,
	balanceService balanceService,
	accrualService accrualService,
) chi.Router {
	r := chi.NewRouter()
	h := NewHandler(authService, orderService, balanceService, accrualService)

	r.Post(`/api/user/register`, h.Register)
	r.Post(`/api/user/login`, h.Login)

	r.Group(func(r chi.Router) {
		r.Use(authService.NewJWTVerifyMiddleware())
		r.Post(`/api/user/orders`, h.UploadOrder)
		r.Get(`/api/user/orders`, h.GetOrders)
		r.Get(`/api/user/balance`, h.GetBalance)
		r.Post(`/api/user/balance/withdraw`, h.Withdraw)
		r.Get(`/api/user/withdrawals`, h.GetWithdrawals)
	})

	return r
}
