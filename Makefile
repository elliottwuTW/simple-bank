redis:
	docker run --name redis -p 6379:6379 -d redis:8-alpine
mock:
	mockgen -package mockdb -destination ./database/mock/database.go github.com/simple_bank/database Database
test:
	go test -v -cover ./...

.PHONY: test mock redis