module github.com/danilogalisteu/bd-07-gp-chirpy

go 1.22.1

require internal/api v1.0.0

require internal/database v1.0.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.22.0 // indirect
)

replace internal/api => ./internal/api

replace internal/database => ./internal/database
