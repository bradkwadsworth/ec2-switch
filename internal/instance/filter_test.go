package instance

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var tags ec2Tags = ec2Tags{
	&ec2.Filter{
		Name:   aws.String("foo"),
		Values: aws.StringSlice([]string{"bar", "baz"}),
	},
	&ec2.Filter{
		Name:   aws.String("stooges"),
		Values: aws.StringSlice([]string{"Larry", "Moe", "Curly", "Shemp"}),
	},
}

var filters ec2Filters = ec2Filters{
	&ec2.Filter{
		Name:   aws.String("foo"),
		Values: aws.StringSlice([]string{"bar", "baz"}),
	},
	&ec2.Filter{
		Name:   aws.String("stooges"),
		Values: aws.StringSlice([]string{"Larry", "Moe", "Curly", "Shemp"}),
	},
}

func TestEc2TagsString(t *testing.T) {
	w := `{
  Name: "foo",
  Values: ["bar","baz"]
},{
  Name: "stooges",
  Values: [
    "Larry",
    "Moe",
    "Curly",
    "Shemp"
  ]
}`
	if s := fmt.Sprint(tags); s != w {
		t.Errorf("%s does not equal %s", s, w)
	}
}

func TestEc2TagsSet(t *testing.T) {
	w := ec2Tags{
		&ec2.Filter{
			Name:   aws.String("tag:foo"),
			Values: aws.StringSlice([]string{"bar", "baz"}),
		},
		&ec2.Filter{
			Name:   aws.String("tag:stooges"),
			Values: aws.StringSlice([]string{"Larry", "Moe", "Curly", "Shemp"}),
		},
	}
	tags := make(ec2Tags, 0)
	for _, v := range []string{"foo:bar,baz", "stooges:Larry,Moe,Curly,Shemp"} {
		tags.Set(v)
	}
	if reflect.DeepEqual(tags, w) {
		t.Errorf("%s does not equal %s", tags, w)
	}
}

func TestEc2FiltersString(t *testing.T) {
	w := `{
  Name: "foo",
  Values: ["bar","baz"]
},{
  Name: "stooges",
  Values: [
    "Larry",
    "Moe",
    "Curly",
    "Shemp"
  ]
}`
	if s := fmt.Sprint(filters); s != w {
		t.Errorf("%s does not equal %s", s, w)
	}
}

func TestEc2FiltersSet(t *testing.T) {
	w := ec2Tags{
		&ec2.Filter{
			Name:   aws.String("tag:foo"),
			Values: aws.StringSlice([]string{"bar", "baz"}),
		},
		&ec2.Filter{
			Name:   aws.String("stooges"),
			Values: aws.StringSlice([]string{"Larry", "Moe", "Curly", "Shemp"}),
		},
	}
	filters := make(ec2Filters, 0)
	for _, v := range []string{"tag:foo:bar,baz", "stooges:Larry,Moe,Curly,Shemp"} {
		filters.Set(v)
	}
	if reflect.DeepEqual(filters, w) {
		t.Errorf("%s does not equal %s", filters, w)
	}
}
