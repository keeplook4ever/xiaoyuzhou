package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

func Test() {
	//sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	//if err != nil {
	//	fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
	//	fmt.Println(err)
	//	return
	//}
	//r := sdkConfig.Region
	//fmt.Println(r)
	//s3Client := s3.NewFromConfig(sdkConfig)
	//count := 10
	//fmt.Printf("Let's list up to %v buckets for your account.\n", count)
	//result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	//if err != nil {
	//	fmt.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
	//	return
	//}
	//if len(result.Buckets) == 0 {
	//	fmt.Println("You don't have any buckets!")
	//} else {
	//	if count > len(result.Buckets) {
	//		count = len(result.Buckets)
	//	}
	//	for _, bucket := range result.Buckets[:count] {
	//		fmt.Printf("\t%v\n", *bucket.Name)
	//	}
	//}

	// snippet-start:[sts.go.take_role.args]
	roleARN := "arn:aws:iam::709698238633:user/s3-xiaoyuzhou"
	sessionName := "s3token"

	// snippet-start:[sts.go.take_role.session]
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sts.New(sess)
	// snippet-end:[sts.go.take_role.session]

	result, err := TakeRole(svc, &roleARN, &sessionName)
	if err != nil {
		fmt.Println("Got an error assuming the role:")
		fmt.Println(err)
		return
	}

	// snippet-start:[sts.go.take_role.display]
	fmt.Println(result.AssumedRoleUser)
	credentials := *result.Credentials
	fmt.Println(credentials)
	// snippet-end:[sts.go.take_role.display]
}

// TakeRole gets temporary security credentials to access resources
// Inputs:
//     svc is an AWS STS service client
//     roleARN is the Amazon Resource Name (ARN) of the role to assume
//     sessionName is a unique identifier for the session
// Output:
//     If success, information about the assumed role and nil
//     Otherwise, nil and an error from the call to AssumeRole
func TakeRole(svc stsiface.STSAPI, roleARN, sessionName *string) (*sts.AssumeRoleOutput, error) {
	// snippet-start:[sts.go.take_role.call]
	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         roleARN,
		RoleSessionName: sessionName,
	})
	// snippet-end:[sts.go.take_role.call]

	return result, err
}
