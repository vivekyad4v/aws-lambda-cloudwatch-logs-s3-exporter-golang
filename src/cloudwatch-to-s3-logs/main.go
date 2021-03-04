package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var (
	limitList            int64    = 1
	logGroupNamePrefixes []string = []string{"/aws/transfer", "/aws/vendedlogs"}
	logS3BucketID        string   = "hpy-uat-exported-cloudwatch-logs-ap-south-1"
)

func initAwsCloudwatchLogsSession() (svcCloudwatchLogs *cloudwatchlogs.CloudWatchLogs) {
	mySession := session.Must(session.NewSession())
	return cloudwatchlogs.New(mySession)
}

func fetchCloudwatchLogGroups(logGroupNamePrefixes []string) (cloudwatchLogGroupNames []string, err error) {
	svcCloudwatchLogs := initAwsCloudwatchLogsSession()

	for k, v := range logGroupNamePrefixes {
		input := &cloudwatchlogs.DescribeLogGroupsInput{Limit: &limitList, LogGroupNamePrefix: &v}
		describeLogGroups, err := svcCloudwatchLogs.DescribeLogGroups(input)
		for describeLogGroups.NextToken != nil {
			newInput := &cloudwatchlogs.DescribeLogGroupsInput{Limit: &limitList, LogGroupNamePrefix: &logGroupNamePrefixes[k], NextToken: describeLogGroups.NextToken}
			describeLogGroups, err = svcCloudwatchLogs.DescribeLogGroups(newInput)
			cloudwatchLogGroupNames = append(cloudwatchLogGroupNames, *describeLogGroups.LogGroups[0].LogGroupName)

			if err != nil {
				fmt.Println("Error in - ", describeLogGroups.LogGroups, "-", err)
			}
		}
		cloudwatchLogGroupNames = append(cloudwatchLogGroupNames, *describeLogGroups.LogGroups[0].LogGroupName)
	}
	return cloudwatchLogGroupNames, err
}

func mainHandler(ctx context.Context) (output string, err error) {

	var dateDetail getDate
	dateDetail = dateDetail.ymdyesterday()

	logGroupsToExport, err := fetchCloudwatchLogGroups(logGroupNamePrefixes)
	if len(logGroupsToExport) < 1 {
		fmt.Println("No log groups matched for export", logGroupsToExport)
	}
	fmt.Println("Log groups to export => ", logGroupsToExport)

	taskIDList := s3ExportTaskStart(logGroupsToExport, dateDetail)

	fmt.Println(taskIDList)

	return "Sucess", err
}

func main() {
	lambda.Start(mainHandler)
}

// snippets
// fmt.Println("Checking for yesterday's export status.......")
// checkForFailedTasks, err := checkForFailedTasks()
// if checkForFailedTasks == "FAILED" || err != nil {
// 	fmt.Println("Tasks FAILED yesterday, status => ", checkForFailedTasks)
// } else {
// 	fmt.Println("Tasks SUCCEEDED yesterday")
// }

// func checkForFailedTasks() (status string, err error) {
// 	svcCloudwatchLogs := initAwsCloudwatchLogsSession()

// 	exportTaskStatusCodeFailed := cloudwatchlogs.ExportTaskStatusCodeFailed
// 	describeExportTasksInputFailed := &cloudwatchlogs.DescribeExportTasksInput{StatusCode: &exportTaskStatusCodeFailed}

// 	exportTasksFailed, err := svcCloudwatchLogs.DescribeExportTasks(describeExportTasksInputFailed)

// 	if len(exportTasksFailed.ExportTasks) != 0 {
// 		fmt.Println("There are FAILED export tasks.", exportTasksFailed)
// 		return "FAILED", err
// 	}

// 	return "ALL TASKS SUCCEEDED", err
// }
