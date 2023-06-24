package review

import (
	"github.com/go-chi/chi"
)

var ExchangeInfoServiceSubRoute = chi.NewRouter()

func init() {
	ExchangeInfoServiceSubRoute.Group(func(r chi.Router) {

	})
}
