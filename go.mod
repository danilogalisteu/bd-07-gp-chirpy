module github.com/danilogalisteu/bd-07-gp-chirpy

go 1.22.1

require internal/database v1.0.0

require golang.org/x/crypto v0.22.0 // indirect

replace internal/database => ./internal/database
