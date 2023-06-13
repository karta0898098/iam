.PHONY: iam.local
iam.local:
	go run ./cmd/app

.PHONY: docker.build
docker.build:
	docker build -f ./deployments/build/Dockerfile -t iam .

.PHONY: docker.deploy
docker.deploy:
	docker build -f ./deployments/build/Dockerfile -t iam . && \
	docker-compose -p iam -f ./deployments/environment/docker-compose.deploy.yml up -d

.PHONY: docker.deploy.down
docker.deploy.down:
	docker-compose -p iam -f ./deployments/environment/docker-compose.deploy.yml down

.PHONY: iam.dev.env
iam.dev.env:
	docker-compose -p iam -f ./deployments/environment/docker-compose.dev.yml up -d

.PHONY: iam.dev.env.down
iam.dev.env.down:
	docker-compose -p iam -f ./deployments/environment/docker-compose.dev.yml down

.PHONY: uint.testing
uint.testing:
	go clean -testcache && go test ./...

.PHONY: migrate
DBUSER=iam_program
DBPASSWORD=ETqdG59zTQ4zTrCV
DBNAME=iam
DBHOST=127.0.0.1:8081
migrate:
	goose -dir ./deployments/migrate/ -v mysql "$(DBUSER):$(DBPASSWORD)@tcp($(DBHOST))/$(DBNAME)?charset=utf8&parseTime=True" up

mocks:
	mockery --all --with-expecter --dir ./pkg/app/identity  --output ./pkg/app/identity/mocks

proto:
	$(foreach dir, protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative $(shell find . -name *.proto),
	$(dir))
