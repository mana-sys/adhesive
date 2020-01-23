package adhesive

//
//import (
//	"fmt"
//	"time"
//
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/aws/endpoints"
//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/cloudformation"
//	"github.com/urfave/cli"
//)
//
//func (a *Adhesive) remove(ctx *cli.Context) error {
//	return nil
//}
//
//func remove(ctx *cli.Context) error {
//	sess, err := session.NewSession(&aws.Config{
//		Region: aws.String(endpoints.UsWest2RegionID),
//	})
//
//	if err != nil {
//		return err
//	}
//
//	svc := cloudformation.New(sess)
//	_, err = svc.DeleteStack(&cloudformation.DeleteStackInput{
//		StackName: aws.String("adhesive-hello-world"),
//	})
//
//	fmt.Println("adhesive-hello-world: Deleting stack...")
//
//	start := time.Now().Round(time.Second)
//	ticker := time.NewTicker(5 * time.Second)
//
//loop:
//	for {
//		select {
//		case t := <-ticker.C:
//			// Check if stack is done creating or if an error has occurred.
//			out2, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
//				StackName: aws.String("adhesive-hello-world"),
//			})
//			if err != nil {
//				return err
//			}
//
//			if *out2.Stacks[0].StackStatus != cloudformation.StackStatusDeleteInProgress {
//				break loop
//			}
//
//			fmt.Printf("adhesive-hello-world: Still deleting stack... (%s elapsed)\n", t.Round(time.Second).Sub(start))
//		}
//	}
//
//	return err
//}
