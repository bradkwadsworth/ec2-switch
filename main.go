package main

import (
	"flag"
	"fmt"
	"os"

	"git.digitalmeasures.com/devops/ec2-switch/internal/instance"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	flag.Var(&instance.Tags, "tag", "Tags key:value. example MyTag:Value1,Value2")
	flag.Var(&instance.Filters, "filter", "Filters key:value. example MyKey:Value1,Value2")
	force := flag.Bool("force", false, "Do not ask to verify action")
	flag.Parse()
	insts := instance.NewInstanceActions()
	// Verify command argument
	action := flag.Arg(0)
	if err := instance.CheckArgs(action); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	insts.Action = action
	// Create AWS session. Must provide credentials through environment variables
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	svc := ec2.New(sess)
	insts.Conn = svc
	// Combine filter and tag flags
	allFilters := append(instance.Filters, instance.Tags...)
	insts.Filters = allFilters
	insts.SetInstanceIds()
	switch action {
	case "list":
		if err := insts.ListInstances(); err != nil {
			fmt.Println(err)
		}
	case "start":
		verified, err := instance.VerifyAction("start")
		if err != nil {
			fmt.Println(err)
		}
		insts.Verified = verified
		// 	if verified {
		// 		if err := instance.StartInstances(svc, query); err != nil {
		// 			fmt.Println(err)
		// 		}
		// 	}
		// case "stop":
		// 	if err := instance.StopInstances(svc, query); err != nil {
		// 		fmt.Println(err)
		// 	}
	}
	// // Verify requested action
	//
	// switch verify {
	// case "y":
	// 	fmt.Printf("Performing %s action on instances\n", action)
	// 	switch action {
	// 	case "stop":
	// 		output, err := svc.StopInstances(newStopInstanceInput(query))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		if err := pollInstances(svc, output.StoppingInstances, "stopped"); err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	case "start":
	// 		output, err := svc.StartInstances(newStartInstanceInput(query))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		if err := pollInstances(svc, output.StartingInstances, "running"); err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	default:
	// 		fmt.Printf("Action %s not defined", action)
	// 	}
	// case "n":
	// 	fmt.Println("Exiting and taking no action")
	// default:
	// 	fmt.Println("Answer must be y or n")
	// }
}
