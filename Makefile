docker-integration-test:
	docker run --name postgres_test -p 127.0.0.1:5432:5432/tcp -e POSTGRES_PASSWORD=pass -d postgres
compose-up-integration-test:
	docker-compose up --build --abort-on-container-exit --exit-code-from integration-test
mockgen:
	mockgen -source=./internal/handler/handler.go -destination=./internal/handler/mocks_test.go