package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// Slice of *ec2.Filter types
type ec2Tags []*ec2.Filter

// Returns string of ec2Tags
func (s *ec2Tags) String() string {
	return fmt.Sprint(*s)
}

// Combine multiple tag flags
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

// Slice of ec2.Filter types
type ec2Filters []*ec2.Filter

// Returns string of ec2Filters
func (s *ec2Filters) String() string {
	return fmt.Sprint(*s)
}

// Combine multiple filter flags
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

var actions = []string{"list", "start", "stop"}

// Check command arguments
func checkArgs(arg string) error {
	for _, v := range actions {
		if v == arg {
			return nil
		}
	}
	return errors.New("Specified action not defined")
}

// Get slice of ec2.Instance pointer objects from ec2.Reservation objects
func instances(res []*ec2.Reservation) []*ec2.Instance {
	inst := make([]*ec2.Instance, 0)
	for _, v := range res {
		for _, i := range v.Instances {
			// Append ec2.Instance objects to separate slice
			inst = append(inst, i)
		}
	}
	return inst
}

// Get slice of InstanceId string pointers from ec2.Instance objects
func instanceIds(inst []*ec2.Instance) []*string {
	ids := make([]*string, len(inst))
	for i, v := range inst {
		// Append instanceId to new slice
		ids[i] = v.InstanceId
	}
	return ids
}

// Output info on instances that match filters and/or tags
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

// Output info on instance status
func instanceStatusOutput(status *ec2.InstanceStatus) string {
	var str string
	id := fmt.Sprintf("Instance ID: %s", *status.InstanceId)
	str += fmt.Sprintln(strings.Repeat("-", len(id)))
	str += fmt.Sprintln(id)
	str += fmt.Sprintf("Status: %s\n", *status.InstanceState.Name)
	return str
}

// Wait for instances to become desired state
func pollInstances(conn *ec2.EC2, states []*ec2.InstanceStateChange, reqState string) error {
	instanceIds := make([]*string, len(states))
	readyInstances := make([]*string, len(states))
	//Get instances ids from state change
	for i := range states {
		instanceIds[i] = states[i].InstanceId
	}
	for i := 0; i < len(states); {
		if readyInstances[i] == states[i].InstanceId {
			continue
		}
		// Query api for updates to instance statuses
		instances, err := conn.DescribeInstanceStatus(newDescribeInstanceStatus(instanceIds))
		if err != nil {
			return err
		}
		// When instance is in desired state add to readyInstances slice and increment counter
		if *instances.InstanceStatuses[i].InstanceState.Name == reqState {
			readyInstances[i] = instanceIds[i]
			fmt.Println(instanceStatusOutput(instances.InstanceStatuses[i]))
			i++
		} else {
			fmt.Println(instanceStatusOutput(instances.InstanceStatuses[i]))
		}
		// Sleep between api calls
		time.Sleep(1 * time.Second)
	}
	return nil
}

// Create new ec2.DescribeInstancesInput pointer object
func newDescribeInstanceInput(filters []*ec2.Filter) *ec2.DescribeInstancesInput {
	input := new(ec2.DescribeInstancesInput)
	// Only include instances who's states are running or stopped
	st := "instance-state-name"
	stRun := "running"
	stStop := "stopped"
	filters = append(filters, &ec2.Filter{Name: &st, Values: []*string{&stRun, &stStop}})
	input.Filters = filters
	return input
}

// Create new ec2.DescribeInstanceStatusInput pointer object
func newDescribeInstanceStatus(instances []*string) *ec2.DescribeInstanceStatusInput {
	input := new(ec2.DescribeInstanceStatusInput)
	input.SetInstanceIds(instances)
	// Include all instances no matter what state they are in
	input.SetIncludeAllInstances(true)
	return input
}

// Create new ec2.StopInstancesInput pointer object
func newStopInstanceInput(query *ec2.DescribeInstancesOutput) *ec2.StopInstancesInput {
	input := new(ec2.StopInstancesInput)
	input.SetInstanceIds(instanceIds(instances(query.Reservations)))
	return input
}

// Create new ec2.StartInstancesInput pointer object
func newStartInstanceInput(query *ec2.DescribeInstancesOutput) *ec2.StartInstancesInput {
	input := new(ec2.StartInstancesInput)
	input.SetInstanceIds(instanceIds(instances(query.Reservations)))
	return input
}
