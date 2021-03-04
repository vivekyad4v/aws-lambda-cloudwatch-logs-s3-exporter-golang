package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/s3"
)

func s3ExportTaskStart(cloudwatchLogGroupNames []string, getDate getDate) (taskIDList []string) {

	svcCloudwatchLogs := initAwsCloudwatchLogsSession()

	var destinationPrefix string = getDate.year + "/" + getDate.month + "/" + getDate.day
	var destinationKey string = getDate.year + "-" + getDate.month + "-" + getDate.day + ".json"

	type exportTaskIDList struct {
		TaskIDs []string
	}

	for k := range cloudwatchLogGroupNames {
		createExportTaskInput := &cloudwatchlogs.CreateExportTaskInput{
			Destination:       &logS3BucketID,
			DestinationPrefix: &destinationPrefix,
			From:              &getDate.yms,
			To:                &getDate.tms,
			LogGroupName:      &cloudwatchLogGroupNames[k],
		}
		fmt.Println("Exporting log group =>", cloudwatchLogGroupNames[k])
		createExportTaskOutput, _ := svcCloudwatchLogs.CreateExportTask(createExportTaskInput)
		taskID, Status, err := waitForExportTaskCompletion(*createExportTaskOutput.TaskId)
		fmt.Println("Executed Task ID => ", taskID, ", Status => ", Status, ", Error => ", err)
		taskIDList = append(taskIDList, *createExportTaskOutput.TaskId)
		fmt.Println("Task Id list - ", taskIDList)
	}

	// PUT TASKIDs to S3 BUCKET
	fmt.Println("Putting taskIDs to S3 bucket....")

	// Converting to JSON
	taskIDToJSON := exportTaskIDList{TaskIDs: taskIDList}
	taskIDListJSON, _ := json.Marshal(taskIDToJSON)
	readerJSON := strings.NewReader(string(taskIDListJSON))

	// Creating S3 session
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Error - Unable to create s3 AWS session!")
	}
	s3svc := s3.New(sess)

	// Putting object to S3
	s3PutObjectInput := &s3.PutObjectInput{
		Bucket: &logS3BucketID,
		Key:    &destinationKey,
		Body:   readerJSON,
	}

	_, err = s3svc.PutObject(s3PutObjectInput)

	return taskIDList
}

func waitForExportTaskCompletion(taskID string) (string, string, error) {
	svcCloudwatchLogs := initAwsCloudwatchLogsSession()

	describeExportTasksInput := &cloudwatchlogs.DescribeExportTasksInput{
		TaskId: &taskID,
	}

	time.Sleep(time.Second * 1)
	describeExportTasksOutput, err := svcCloudwatchLogs.DescribeExportTasks(describeExportTasksInput)
	task := describeExportTasksOutput.ExportTasks[0]
	fmt.Println("Current Task ID => ", *task.TaskId, ", Status => ", *task.Status.Code)

	waitErrorCount := 0
	exportTaskCompleted := "COMPLETED"

	for *task.Status.Code != exportTaskCompleted {
		waitErrorCount++
		fmt.Println("Current Status => ", *task.Status.Code)
		fmt.Println("Waiting for task to be completed....")
		time.Sleep(time.Second * 5)
		describeExportTasksOutput, err = svcCloudwatchLogs.DescribeExportTasks(describeExportTasksInput)
		task = describeExportTasksOutput.ExportTasks[0]
		if waitErrorCount > 3 {
			fmt.Println("Status => ", *task.Status.Code, "Error => ", err)
			return taskID, "FAILED", err
		}
	}
	return taskID, "SUCESS", err
}
