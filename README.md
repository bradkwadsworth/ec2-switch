# EC2 Switch

EC2 Switch is a utility that simplifies starting up or shutting down Amazon Web Services EC2 instances based on typical EC2 instance specific filters.

## Dependencies
EC2 Switch uses aws-sdk-go, therefore AWS credentials must either be sourced with AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY, or AWS_PROFILE environment variables.

## Installation
### Download Binary
TODO
### Build from Source
#### Get dependencies
`go get github.com/aws/aws-sdk-go/service/ec2`

#### Build
`go build`
#### Install
`go install`

Copy the ec2-switch binary somewhere in your path

## List Instances by Tag
`ec2-switch -tag Environment:myenv list`

## List Instances by Instance-Id Filter
`ec2-switch -filter instance-id:someid list`

## Start Instances by Tags
`ec2-switch -tag Environment:myenv -Type:myspecialserver start`

## Stop Instances by Tags and Filters
`ec2-switch -tag Environment:myenv -filter instance-id:someid stop`
