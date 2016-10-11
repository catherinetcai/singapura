// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/grindrllc/singapura/singapura"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var profile string
var role string
var env string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "singapura",
	Short: "Singapura is an application that creates and deletes a user in AWS.",
	Long: `CreateUser <user> creates a user, generates a password for them, and puts them into the default groups set in the groups.yaml
	deleteUser <user> deletes the specified user in AWS
	Pass the --profile <profile> flag in order to specify the profile in your ~/.aws/credentials file that you want to run Singapura with.`,
	DisableAutoGenTag: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	//Remove this, debug
	singapura.GroupsByRoleAndEnv(&singapura.GroupConfig{})

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.singapura.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.AddCommand(createUserCmd)
	RootCmd.PersistentFlags().StringVar(&profile, "profile", "", "AWS profile to set credentials from your ~/.aws/credentials file")
	createUserCmd.Flags().StringVar(&role, "role", "", "Role to apply groups to - see the config/groups.yaml file")
	createUserCmd.Flags().StringVar(&env, "env", "", "Environment to apply groups to - see the config/groups.yaml file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".singapura") // name of config file (without extension)
	viper.AddConfigPath("$HOME")      // adding home directory as first search path
	viper.AutomaticEnv()              // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
