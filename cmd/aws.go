package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/grindrllc/singapura/singapura"
	"github.com/spf13/cobra"
)

var createUserCmd = &cobra.Command{
	Use:   "createUser",
	Short: "Create a new AWS User",
	Long:  `Create a new AWS User - username, path`,
	Run: func(cmd *cobra.Command, args []string) {
		//Need a --profile flag in order to get credentials other than default
		if len(args) < 1 {
			fmt.Println("Imma need a username doofus")
			return
		}
		i := singapura.IamInstance()
		u := &iam.CreateUserInput{
			UserName: aws.String(args[0]),
		}
		res, err := i.CreateUserRequest(u)
		if err != nil {
			fmt.Println("Error creating user: %v\n", err)
		}
		fmt.Printf("Response from creating user: %v\n", res)
	},
}
