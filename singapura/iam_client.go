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

// The Singapura struct keeps track of all the user's input flags
type Singapura struct {
	Role     string
	Env      string
	UserName string
	Profile  string
}

type Environment struct {
	Roles map[string]Role
}

type Role struct {
	Groups []string `yaml:"groups"`
}

const (
	passwordLen    = 10
	defaultEnv     = "preprod"
	defaultRole    = "developer"
	groupsFileName = "configs/groups.yaml"
)

// UnmarshalYAML allows us to keep our custom type of YAML without having to name the groups in the Environment struct
// http://stackoverflow.com/questions/32147325/how-to-parse-yaml-with-dyanmic-key-in-golang
func (e *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var roles map[string]Role
	if err := unmarshal(&roles); err != nil {
		fmt.Printf("unmarshaling in this environment struct")
		if _, ok := err.(*yaml.TypeError); !ok {
			return err
		}
	}
	e.Roles = roles
	return nil
}

// IamInstance creates an instance of a Iam Client with the default ~/.aws/credentials
func IamInstance(s *Singapura) (*IAM, error) {
	sess, err := awsSession(s.Profile)
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
func (i *IAM) CreateUser(s *Singapura) (*iam.CreateUserOutput, error) {
	res, err := i.Iam.CreateUser(&iam.CreateUserInput{
		UserName: aws.String(s.UserName),
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateUserPassword generates a random password with the securerandom lib
// for the specified username
func (i *IAM) CreateUserPassword(s *Singapura) (string, error) {
	var err error
	var password string
	password, err = securerandom.Base64(passwordLen, true)
	if err != nil {
		return "", err
	}

	_, err = i.Iam.CreateLoginProfile(&iam.CreateLoginProfileInput{
		UserName:              aws.String(s.UserName),
		Password:              aws.String(password),
		PasswordResetRequired: aws.Bool(true),
	})

	if err != nil {
		return "", err
	}
	return password, nil
}

// AddUserGroups adds groups to a specified user based off the flags that are passed in
func (i *IAM) AddUserGroups(s *Singapura) ([]string, error) {
	var outputs []string
	groups, err := GroupsByRoleAndEnv(s)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		_, err = i.Iam.AddUserToGroup(&iam.AddUserToGroupInput{
			GroupName: aws.String(group),
			UserName:  aws.String(s.UserName),
		})
		if err != nil {
			fmt.Printf("Error adding user to group: %v\n", group)
			continue
		}
		outputs = append(outputs, group)
	}
	return outputs, nil
}

// CreateAccessKey generates access keys for a user
func (i *IAM) CreateAccessKey(s *Singapura) (*iam.CreateAccessKeyOutput, error) {
	k, err := i.Iam.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(s.UserName),
	})
	if err != nil {
		return nil, err
	}
	return k, nil
}

// GroupsByRoleAndEnv returns a list of groups based off of environment and passed in role
func GroupsByRoleAndEnv(s *Singapura) ([]string, error) {
	setDefaultEnvRole(s)
	env, err := allGroups()
	if err != nil {
		return nil, err
	}
	roles, ok := env[s.Env]
	if !ok {
		return nil, fmt.Errorf("Unable to find any roles related to env: %v\n", s.Env)
	}
	var role Role
	role, ok = roles.Roles[s.Role]
	if !ok {
		return nil, fmt.Errorf("Unable to find role: %v\n", s.Role)
	}
	return role.Groups, nil
}

// AllRoles returns a map keyed by Environment, and contains roles/groups that are keyed
func allGroups() (map[string]Environment, error) {
	filename, _ := filepath.Abs(groupsFileName)
	file, err := ioutil.ReadFile(filename)

	var envRoles map[string]Environment
	err = yaml.Unmarshal(file, &envRoles)

	if err != nil {
		return nil, err
	}

	return envRoles, nil
}

// setDefaultEnvRole sets default environment and roles if not set
func setDefaultEnvRole(g *Singapura) {
	if g.Env == "" {
		g.Env = defaultEnv
	}
	if g.Role == "" {
		g.Role = defaultRole
	}
}
