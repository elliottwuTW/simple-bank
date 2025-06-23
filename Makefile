mock:
	mockgen -package mockdb -destination ./database/mock/database.go github.com/simple_bank/database Database
test:
	go test -v -cover ./...

.PHONY: test mock