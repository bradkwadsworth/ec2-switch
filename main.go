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
	// Verify command argument
	action := flag.Arg(0)
	if err := checkArgs(action); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Create AWS session. Must provide credentials through environment variables
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	svc := ec2.New(sess)
	// Combine filter and tag flags
	allFilters := append(filters, tags...)
	// Get info about instances that match filters
	query, err := svc.DescribeInstances(newDescribeInstanceInput(allFilters))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(instanceOutput(query.Reservations))
	// Quit if requested action is list
	if action == "list" {
		os.Exit(0)
	}
	// Verify requested action
	fmt.Printf("Are you sure you would like to %s the above instances (y/n)\n", action)
	var verify string
	_, err = fmt.Scan(&verify)
	if err != nil {
		fmt.Println(err)
	}
	switch verify {
	case "y":
		fmt.Printf("Performing %s action on instances\n", action)
		switch action {
		case "stop":
			output, err := svc.StopInstances(newStopInstanceInput(query))
			if err != nil {
				fmt.Println(err)
			}
			if err := pollInstances(svc, output.StoppingInstances, "stopped"); err != nil {
				fmt.Println(err)
			}
		case "start":
			output, err := svc.StartInstances(newStartInstanceInput(query))
			if err != nil {
				fmt.Println(err)
			}
			if err := pollInstances(svc, output.StartingInstances, "running"); err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Printf("Action %s not defined", action)
		}
	case "n":
		fmt.Println("Exiting and taking no action")
	default:
		fmt.Println("Answer must be y or n")
	}
}
