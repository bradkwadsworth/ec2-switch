package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type ec2Tags []*ec2.Filter

func (s *ec2Tags) String() string {
	return fmt.Sprint(*s)
}

func (s *ec2Tags) Set(value string) error {
	strs := strings.Split(value, ":")
	filter := new(ec2.Filter)
	filter.SetName("tag:" + strs[0])
	vals := strings.Split(strs[1], ",")
	filterVals := make([]*string, len(vals))
	for i := range vals {
		filterVals[i] = &vals[i]
	}
	filter.SetValues(filterVals)
	*s = append(*s, filter)
	return nil
}

type ec2Filters []*ec2.Filter

func (s *ec2Filters) String() string {
	return fmt.Sprint(*s)
}

func (s *ec2Filters) Set(value string) error {
	strs := strings.Split(value, ":")
	filter := new(ec2.Filter)
	var vals []string
	if strs[0] == "tag" {
		filter.SetName(strs[0] + ":" + strs[1])
		vals = strings.Split(strs[2], ",")
	} else {
		filter.SetName(strs[0])
		vals = strings.Split(strs[1], ",")
	}
	filterVals := make([]*string, len(vals))
	for i := range vals {
		filterVals[i] = &vals[i]
	}
	filter.SetValues(filterVals)
	*s = append(*s, filter)
	return nil
}

var tags ec2Tags
var filters ec2Filters

func instances(res []*ec2.Reservation) []*ec2.Instance {
	inst := make([]*ec2.Instance, 0)
	for _, v := range res {
		for _, i := range v.Instances {
			inst = append(inst, i)
		}
	}
	return inst
}

func instanceIds(inst []*ec2.Instance) []*string {
	ids := make([]*string, len(inst))
	for i, v := range inst {
		ids[i] = v.InstanceId
	}
	return ids
}

func instanceOutput(res []*ec2.Reservation) string {
	var str string
	inst := instances(res)
	for _, v := range inst {
		id := fmt.Sprintf("Instance ID: %s\n", *v.InstanceId)
		str += fmt.Sprintln(strings.Repeat("-", len(id)))
		str += fmt.Sprintf(id)
		str += fmt.Sprintf("State: %s\n", *v.State.Name)
		str += fmt.Sprintln("Tags")
		for _, t := range v.Tags {
			str += fmt.Sprintf("  %s: %s\n", *t.Key, *t.Value)
		}
	}
	return str
}

func instanceStateOutput(states []*ec2.InstanceStateChange) string {
	var str string
	for _, v := range states {
		id := fmt.Sprintf("Instance ID: %s\n", *v.InstanceId)
		str += fmt.Sprintln(strings.Repeat("-", len(id)))
		str += fmt.Sprintln(id)
		str += fmt.Sprintf("  Last State: %s\n", *v.PreviousState.Name)
		str += fmt.Sprintf("  Current State: %s\n", *v.CurrentState.Name)
	}
	return str
}

func newDescribeInstanceInput(filters []*ec2.Filter) *ec2.DescribeInstancesInput {
	input := new(ec2.DescribeInstancesInput)
	st := "instance-state-name"
	stRun := "running"
	stStop := "stopped"
	filters = append(filters, &ec2.Filter{Name: &st, Values: []*string{&stRun, &stStop}})
	input.Filters = filters
	return input
}

func newStopInstanceInput(query *ec2.DescribeInstancesOutput) *ec2.StopInstancesInput {
	input := new(ec2.StopInstancesInput)
	input.SetInstanceIds(instanceIds(instances(query.Reservations)))
	return input
}

func newStartInstanceInput(query *ec2.DescribeInstancesOutput) *ec2.StartInstancesInput {
	input := new(ec2.StartInstancesInput)
	input.SetInstanceIds(instanceIds(instances(query.Reservations)))
	return input
}
