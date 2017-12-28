package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	flag.Var(&tags, "tag", "Tags key:value. example MyTag:Value1,Value2")
	flag.Var(&filters, "filter", "Filters key:value. example MyKey:Value1,Value2")
	flag.Parse()
	action := flag.Arg(0)
	if action == "" {
		os.Exit(1)
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	svc := ec2.New(sess)
	allFilters := append(filters, tags...)
	query, err := svc.DescribeInstances(newDescribeInstanceInput(allFilters))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(instanceOutput(query.Reservations))
	if action == "list" {
		os.Exit(0)
	}
	fmt.Printf("Are you sure you would like to %s the above instances (y/n)\n", action)
	var verify string
	_, err = fmt.Scan(&verify)
	if err != nil {
		fmt.Println(err)
	}
	switch verify {
	case "y":
		fmt.Printf("Performing %s action on instances\n", action)
		if action == "stop" {
			output, err := svc.StopInstances(newStopInstanceInput(query))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(instanceStateOutput(output.StoppingInstances))
		} else if action == "start" {
			output, err := svc.StartInstances(newStartInstanceInput(query))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(instanceStateOutput(output.StartingInstances))
		} else {
			fmt.Printf("Action %s not defined", action)
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
