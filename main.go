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
	insts := instance.NewAction()
	// Verify command argument
	action := flag.Arg(0)
	if err := instance.CheckArgs(action); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	insts.Name = action
	if *force {
		insts.Verified = true
	}
	// Create AWS session. Must provide credentials through environment variables
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	svc := ec2.New(sess)
	insts.Conn = svc
	// Combine filter and tag flags
	allFilters := append(instance.Filters, instance.Tags...)
	insts.Filters = allFilters
	insts.SetIDs()
	switch action {
	case "list":
		if err := insts.List(); err != nil {
			fmt.Println(err)
		}
	case "start":
		if err := insts.Start(); err != nil {
			fmt.Println(err)
		}
	case "stop":
		if err := insts.Stop(); err != nil {
			fmt.Println(err)
		}
	}
}
