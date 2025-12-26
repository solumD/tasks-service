# запуск локально
run-locally:
	go build -o bin/tasks-service cmd/app/main.go
	bin/tasks-service

# запуск в контейнере docker
run:
	docker build -t tasks-service .
	docker run -p 8080:8080 tasks-service

# тесты транспортного слоя
test-handler:
	go test -v ./internal/handler/v1/tests/

# тесты бизнес логики
test-usecase:
	go test -v ./internal/usecase/tests/

# выполнение всех тестов
test:
	make test-handler
	make test-usecase