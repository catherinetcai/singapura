package cmd

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/grindrllc/singapura/singapura"
	"github.com/spf13/cobra"
)

var createUserCmd = &cobra.Command{
	Use:   "CreateUser",
	Short: "Create a new AWS User",
	Long: `Create a new AWS User.
	Usage: CreateUser --profile <profile> --env <environment> --role <role>.
	Defaults to preprod and developer for role if unspecified`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var i *singapura.IAM
		var err error
		if len(args) <= 0 {
			return errors.New("Provide a username to create account.")
		}
		username := args[0]
		if username == "" {
			return errors.New("Provide a username to create account.")
		}

		s := createConfig(username)
		i, err = singapura.IamInstance(s)
		if err != nil {
			return errors.New("Error connecting to AWS with the profile specified.")
		}

		fmt.Println("Creating user...")
		var userRes *iam.CreateUserOutput
		userRes, err = i.CreateUser(s)
		if err != nil {
			return err
		}

		fmt.Println("Generating password...")
		var loginRes string
		loginRes, err = i.CreateUserPassword(s)
		if err != nil {
			return err
		}

		var accessRes *iam.CreateAccessKeyOutput
		if accesskey {
			fmt.Println("Creating access keys...")
			accessRes, err = i.CreateAccessKey(s)
			if err != nil {
				return err
			}
		}

		fmt.Println("Adding user to groups...")
		var groupRes []string
		groupRes, err = i.AddUserGroups(s)
		if err != nil {
			return err
		}

		fmt.Printf("User Info: %v\n", userRes)
		fmt.Printf("Password: %v\n", loginRes)
		fmt.Printf("Access key: %v\n", accessRes)
		fmt.Printf("Groups: %v\n", groupRes)
		return nil
	},
}

// createConfig creates and populates the struct for Singapura
func createConfig(username string) *singapura.Singapura {
	return &singapura.Singapura{
		Profile:  profile,
		Env:      env,
		Role:     role,
		UserName: username,
	}
}

var deleteUserCmd = &cobra.Command{
	Use:   "DeleteUser",
	Short: "Delete an AWS User",
	Long:  `Delete an AWS User - username`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
