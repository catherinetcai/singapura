package singapura

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/tuvistavie/securerandom"
)

type IAM struct {
	Iam  *iam.IAM
	sess *session.Session
}

const passwordLen = 10

//IamInstance creates an instance of a Iam Client with the default ~/.aws/credentials
func IamInstance(profile string) (*IAM, error) {
	sess, err := awsSession(profile)
	if err != nil {
		return nil, err
	}
	i := &IAM{
		Iam:  iam.New(sess),
		sess: sess,
	}
	return i, nil
}

func awsSession(profile string) (*session.Session, error) {
	var opts session.Options
	var sess *session.Session
	var err error
	if profile != "" {
		opts = session.Options{
			Profile: profile,
		}
		sess, err = session.NewSessionWithOptions(opts)
	} else {
		sess, err = session.NewSession()
	}
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (i *IAM) CreateUser(username *string) (*iam.CreateUserOutput, error) {
	u := &iam.CreateUserInput{
		UserName: username,
	}
	res, err := i.Iam.CreateUser(u)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *IAM) CreateUserPassword(username *string) (*iam.CreateLoginProfileOutput, error) {
	var res *iam.CreateLoginProfileOutput
	var err error
	var password string
	password, err = securerandom.Base64(passwordLen, true)
	fmt.Printf("Password: %v\n", password)
	if err != nil {
		return nil, err
	}
	p := &iam.CreateLoginProfileInput{
		UserName:              username,
		Password:              aws.String(password),
		PasswordResetRequired: aws.Bool(true),
	}
	res, err = i.Iam.CreateLoginProfile(p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *IAM) AddUserGroups() {
	*iam.AddUserToGroupInput
	i.Iam.AddUserToGroup
}
