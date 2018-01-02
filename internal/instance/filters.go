package instance

import (
	"fmt"
	"strings"

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

// Tags is a slice of *ec2.Filter. Used as a shortcut to EC2 instance tags
var Tags ec2Tags

// Filters is a slice of  *ec2.Filter
var Filters ec2Filters
