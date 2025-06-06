.PHONY: run-service test run-local

run-service:
	docker-compose -f deployment/docker-compose.yml --project-directory . up

test:
	./coverage.sh

run-local:
	go run main.go run-http --config_file=config-local.json

run-db:
	docker-compose -f deployment/docker-compose-db.yml up

run-migration:
	go run main.go run-migration --config_file=config-local.json

run-test-api:
	k6 run example-test-script.js