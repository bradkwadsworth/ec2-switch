// Package instance provides control primitives for EC2 instances
package instance

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// Action object for EC2 instance data
type Action struct {
	Conn     *ec2.EC2
	Filters  []*ec2.Filter
	IDs      []*string
	Name     string
	Verified bool
}

// Get slice of EC2 instance info pointer objects from ec2.Reservation objects
func info(res []*ec2.Reservation) []*ec2.Instance {
	inst := make([]*ec2.Instance, 0)
	for _, v := range res {
		for _, i := range v.Instances {
			// Append ec2.Instance objects to separate slice
			inst = append(inst, i)
		}
	}
	return inst
}

// Get slice of EC2 instance id string pointers from ec2.Instance objects
func ids(inst []*ec2.Instance) []*string {
	ids := make([]*string, len(inst))
	for i, v := range inst {
		// Append instanceId to new slice
		ids[i] = v.InstanceId
	}
	return ids
}

// Output info on EC2 instances that match filters and/or tags
func infoOutput(insts []*ec2.Instance) string {
	var str string
	for _, v := range insts {
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

// Output info on EC2 instance status
func statusOutput(status *ec2.InstanceStatus) string {
	var str string
	id := fmt.Sprintf("Instance ID: %s", *status.InstanceId)
	str += fmt.Sprintln(strings.Repeat("-", len(id)))
	str += fmt.Sprintln(id)
	str += fmt.Sprintf("Status: %s\n", *status.InstanceState.Name)
	return str
}

// Create new ec2.DescribeInstancesInput pointer object for describing EC2 instances
func newDescribeInput(filters []*ec2.Filter) *ec2.DescribeInstancesInput {
	input := new(ec2.DescribeInstancesInput)
	// Only include instances who's states are running or stopped
	st := "instance-state-name"
	stRun := "running"
	stStop := "stopped"
	filters = append(filters, &ec2.Filter{Name: &st, Values: []*string{&stRun, &stStop}})
	input.Filters = filters
	return input
}

// Create new ec2.DescribeInstanceStatusInput pointer object for describing EC2 instance statuses
func newDescribeStatus(insts []*string) *ec2.DescribeInstanceStatusInput {
	input := new(ec2.DescribeInstanceStatusInput)
	input.SetInstanceIds(insts)
	// Include all instances no matter what state they are in
	input.SetIncludeAllInstances(true)
	return input
}

// Create new ec2.StopInstancesInput pointer object for stopping EC2 instances
func newStopInput(ids []*string) *ec2.StopInstancesInput {
	input := new(ec2.StopInstancesInput)
	input.SetInstanceIds(ids)
	return input
}

// Create new ec2.StartInstancesInput pointer object for starting EC2 instances
func newStartInput(ids []*string) *ec2.StartInstancesInput {
	input := new(ec2.StartInstancesInput)
	input.SetInstanceIds(ids)
	return input
}

// NewAction creates a new Action object
func NewAction() *Action {
	return new(Action)
}

// SetIDs set the IDs field for an Action object with EC2 instance IDs
func (s *Action) SetIDs() error {
	query, err := s.Conn.DescribeInstances(newDescribeInput(s.Filters))
	if err != nil {
		return err
	}
	s.IDs = ids(info(query.Reservations))
	return nil
}

// List creates a list of EC2 instances that match desired Filters field
func (s *Action) List() error {
	query, err := s.Conn.DescribeInstances(newDescribeInput(s.Filters))
	if err != nil {
		return err
	}
	ids := info(query.Reservations)
	fmt.Println(infoOutput(ids))
	return nil
}

// Verify if desired EC2 instance action should be performed
func (s *Action) verifyAction() error {
	var verify string
	if s.Name == "" {
		return errors.New("No action defined")
	}
	if s.Verified == true {
		return nil
	}
	fmt.Printf("Are you sure you would like to %s the above instances (y/n)\n", s.Name)
	_, err := fmt.Scan(&verify)
	if err != nil {
		return err
	}
	switch verify {
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

// Start commands EC2 instances to start
func (s *Action) Start() error {
	if err := s.List(); err != nil {
		return err
	}
	if err := s.verifyAction(); err != nil {
		return err
	}
	if s.Verified {
		_, err := s.Conn.StartInstances(newStartInput(s.IDs))
		if err != nil {
			return err
		}
		if err := s.poll("running"); err != nil {
			return err
		}
	}
	return nil
}

// Stop commands EC2 instances to stop
func (s *Action) Stop() error {
	if err := s.List(); err != nil {
		return err
	}
	if err := s.verifyAction(); err != nil {
		return err
	}
	if s.Verified {
		_, err := s.Conn.StopInstances(newStopInput(s.IDs))
		if err != nil {
			return err
		}
		if err := s.poll("stopped"); err != nil {
			return err
		}
	}
	return nil
}

// Wait for EC2 instances to become desired state
func (s *Action) poll(reqState string) error {
	readyInstances := make([]string, len(s.IDs))
	//Get instances ids from state change
	for i := 0; i < len(readyInstances); {
		// Query api for updates to instance statuses
		insts, err := s.Conn.DescribeInstanceStatus(newDescribeStatus(s.IDs))
		if err != nil {
			return err
		}
		// When instance is in desired state add to readyInstances slice and increment counter
		for k := range insts.InstanceStatuses {
			if readyInstances[k] == *insts.InstanceStatuses[k].InstanceId {
				continue
			}
			if *insts.InstanceStatuses[k].InstanceState.Name == reqState {
				readyInstances[k] = *insts.InstanceStatuses[k].InstanceId
				fmt.Println(statusOutput(insts.InstanceStatuses[k]))
				i++
			} else {
				fmt.Println(statusOutput(insts.InstanceStatuses[k]))
			}
		}
		// Sleep between api calls
		time.Sleep(1 * time.Second)
	}
	return nil
}
