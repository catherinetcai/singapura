package singapura

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

//IamInstance creates an instance of a Iam Client with the default ~/.aws/credentials
func IamInstance() *iam.IAM {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Failed to create IAM session,", err)
		return nil
	}
	return iam.New(sess)
}
