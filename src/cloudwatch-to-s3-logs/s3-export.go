package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Happay-DevSecOps/go-utils/cdate"
	"github.com/Happay-DevSecOps/go-utils/cerrors"
	"github.com/Happay-DevSecOps/go-utils/couts"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/s3"
)

func s3ExportTaskStart(cloudwatchLogGroupNames []string, getDate cdate.GetDate) (coutgeneral *couts.GeneralOutput) {

	var taskIDList []string
	type exportTaskIDList struct {
		TaskIDs []string
	}

	svcCloudwatchLogs := initAwsCloudwatchLogsSession()

	var destinationPrefix string = getDate.Year + "/" + getDate.Month + "/" + getDate.Day
	var destinationKey string = getDate.Year + "-" + getDate.Month + "-" + getDate.Day + ".json"

	for k := range cloudwatchLogGroupNames {
		createExportTaskInput := &cloudwatchlogs.CreateExportTaskInput{
			Destination:       &logS3BucketID,
			DestinationPrefix: &destinationPrefix,
			From:              &getDate.Yms,
			To:                &getDate.Tms,
			LogGroupName:      &cloudwatchLogGroupNames[k],
		}
		fmt.Println("")
		logInfo.Println("Exporting log group =>", cloudwatchLogGroupNames[k])
		createExportTaskOutput, err := svcCloudwatchLogs.CreateExportTask(createExportTaskInput)
		if err != nil {
			cerrorg := cerrors.GeneralError{Code: "502", Message: "Error in exporting =>", Err: err}
			logError.Println(cerrorg)
			continue
		}
		taskID, Status, err := waitForExportTaskCompletion(*createExportTaskOutput.TaskId)
		logInfo.Println("Executed Task ID => ", taskID, ", Status => ", Status, ", Error => ", err)
		taskIDList = append(taskIDList, *createExportTaskOutput.TaskId)
		logInfo.Println("Task Id list =>", taskIDList)
		fmt.Println("")
	}

	// PUT TASKIDs to S3 BUCKET
	logInfo.Println("Putting taskIDs to S3 bucket....")

	// Converting to JSON
	taskIDToJSON := exportTaskIDList{TaskIDs: taskIDList}
	taskIDListJSON, _ := json.Marshal(taskIDToJSON)
	readerJSON := strings.NewReader(string(taskIDListJSON))

	// Creating S3 session
	sess, err := session.NewSession()
	if err != nil {
		logError.Println("Error - Unable to create s3 AWS session!")
	}
	s3svc := s3.New(sess)

	// Putting object to S3
	s3PutObjectInput := &s3.PutObjectInput{
		Bucket: &logS3BucketID,
		Key:    &destinationKey,
		Body:   readerJSON,
	}

	_, err = s3svc.PutObject(s3PutObjectInput)
	if err != nil {
		logError.Println("Error - Unable to put object to s3!")
	}

	fmt.Println("")
	coutgeneral = &couts.GeneralOutput{Code: "200", Message: "Log groups exported successully!", Out: taskIDList}

	return coutgeneral
}

func waitForExportTaskCompletion(taskID string) (string, string, error) {
	svcCloudwatchLogs := initAwsCloudwatchLogsSession()

	describeExportTasksInput := &cloudwatchlogs.DescribeExportTasksInput{
		TaskId: &taskID,
	}

	time.Sleep(time.Second * 10)
	describeExportTasksOutput, err := svcCloudwatchLogs.DescribeExportTasks(describeExportTasksInput)
	task := describeExportTasksOutput.ExportTasks[0]
	logInfo.Println("Current Task ID => ", *task.TaskId, ", Status => ", *task.Status.Code)

	waitErrorCount := 0
	exportTaskCompleted := "COMPLETED"

	for *task.Status.Code != exportTaskCompleted {
		waitErrorCount++
		logInfo.Println("Current Status => ", *task.Status.Code)
		logInfo.Println("Waiting for task to be completed....")
		time.Sleep(time.Second * 15)
		describeExportTasksOutput, err = svcCloudwatchLogs.DescribeExportTasks(describeExportTasksInput)
		task = describeExportTasksOutput.ExportTasks[0]
		if waitErrorCount > 3 {
			logInfo.Println("Status => ", *task.Status.Code, "Error => ", err)
			return taskID, "FAILED", err
		}
	}
	return taskID, "SUCESS", err
}
