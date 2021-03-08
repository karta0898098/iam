.PHONY: iam.local
iam.local:
	go run ./cmd/app

.PHONY: docker.build
docker.build:
	docker build -f ./deployments/build/Dockerfile -t iam .

.PHONY: docker.run
CONFIGPATH=
docker.run:
	docker run -v $(CONFIGPATH):/app/deployments/config iam:latest

.PHONY: iam.dev.env
iam.dev.env:
	docker-compose -p iam -f ./deployments/environment/docker-compose.dev.yml up -d

.PHONY: iam.dev.env.down
iam.dev.env.down:
	docker-compose -p iam -f ./deployments/environment/docker-compose.dev.yml down

.PHONY: migrate
DBUSER=iam_program
DBPASSWORD=ETqdG59zTQ4zTrCV
DBNAME=iam
DBHOST=127.0.0.1:8081
migrate:
	goose -dir ./deployments/migrate/ -v mysql "$(DBUSER):$(DBPASSWORD)@tcp($(DBHOST))/$(DBNAME)?charset=utf8&parseTime=True" up