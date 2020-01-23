package main

// ensureBaseStackExists will verify that the deployment stack exists. If
// the stack does not exist, the stack will be created using the base template.
//func ensureBaseStackExists(a *adhesive.Adhesive) (string, error) {
//	out, err := a.Cfn.DescribeStacks(&cloudformation.DescribeStacksInput{
//		StackName: aws.String(a.Config.Workflow),
//	})
//
//	if err2, ok := err.(awserr.Error); ok {
//		fmt.Println("cloudformation error")
//		fmt.Println(err2.Code())
//		if err2.Code() == "ValidationError" {
//			fmt.Println("sadf")
//			coreTemplate := adhesive.GenerateCoreTemplate()
//			j, err := coreTemplate.JSON()
//			if err != nil {
//				return "", err
//			}
//
//			capabilities := []*string{
//				aws.String(cloudformation.CapabilityCapabilityIam),
//				aws.String(cloudformation.CapabilityCapabilityNamedIam),
//			}
//
//			// Create the initial stack from the core template.
//			out2, err := a.Cfn.CreateStack(&cloudformation.CreateStackInput{
//				Capabilities: capabilities,
//				StackName:    aws.String(a.Config.Workflow),
//				TemplateBody: aws.String(string(j)),
//			})
//
//			if err != nil {
//				return "", err
//			}
//
//			// Monitor the stack creation.
//			fmt.Printf("%s: Creating stack...\n", a.Config.Workflow)
//
//			if err := adhesive.MonitorStack(a.Cfn, *out2.StackId, a.Config.Workflow); err != nil {
//				return "", nil
//			}
//
//			fmt.Printf("%s: Creation complete...\n", a.Config.Workflow)
//
//			return *out2.StackId, nil
//		}
//		return "", err2
//	}
//	if err != nil {
//		return "", err
//	}
//
//	return *out.Stacks[0].StackId, nil
//}

//var (
//	deployCommand = &cobra.Command{
//		Use:   "deploy",
//		Short: "Deploy your AWS Glue jobs using CloudFormation.",
//		Run: func(cmd *cobra.Command, args []string) {
//			fmt.Println("deploy")
//		},
//	}
//)
//
//func deploy(ctx *cli.Context) error {
//	config, err := command.LoadConfigFile(ctx.GlobalString("config"))
//	if err != nil {
//		return err
//	}
//
//	a, err := adhesive.New(config)
//	if err != nil {
//		return err
//	}
//
//	return a.Deploy()
//}
