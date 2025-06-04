.PHONY: clean docker-up docker-down swagger run-local

clean:
	rm -rf ./sql/data
	rm -rf ./log/*

run-local: swagger
	go run ./cmd/simpleaccount

docker-up:
	docker compose up -d --remove-orphans

docker-down:
	docker compose down

swagger:
	swag init -g ./cmd/simpleaccount/main.go --markdownFiles swagger-markdown --parseDependency true
