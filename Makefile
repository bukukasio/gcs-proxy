IMG := gcr.io/beecash-prod/infra/gcs-proxy
SHA := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
TAG := ${BRANCH}-${SHA}

docker-build:
	docker build -t ${IMG}:${TAG} .

docker-push: docker-build
	docker push ${IMG}:${TAG}
