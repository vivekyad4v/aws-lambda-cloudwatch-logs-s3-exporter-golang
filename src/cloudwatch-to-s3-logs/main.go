package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Happay-DevSecOps/go-utils/cdate"
	"github.com/Happay-DevSecOps/go-utils/cerrors"
	"github.com/Happay-DevSecOps/go-utils/couts"
	"github.com/Happay-DevSecOps/go-utils/logger"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	limitList            int64    = 1
	logGroupNamePrefixes []string = []string{"/aws/transfer", "/aws/vendedlogs"}
	// logGroupNamePrefixes       []string = []string{"/s/transfer", "/alogs"}
	logS3BucketID              string = os.Getenv("LOG_BUCKET_ID")
	dateDetail                 cdate.GetDate
	logInfo, logWarn, logError *log.Logger
	cerrorgeneral              *cerrors.GeneralError
	coutgeneral                *couts.GeneralOutput
)

func init() {
	logInfo, logWarn, logError = logger.InitLogger()
}

func mainHandler(ctx context.Context) (new interface{}, err error) {

	fmt.Println("")
	dateDetail = dateDetail.Ymdyesterday()

	logGroupsToExport := fetchCloudwatchLogGroups(logGroupNamePrefixes)
	if len(logGroupsToExport) < 1 {
		cerrorg := cerrors.GeneralError{Code: "501", Message: "No log groups matched for export", Err: errors.New("NoMatch")}
		logError.Println(&cerrorg)
		return cerrorgeneral, errors.New("NoMatch")
	}

	logInfo.Println("Log groups to export => ", logGroupsToExport)
	fmt.Println("")

	coutgeneral = s3ExportTaskStart(logGroupsToExport, dateDetail)
	return coutgeneral, nil
}

func main() {
	lambda.Start(mainHandler)
}
