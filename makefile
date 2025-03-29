.PHONY: run migrate

run:
	go run cmd/main.go

migrate:
	go run cmd/main.go --migrate

seed:
	go run cmd/main.go --seed

start:
	go run cmd/main.go --migrate --seed
	
