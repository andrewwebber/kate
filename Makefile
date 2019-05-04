all: build

build: bin/kate

.PHONY: bin/kate
bin/kate:
	@CGO_ENABLED=0 GOOS=linux go build -v -a -tags netgo -o ./bin/kate github.com/andrewwebber/kate

.PHONY: docker-image
docker-image: build
	@docker build -t andrewwebber/kate .
.PHONY: docker-push
docker-push: docker-image
	@docker push andrewwebber/kate
