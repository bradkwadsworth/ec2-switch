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

func Instances(res []*ec2.Reservation) []*ec2.Instance {
	inst := make([]*ec2.Instance, 0)
	for _, v := range res {
		for _, i := range v.Instances {
			inst = append(inst, i)
		}
	}
	return inst
}

func InstanceIds(inst []*ec2.Instance) []*string {
	ids := make([]*string, len(inst))
	for i, v := range inst {
		ids[i] = v.InstanceId
	}
	return ids
}

func InstanceOutput(res []*ec2.Reservation) string {
	var str string
	inst := Instances(res)
	for _, v := range inst {
		id := fmt.Sprintf("Instance ID: %s\n", *v.InstanceId)
		str += fmt.Sprintln(strings.Repeat("-", len(id)))
		str += fmt.Sprintf(id)
		str += fmt.Sprintln("Tags")
		for _, t := range v.Tags {
			str += fmt.Sprintf("  %s: %s\n", *t.Key, *t.Value)
		}
	}
	return str
}

func InstanceStateOutput(states []*ec2.InstanceStateChange) string {
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
	fmt.Println(InstanceOutput(query.Reservations))
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
			stopInstances := new(ec2.StopInstancesInput)
			stopInstances.SetInstanceIds(InstanceIds(Instances(query.Reservations)))
			output, err := svc.StopInstances(stopInstances)
			if err != nil {
				fmt.Print(err)
			}
			fmt.Println(InstanceStateOutput(output.StoppingInstances))
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
