package middlewares

import "github.com/go-chi/cors"

func GetCors(domains ...string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   domains,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
