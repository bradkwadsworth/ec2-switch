package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ec2Filters []*ec2.Filter

func (s *ec2Filters) String() string {
	return fmt.Sprint(*s)
}

func (s *ec2Filters) Set(value string) error {
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

var filters ec2Filters

type instances []*ec2.Instance

func NewInstances(instances *ec2.DescribeInstancesOutput) instances {
	inst := make([]*ec2.Instance, 0)
	for _, v := range instances.Reservations {
		for _, v := range v.Instances {
			inst = append(inst, v)
		}
	}
	return inst
}

func (s *instances) InstanceIds() []string {
	ids := make([]string, len(*s))
	for i, v := range *s {
		ids[i] = v.GoString()
	}
	return ids
}

func tagOutput(tag []*ec2.Tag) string {
	var str string
	for _, v := range tag {
		str += fmt.Sprintf("  %s: %s\n", *v.Key, *v.Value)
	}
	return str
}

func (s *instances) shortOutput() {
	num := fmt.Sprintf("Found %d instances\n", len(*s))
	fmt.Println(strings.Repeat("-", len(num)))
	fmt.Print(num)
	for _, v := range *s {
		name := fmt.Sprintf("Instance ID: %s\n", *v.InstanceId)
		tags := fmt.Sprintf("%s", tagOutput(v.Tags))
		fmt.Println(strings.Repeat("-", len(name)))
		fmt.Print(name)
		fmt.Println("Tags:")
		fmt.Printf("%s\n", tags)
	}
}

type instanceStates []*ec2.InstanceStateChange

func NewInstanceStates(states *ec2.StopInstancesOutput) instanceStates {
	return states.StoppingInstances
}

func (s *instanceStates) InstanceIds() []string {
	inst := make([]string, len(*s))
	for i, v := range *s {
		inst[i] = *v.InstanceId
	}
	return inst
}

func main() {
	flag.Var(&filters, "tag", "Tags key:value. example MyKey:Value1,Value2")
	flag.Parse()
	action := flag.Arg(0)
	if action == "" {
		os.Exit(1)
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	svc := ec2.New(sess)
	instanceInput := new(ec2.DescribeInstancesInput)
	instanceInput.Filters = filters
	query, err := svc.DescribeInstances(instanceInput)
	if err != nil {
		fmt.Println(err)
	}
	insts := NewInstances(query)
	insts.shortOutput()
	if action == "list" {
		os.Exit(0)
	}
	fmt.Printf("Are you sure you would like to %s the above instances (y|n)\n", action)
	var verify string
	_, err = fmt.Scan(&verify)
	if err != nil {
		fmt.Println(err)
	}
	switch verify {
	case "y":
		fmt.Printf("Performing %s action on instances\n", action)
		if action == "stop" {
			// 	stopInstances := new(ec2.StopInstancesInput)
			// 	stopInstances.SetInstanceIds(insts.InstanceIds())
			// 	output, err := svc.StopInstances(stopInstances)
			// 	if err != nil {
			// 		fmt.Print(err)
			// 	}
		}
	case "n":
		fmt.Println("Exiting and taking no action")
	default:
		fmt.Println("Answer must be y or n")
	}

	// for _, s := range output.StoppingInstances {
	// 	fmt.Printf("%s %s", *s.InstanceId, s.CurrentState.String())
	// }
}
