package instance

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type InstanceActions struct {
	Conn      *ec2.EC2
	Filters   []*ec2.Filter
	Instances []*string
	Action    string
	Verified  bool
}

func NewInstanceActions() *InstanceActions {
	return new(InstanceActions)
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
func (s *InstanceActions) pollInstances(reqState string) error {
	readyInstances := make([]*string, len(s.Instances))
	//Get instances ids from state change
	for i := 0; i < len(readyInstances); {
		if readyInstances[i] == s.Instances[i] {
			continue
		}
		// Query api for updates to instance statuses
		instances, err := s.Conn.DescribeInstanceStatus(newDescribeInstanceStatus(s.Instances))
		if err != nil {
			return err
		}
		// When instance is in desired state add to readyInstances slice and increment counter
		if *instances.InstanceStatuses[i].InstanceState.Name == reqState {
			readyInstances[i] = s.Instances[i]
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
func NewDescribeInstanceInput(filters []*ec2.Filter) *ec2.DescribeInstancesInput {
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
func newStartInstanceInput(ids []*string) *ec2.StartInstancesInput {
	input := new(ec2.StartInstancesInput)
	input.SetInstanceIds(ids)
	return input
}

func (s *InstanceActions) SetInstanceIds() {
	query, err := s.Conn.DescribeInstances(NewDescribeInstanceInput(s.Filters))
	s.Instances = instanceIds(instances(query.Reservations))
}

func (s *InstanceActions) ListInstances() error {
	query, err := s.Conn.DescribeInstances(NewDescribeInstanceInput(s.Filters))
	if err != nil {
		return err
	}
	fmt.Println(instanceOutput(query.Reservations))
	return nil
}

func (s *InstanceActions) verifyAction() error {
	if s.Action == "" {
		return errors.New("No action defined")
	}
	fmt.Printf("Are you sure you would like to %s the above instances (y/n)\n", s.Action)
	_, err := fmt.Scan(&verify)
	if err != nil {
		return err
	}
	switch s.Action {
	case "y":
		s.Verified = true
		return nil
	case "n":
		s.Verified = false
		return nil
	default:
		return errors.New("Answer must be y or n")
	}
}

func (s *InstanceActions) StartInstances() error {
	if err := s.verifyAction(); err != nil {
		return err
	}
	if s.Verified {
		output, err := s.Conn.StartInstances(newStartInstanceInput(s.Instances))
		if err != nil {
			return err
		}
	}
	if err := s.pollInstances("running"); err != nil {
		return err
	}
	return nil
}

func (s *InstanceActions) StopInstances() error {
	output, err := conn.StartInstances(newStartInstanceInput(query))
	if err != nil {
		return err
	}
	if err := pollInstances(conn, output.StartingInstances, "stopped"); err != nil {
		return err
	}
	return nil
}
