make:
	@templ generate 
	@docker compose up -d
	@go build -o bin/main cmd/main/main.go
	@./bin/main
