REGISTRY := "docker.io/despitehowever"
CONTAINER := argo-workflows-webhook
CONTAINER_TAG := latest

.PHONY: logs
logs:
	kubectl logs -f $(shell kubectl get pods -l 'app==argo-webhook' -o yaml | yq ".items[0].metadata.name")

.PHONY: deploy
deploy:
	kustomize build deploy/ | kubectl apply -f -

.PHONY: redeploy
redeploy:
	kubectl delete pod $(shell kubectl get pods -l 'app==argo-webhook' -o yaml | yq ".items[0].metadata.name")

.PHONY: build
build:
	docker buildx build --platform linux/amd64 -t $(REGISTRY)/$(CONTAINER):$(CONTAINER_TAG) --push .

.PHONY: html
html: test
	go tool cover -html=coverage.out -o coverage.html && open coverage.html

.PHONY: test
test:
	#go test -shuffle=on -race -coverprofile=coverage.txt -covermode=atomic ./...
	go test -shuffle=on -race -coverprofile=coverage.txt -covermode=atomic $$(go list ./... | grep -v /cmd/)
