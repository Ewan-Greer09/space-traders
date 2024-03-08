make:
	@templ generate 
	@docker compose up -d
	@go run cmd/main/main.go
