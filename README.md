# Singapura
> A tiny cat breed name for a tiny AWS tool.

### What does Singapura do?
Singapura is a CLI tool that automates the process of creating an AWS user. It generates a password for them, adds them to groups that you specify, and creates access keys for them.

### Great! How do I do that?
You can build the binary yourself: 

```
go build
```

...or you can run it if you have Go installed:

```
go run main.go //interactions go here
```

Singapura uses [spf13's Cobra](https://github.com/spf13/cobra) for CLI interactions, so you can see all the helpful commands by running this:

```
go run main.go help
```

#### Creating a User

```
./singapura CreateUser ausernamehere
```

This will create a username, generate a password for them, generate an access key for them, and add them to the default groups set in config/groups.yaml.

There's flags to allow you to pass in certain options:

```
-p, --p - string - Specify the AWS profile to use from your ~/.aws/credentials file. It will default to the default aws profile if not specified
-e, --env - string - Specify the environment in your config/groups.yaml. Defaults to preprod if not specified
-r, --role - string - Specify the role to apply to groups. Defaults to developer if not specified
-k, --key - bool - Specify whether or not access keys should be generated for the user. Defaults to true if not specified
```

#### Groups File

Singapura reads from the config/groups.yaml file to determine what groups to add. Here's an example of how to configure it:

```
env_one:
  role_one:
    groups:
      - "Developers"
  role_two:
    groups:
      - "devops"
      - "DevOps-NOC-Dynamo-Ec2"
      - "DevOps-NOC-ElasticCache-Route53"
      - "DevOps-NOC-NOC-SUPPORT-DASHBOARD"
      - "DevOPS-ReadOnly-Cloudwatch-Opsworks"
      - "MFA_ENFORCED"
      - "DevOps-NOC-EBS-VPC"
      - "DevOps-NOC-CloudFormation-Cloudfront"
      - "Grindr-ADMIN-ALL-ACCESS-PLAN-fixed-ma-03052016"
      - "DevOps-SNS-SQS"
env_two:
```
