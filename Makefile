.PHONY: html
html: test
	go tool cover -html=coverage.out -o coverage.html && open coverage.html

.PHONY: test
test:
	#go test -shuffle=on -race -coverprofile=coverage.txt -covermode=atomic $$(go list ./... | grep -v /cmd/)
	go test -shuffle=on -race -coverprofile=coverage.txt -covermode=atomic ./...
