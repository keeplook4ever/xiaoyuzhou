package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	S3userId  = "AKIA2KPKPXCU2E5423ML"
	S3userKey = "RW8hOcFnkb1wGPUflF9HfwGmVStkk3FmYvllqFli"
	ValidTime = 3600 // 1小时
)

func GetToken() (*TmpTokenStruct, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(endpoints.ApNortheast1RegionID),
		Credentials: credentials.NewStaticCredentials(S3userId, S3userKey, ""),
	}))

	//roleARN := "arn:aws:iam::709698238633:user/s3-xiaoyuzhou"
	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the "myRoleARN" ARN.
	//creds := stscreds.NewCredentials(sess, roleARN)

	svc := sts.New(sess)

	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(ValidTime),
	}
	res, err := svc.GetSessionToken(input)
	if err != nil {
		return nil, err
	}

	resp := TmpTokenStruct{
		AccessKeyId:     *res.Credentials.AccessKeyId,
		Expiration:      res.Credentials.Expiration.Unix(),
		SecretAccessKey: *res.Credentials.SecretAccessKey,
		SessionToken:    *res.Credentials.SessionToken,
	}
	return &resp, err
}

type TmpTokenStruct struct {
	AccessKeyId     string `json:"access_key_id"`
	Expiration      int64  `json:"expiration"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`
}
