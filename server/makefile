all:
	go run main.go

.PHONY: migrate-up
migrate-up:
	migrate -verbose -source file://./database/migration -database mysql://root:uwuowo@tcp\(db:3306\)/dcard-backend-intern-assignment-2024 up

.PHONY: migrate-version
migrate-version:
	migrate -source file://./database/migration -database mysql://root:uwuowo@tcp\(db:3306\)/dcard-backend-intern-assignment-2024 version

.PHONY: migrate-down
migrate-down:
	migrate -verbose -source file://./database/migration -database mysql://root:uwuowo@tcp\(db:3306\)/dcard-backend-intern-assignment-2024 down

.PHONY: build
build:
	go build -o ./dcard-ad-assignment
	mv dcard-ad-assignment /usr/local/bin

.PHONY: test
test:
	go test ./test/...

.PHONY: swag
swag:
	swag init