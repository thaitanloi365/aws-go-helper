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
	@rm -rf /var/folders/mc/fbckgsmd65qb98njnml061nc0000gn/T//output.zip

update_func: remove_junk build_linux_lambda zip
	@read -p "AWS_ACCESS_KEY_ID = " access_key; \
	read -p "AWS_SECRET_ACCESS_KEY = " secret_key; \
	read -p "AWS_REGION = " region; \
	AWS_ACCESS_KEY_ID=$$access_key \
	AWS_SECRET_ACCESS_KEY=$$secret_key \
	AWS_REGION=$$region \
	drone-lambda --region $$region \
	--function-name $(DEPLOY_IMAGE) \
	--source release/linux/lambda/$(DEPLOY_IMAGE)

update_func_without_rebuid: remove_junk
	@read -p "AWS_ACCESS_KEY_ID = " access_key; \
	read -p "AWS_SECRET_ACCESS_KEY = " secret_key; \
	read -p "AWS_REGION = " region; \
	AWS_ACCESS_KEY_ID=$$access_key \
	AWS_SECRET_ACCESS_KEY=$$secret_key \
	AWS_REGION=$$region \
	drone-lambda --region $$region \
	--function-name $(DEPLOY_IMAGE) \
	--source release/linux/lambda/$(DEPLOY_IMAGE)

