VERSION := $(shell git describe --tags --dirty)

RELEASE_DIR := releases/$(VERSION)
RELEASE_BIN := releases/$(VERSION)/bin

build:
	go fmt alkira/*
	go build -o bin/terraform-provider-alkira

test:
	go fmt alkira/*.go
	go test ./alkira/... $(VERBOSE)

release:
	GOOS=windows GOARCH=amd64 go build -o $(RELEASE_BIN)/windows-amd64/terraform-provider-alkira_$(VERSION)
	GOOS=linux GOARCH=amd64 go build -o $(RELEASE_BIN)/linux-amd64/terraform-provider-alkira_$(VERSION)
	GOOS=linux GOARCH=arm64 go build -o $(RELEASE_BIN)/linux-arm64/terraform-provider-alkira_$(VERSION)
	GOOS=darwin GOARCH=amd64 go build -o $(RELEASE_BIN)/darwin-amd64/terraform-provider-alkira_$(VERSION)
	GOOS=darwin GOARCH=arm64 go build -o $(RELEASE_BIN)/darwin-arm64/terraform-provider-alkira_$(VERSION)
	tar czf $(RELEASE_DIR)/terraform-provider-alkira-$(VERSION)-linux-amd64.tar.gz -C $(RELEASE_BIN)/linux-amd64 terraform-provider-alkira_$(VERSION)
	tar czf $(RELEASE_DIR)/terraform-provider-alkira-$(VERSION)-linux-arm64.tar.gz -C $(RELEASE_BIN)/linux-arm64 terraform-provider-alkira_$(VERSION)
	tar czf $(RELEASE_DIR)/terraform-provider-alkira-$(VERSION)-darwin-amd64.tar.gz -C $(RELEASE_BIN)/darwin-amd64 terraform-provider-alkira_$(VERSION)
	tar czf $(RELEASE_DIR)/terraform-provider-alkira-$(VERSION)-darwin-arm64.tar.gz -C $(RELEASE_BIN)/darwin-arm64 terraform-provider-alkira_$(VERSION)
	zip $(RELEASE_DIR)/terraform-provider-alkira-$(VERSION)-windows-amd64.zip -j $(RELEASE_BIN)/windows-amd64/terraform-provider-alkira_$(VERSION)
	rm -rf $(RELEASE_BIN)

fmt:
	go fmt alkira/*

doc:
	tfplugindocs generate

vendor: GOPRIVATE=github.com/alkiranet
vendor:
	go mod tidy
	go mod vendor

superclean:
	git clean -x -d -f
