iam-local:
	go run ./cmd/identity

docker-build:
	docker build -f ./deployments/build/Dockerfile -t iam .

docker-deploy:
	docker build -f ./deployments/build/Dockerfile -t iam . && \
	docker-compose -p iam -f ./deployments/environment/docker-compose.deploy.yml up -d

docker-deploy-down:
	docker-compose -p iam -f ./deployments/environment/docker-compose.deploy.yml down

iam-env:
	docker-compose -p iam -f ./deployments/environment/docker-compose.dev.yml up -d

iam-env-down:
	docker-compose -p iam -f ./deployments/environment/docker-compose.dev.yml down


.PHONY: migrate
DBUSER=iam_program
DBPASSWORD=ETqdG59zTQ4zTrCV
DBNAME=iam
DBHOST=127.0.0.1
DBPORT=8081
migrate:
	goose -dir ./deployments/migrate/ -v postgres "user=$(DBUSER) password=$(DBPASSWORD) dbname=$(DBNAME) host=$(DBHOST) port=$(DBPORT) sslmode=disable" up
mocks:
	mockery --all --with-expecter --dir ./pkg/app/identity --output ./pkg/app/identity/mocks

proto:
	$(foreach dir, protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=require_unimplemented_servers=false:. \
	--go-grpc_opt=paths=source_relative $(shell find . -name *.proto),$(dir))
