VERSION := $(shell git describe --tags --dirty)

RELEASE_DIR := releases/$(VERSION)
RELEASE_BIN := releases/$(VERSION)/bin

build:
	go build -o bin/terraform-provider-alkira

release:
	GOOS=windows GOARCH=amd64 go build -o $(RELEASE_BIN)/windows-amd64/terraform-provider-alkira_$(VERSION)
	GOOS=linux GOARCH=amd64 go build -o $(RELEASE_BIN)/linux-amd64/terraform-provider-alkira_$(VERSION)
	GOOS=darwin GOARCH=amd64 go build -o $(RELEASE_BIN)/darwin-amd64/terraform-provider-alkira_$(VERSION)
	tar cvzf $(RELEASE_DIR)/terraform-provider-alkira-$(VERSION)-linux-amd64.tar.gz -C $(RELEASE_BIN)/linux-amd64 terraform-provider-alkira_$(VERSION)
	tar cvzf $(RELEASE_DIR)/terraform-provider-alkira-$(VERSION)-darwin-amd64.tar.gz -C $(RELEASE_BIN)/darwin-amd64 terraform-provider-alkira_$(VERSION)
	zip $(RELEASE_DIR)/terraform-provider-alkira-$(VERSION)-windows-amd64.zip -j $(RELEASE_BIN)/windows-amd64/terraform-provider-alkira_$(VERSION)

superclean:
	git clean -x -d -f
