.PHONY: sqlctest
test:
	go test ./...

.PHONY: db
db:
	sqlc generate -f model/sqlc.yaml