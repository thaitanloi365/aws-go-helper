LDFLAGS ?= -X 'main.Version=$(VERSION)'
DEPLOY_IMAGE := aws-go-helper
ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

build_linux_lambda:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags 'lambda' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/lambda/$(DEPLOY_IMAGE)

zip:
	zip -j deployment.zip release/linux/lambda/$(DEPLOY_IMAGE);\

remove_junk:
	@rm -rf /var/folders/hv/4hfjh9kx77gd2ympb1_z09wr0000gn/T//output.zip

update_func: remove_junk build_linux_lambda zip
	AWS_ACCESS_KEY_ID=AKIASTMGVDQSSJ7GUIU6 \
	AWS_SECRET_ACCESS_KEY=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 \
	drone-lambda --region us-east-1 \
	--function-name $(DEPLOY_IMAGE) \
	--source release/linux/lambda/$(DEPLOY_IMAGE)

