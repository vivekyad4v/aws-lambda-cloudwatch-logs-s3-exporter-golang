FILE_TEMPLATE = template.yml
FILE_PACKAGE = ./dist/package.yml

clean:
	@ rm -rf ./dist ./src/*/bin/* && mkdir -p ./dist

configure:
	@ aws s3api head-bucket --bucket $(AWS_BUCKET_NAME) && echo "Artifacts S3 bucket already exists - $(AWS_BUCKET_NAME)" ||\
	     aws s3api create-bucket \
			--bucket $(AWS_BUCKET_NAME) \
			--region $(AWS_REGION) \
			--create-bucket-configuration LocationConstraint=$(AWS_REGION)
	@ aws s3api head-bucket --bucket $(LOG_BUCKET_ID) && echo "Artifacts S3 bucket already exists - $(LOG_BUCKET_ID)" ||\
	     aws s3api create-bucket \
			--bucket $(LOG_BUCKET_ID) \
			--region $(AWS_REGION) \
			--create-bucket-configuration LocationConstraint=$(AWS_REGION)

invoke-local:
	@ sam local invoke \
		--template $(FILE_TEMPLATE) \
		--parameter-overrides \
			"ParameterKey=ParamProjectID,ParameterValue=$(PROJECT_ID) \
			 ParameterKey=ParamProjectEnviron,ParameterValue=$(ENVIRON) \
			 ParameterKey=ParamProjectOrgID,ParameterValue=$(ORG_ID) \
			 ParameterKey=ParamProjectName,ParameterValue=$(PROJECT_NAME) \
			 ParameterKey=ParamLogBucketId,ParameterValue=$(LOG_BUCKET_ID)"

run-local-api:
	@ sam local start-api \
		--template $(FILE_TEMPLATE) \
		--parameter-overrides \
			"ParameterKey=ParamProjectID,ParameterValue=$(PROJECT_ID) \
			 ParameterKey=ParamProjectEnviron,ParameterValue=$(ENVIRON) \
			 ParameterKey=ParamProjectOrgID,ParameterValue=$(ORG_ID) \
			 ParameterKey=ParamProjectName,ParameterValue=$(PROJECT_NAME) \
			 ParameterKey=ParamLogBucketId,ParameterValue=$(LOG_BUCKET_ID)"

package:
	@ sam package \
		--template-file $(FILE_TEMPLATE) \
		--s3-bucket $(AWS_BUCKET_NAME) \
		--region $(AWS_REGION) \
		--output-template-file $(FILE_PACKAGE) 1>/dev/null

validate:
	@ sam validate \
		--template-file $(FILE_TEMPLATE) \
		--region $(AWS_REGION)

deploy:
	@ sam deploy \
		--template-file $(FILE_PACKAGE) \
		--region $(AWS_REGION) \
		--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
		--stack-name $(PROJECT_ID) \
		--parameter-overrides \
			ParamProjectID=$(PROJECT_ID) \
			ParamProjectEnviron=$(ENVIRON) \
			ParamProjectOrgID=$(ORG_ID) \
			ParamProjectName=$(PROJECT_NAME) \
			ParamLogBucketId=$(LOG_BUCKET_ID)

describe:
	@ aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(PROJECT_ID)

outputs:
	@ make describe \
		| jq -r '.Stacks[0].Outputs'

destroy:
	@ aws cloudformation delete-stack \
		--stack-name $(PROJECT_ID)

.PHONY: clean configure build invoke-local run-local-api package deploy describe outputs destroy