# aws-lambda-cloudwatch-logs-s3-exporter-golang

#### 1. Export necessary variables
``` 
    export ORG_ID=foo
    export ENVIRON=uat
    export PROJECT_NAME=play-with-stores    
```

#### 2. Deploy locally

```
    make clean build configure run-local
```

#### 3. Deploy on your AWS account

```
    make clean build configure package validate deploy describe outputs
```
