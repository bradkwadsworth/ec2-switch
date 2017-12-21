package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type ec2Filters []*ec2.Filter

func (s *ec2Filters) String() string {
	return fmt.Sprint(*s)
}

func (s *ec2Filters) Set(value string) error {
	strs := strings.Split(value, ":")
	filter := new(ec2.Filter)
	filter.SetName(strs[0])
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

func createFilters(f ec2Filters) []*ec2.Filter {
	fs := make([]*ec2.Filter, 0)
	// for k, v := range f {
	// 	flt := new(ec2.Filter)
	// 	flt.SetName(k)
	// 	flt.SetValues(v)
	// }
	// filter.SetName(name)
	// valuePnt := make([]*string, len(values))
	// for i := range values {
	// 	valuePnt[i] = &values[i]
	// }
	// filter.SetValues(valuePnt)
	return fs
}

func assembleFilters(filter ...*ec2.Filter) []*ec2.Filter {
	filters := make([]*ec2.Filter, 0)
	for _, i := range filter {
		filters = append(filters, i)
	}
	return filters
}

func getInstanceIds(instances *ec2.DescribeInstancesOutput) []*string {
	strPnt := make([]*string, 0)
	for _, v := range instances.Reservations {
		for _, v := range v.Instances {
			strPnt = append(strPnt, v.InstanceId)
		}
	}
	return strPnt
}

func main() {
	flag.Var(&filters, "tag", "Tags key:value. example MyKey:Value1,Value2")
	flag.Parse()
	fmt.Println(filters.String())
	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))
	// svc := ec2.New(sess)
	// envFilter := createFilter("tag:Environment", "alpha")
	// filters := assembleFilters(envFilter)
	// instanceInput := new(ec2.DescribeInstancesInput)
	// instanceInput.Filters = filters
	// instances, _ := svc.DescribeInstances(instanceInput)
	// instancesIds := getInstanceIds(instances)
	// fmt.Println(instancesIds)

	// stopInstances := new(ec2.StopInstancesInput)
	// stopInstances.SetInstanceIds(instances)
	// output, err := svc.StopInstances(stopInstances)
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// for _, s := range output.StoppingInstances {
	// 	fmt.Printf("%s %s", *s.InstanceId, s.CurrentState.String())
	// }
}
