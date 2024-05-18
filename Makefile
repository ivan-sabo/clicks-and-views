ROOT_DIR:=$(realpath $(shell dirname $(firstword $(MAKEFILE_LIST))))

swagger:
	docker run -p 80:8080 -e SWAGGER_JSON=/docs/openapi.yaml -v $(ROOT_DIR)/api:/docs swaggerapi/swagger-ui

run:
	go run cmd/main.go