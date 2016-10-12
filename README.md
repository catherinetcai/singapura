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
-p, --p - Allows you to specify the AWS profile to use from your ~/.aws/credentials file
-e, --env - Allows you to specify the environment in your config/groups.yaml
```
