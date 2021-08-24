package middlewares

import "net/http"

func CorsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			enableCors(&w)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
