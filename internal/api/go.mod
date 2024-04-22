module api

go 1.22.1

require internal/database v1.0.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
)

replace internal/database => ../database
