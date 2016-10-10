package singapura

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/tuvistavie/securerandom"
)

type IAM struct {
	Iam  *iam.IAM
	sess *session.Session
}

type GroupConfig struct {
	Role string
	Env  string
}

type EnvRoles struct {
	//	Groups map[string]map[string][]string
	Envs map[string]Env
}

type Env struct {
	Roles map[string][]string
}

const (
	passwordLen    = 10
	defaultEnv     = "preprod"
	defaultRole    = "developer"
	groupsFileName = "configs/groups.yaml"
)

// IamInstance creates an instance of a Iam Client with the default ~/.aws/credentials
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

// awsSession creates an AWS session. If there is a profile passed in, it'll create a session
// with that option. Otherwise, it defaults to the default in ~/.aws/credentials
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

// CreateUser creates a user with the specified username
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

// CreateUserPassword generates a random password with the securerandom lib
// for the specified username
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

// AddUserGroups adds groups to a specified user
func (i *IAM) AddUserGroups(username string, g *GroupConfig) (*iam.AddUserToGroupOutput, error) {
	return nil, nil
}

// RoleByNameAndEnv returns a list of roles based off of environment and passed in role
func RoleByNameAndEnv(g *GroupConfig) ([]string, error) {
	setDefaultEnvRole(g)
	filename, _ := filepath.Abs(groupsFileName)
	file, err := ioutil.ReadFile(filename)

	var envRoles EnvRoles
	err = yaml.Unmarshal(file, &envRoles)
	fmt.Println(envRoles)

	if err != nil {
		return nil, err
	}
	return []string{}, err
}

// setDefaultEnvRole sets default environment and roles if not set
func setDefaultEnvRole(g *GroupConfig) {
	if g.Env == "" {
		g.Env = defaultEnv
	}
	if g.Role == "" {
		g.Role = defaultRole
	}
}
