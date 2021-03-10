package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func initAwsCloudwatchLogsSession() (svcCloudwatchLogs *cloudwatchlogs.CloudWatchLogs) {
	mySession := session.Must(session.NewSession())
	return cloudwatchlogs.New(mySession)
}

func fetchCloudwatchLogGroups(logGroupNamePrefixes []string) (cloudwatchLogGroupNames []string) {
	svcCloudwatchLogs := initAwsCloudwatchLogsSession()

	for k, v := range logGroupNamePrefixes {
		input := &cloudwatchlogs.DescribeLogGroupsInput{Limit: &limitList, LogGroupNamePrefix: &v}
		describeLogGroups, err := svcCloudwatchLogs.DescribeLogGroups(input)

		for describeLogGroups.NextToken != nil {
			newInput := &cloudwatchlogs.DescribeLogGroupsInput{Limit: &limitList, LogGroupNamePrefix: &logGroupNamePrefixes[k], NextToken: describeLogGroups.NextToken}
			describeLogGroups, err = svcCloudwatchLogs.DescribeLogGroups(newInput)
			cloudwatchLogGroupNames = append(cloudwatchLogGroupNames, *describeLogGroups.LogGroups[0].LogGroupName)

			if err != nil {
				logInfo.Println("Error in - ", describeLogGroups.LogGroups, "-", err)
			}
		}

		if len(describeLogGroups.LogGroups) < 1 {
			logInfo.Println("No match for log group prefix - ", v)
		}

		if len(describeLogGroups.LogGroups) > 0 {
			cloudwatchLogGroupNames = append(cloudwatchLogGroupNames, *describeLogGroups.LogGroups[0].LogGroupName)
		}
	}

	return cloudwatchLogGroupNames
}
