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
	Role     string
	Env      string
	UserName string
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
func (i *IAM) AddUserGroups(g *GroupConfig) ([]*iam.AddUserToGroupOutput, error) {
	var outputs []*iam.AddUserToGroupOutput
	groups, err := GroupsByRoleAndEnv(g)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		var o *iam.AddUserToGroupOutput
		o, err = i.Iam.AddUserToGroup(&iam.AddUserToGroupInput{
			GroupName: aws.String(group),
			UserName:  aws.String(g.UserName),
		})
		if err != nil {
			fmt.Printf("Error adding user to group: %v\n", group)
			continue
		}
		outputs = append(outputs, o)
	}
	return outputs, nil
}

// GroupsByRoleAndEnv returns a list of groups based off of environment and passed in role
func GroupsByRoleAndEnv(g *GroupConfig) ([]string, error) {
	setDefaultEnvRole(g)
	env, err := allGroups()
	if err != nil {
		return nil, err
	}
	roles, ok := env[g.Env]
	if !ok {
		return nil, fmt.Errorf("Unable to find any roles related to env: %v\n", g.Env)
	}
	var role Role
	role, ok = roles.Roles[g.Role]
	if !ok {
		return nil, fmt.Errorf("Unable to find role: %v\n", g.Role)
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
func setDefaultEnvRole(g *GroupConfig) {
	if g.Env == "" {
		g.Env = defaultEnv
	}
	if g.Role == "" {
		g.Role = defaultRole
	}
}
