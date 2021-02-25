ORG_ID ?= hpy
ENVIRON ?= uat
PROJECT_NAME ?= $(notdir $(CURDIR))
PROJECT_ID ?= $(ORG_ID)-$(ENVIRON)-$(PROJECT_NAME)
AWS_BUCKET_NAME ?= $(ORG_ID)-$(ENVIRON)-sls-artifacts-$(AWS_REGION)
AWS_REGION ?= $(AWS_DEFAULT_REGION)
LOG_BUCKET_ID ?= $(ORG_ID)-$(ENVIRON)-exported-cloudwatch-logs-$(AWS_REGION)


export SAM_CLI_TELEMETRY=0
export GO111MODULE=on

### DEFINE MODULE NAMES ###
export MODULE1=cloudwatch-to-s3-logs

###### Sample for environment segregation for local development
ifeq ($(ENVIRON), uat)
	BUILD=umake
endif
ifeq ($(ENVIRON), prd)
	BUILD=pmake
endif

include .make

build:
	@ cd ./src/$(MODULE1) ; go mod init $(MODULE1) ; go mod download ; GOOS=linux go build -o ./bin/$(MODULE1)
