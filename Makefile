compose: 
	docker compose up -d
run:
	go run main.go
.PHONY: 
	compose run