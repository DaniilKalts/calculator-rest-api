.PHONY: docker-run docker-clean start

docker-run:
	@docker rm -f postgres-container >/dev/null 2>&1 || true
	@docker run --name postgres-container \
		-e POSTGRES_PASSWORD=yourpassword \
		-d -p 5433:5432 \
		postgres

docker-clean:
	@docker rm -f postgres-container >/dev/null 2>&1 || true

start: docker-run
	@sleep 1
	@go run ./cmd/server/main.go