# Vars
STG_TAG=stage
PROD_TAG=v0.0.1
IMAGE_NAME=raddare

# Misc
BINARY_NAME=raddare
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	go  build -o ./bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME).go

build-linux:
	CGOENABLED=0 GOOS=linux GOARCH=amd64; go build -o ./bin/$(BINARY_UNIX) ./cmd/$(BINARY_NAME).go

test:
	./scripts/test.sh

clean:
	go clean
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

run:
	./scripts/run.sh

connect-stg:
	gcloud beta container clusters get-credentials ${GC_STG_CLUSTER} --region ${GC_REGION} --project ${GC_STG_PROJECT}

connect-prod:
	gcloud  beta container clusters get-credentials ${GC_PROD_CLUSTER} --region ${GC_REGION} --project ${GC_PROD_PROJECT}

build-stg:
	make build
	docker login
	docker build -t ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(STG_TAG) .
	docker push ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(STG_TAG)

build-prod:
	make build
	docker login
	docker build -t ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(PROD_TAG) .
	docker push ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(PROD_TAG)

template-stg:
	helm template --name $(IMAGE_NAME) -f ./deployments/helm/values-stg.yaml ./deployments/helm

template-prod:
	helm template --name $(IMAGE_NAME) -f ./deployments/helm/values-prod.yaml ./deployments/helm

install-stg:
	make connect-stg
	helm install --name $(IMAGE_NAME) -f ./deployments/helm/values-stg.yaml ./deployments/helm

install-prod:
	make connect-prod
	helm install --name $(IMAGE_NAME) -f ./deployments/helm/values-prod.yaml ./deployments/helm

delete-stg:
	make connect-stg
	helm del --purge $(IMAGE_NAME)

delete-prod:
	make connect-prod
	helm del --purge $(IMAGE_NAME)

deploy-stg:
	make build-stg
	make connect-stg
	make delete-stg
	make install-stg

deploy-prod:
	make build-prod
	make connect-prod
	make delete-prod
	make install-prod

custom-build:
	make mod tidy; go mod vendor; go build ./...

current-conn:
	kubectl config current-context

get-deps:
	go get "gitlab.com/mikrowezel/config"
	go get "gitlab.com/mikrowezel/log"
	go get "gitlab.com/mikrowezel/service"
