.PHONY: sqlctest
test:
	go test ./...

.PHONY: db
db:
	sqlc generate -f model/sqlc.yaml

.PHONY: api
api:
	oapi-codegen -generate types -o "./controller/openapi_types.gen.go" -package "controller" ./controller/openapi.yaml
	oapi-codegen -generate std-http -o "./controller/openapi_api.gen.go" -package "controller" ./controller/openapi.yaml