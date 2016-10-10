package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/grindrllc/singapura/singapura"
	"github.com/spf13/cobra"
)

var createUserCmd = &cobra.Command{
	Use:   "CreateUser",
	Short: "Create a new AWS User",
	Long:  `Create a new AWS User - username`,
	Run: func(cmd *cobra.Command, args []string) {
		var i *singapura.IAM
		var err error
		username := args[0]
		if username == "" {
			fmt.Println("Provide a username to create account.")
			return
		}
		i, err = singapura.IamInstance(profile)

		var userRes *iam.CreateUserOutput
		userRes, err = i.CreateUser(&username)

		var loginRes *iam.CreateLoginProfileOutput
		loginRes, err = i.CreateUserPassword(&username)
		if err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			return
		}
		fmt.Printf("User res: %v\n", userRes)
		fmt.Printf("Login res: %v\n", loginRes)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "DeleteUser",
	Short: "Delete an AWS User",
	Long:  `Delete an AWS User - username`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
