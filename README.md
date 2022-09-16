<h1 align="center">
    <a href="https://www.selefra.io" title="Selefra - Infrastructure as Code for Infrastructure Analysis.">
        <img src=".github/images/logo.png" width="350">
    </a>
</h1>

## Introduction

Selefra's Infrastructure as Code is the easiest and agile way to analysis IAM,EC2,VPC, and infrastructure ,from any cloud control plane.

Simply write code in  YAML and SQL, Selefra automatically  pull control plane data and analysis your AWS,GCP,Azure,and [hosted data source provider](https://github.com/selefra/selefra-provider-sdk/).

For example, analysis  s3 bucket misconfigured from infrastructure:

```yaml
selefra:
  name: selefra-demo
  cli_version: ">=0.0.1"
  providers:
    - name: aws
      source: 'selefra/aws'
      version: "latest"
  connection:
    type: postgres
    username: postgres
    password: pass
    host: localhost
    port: 5432
    database: postgres
    sslmode: disable

providers:
  - name: aws
    regions:
      - us-east-1
      
rules:
  - name: Disabled_MFA
    query: select * from aws_iam_users where user_name = '<root_account>' and mfa_active = 'f'
    labels:
      severity: Critical
    metadata:
      title: "MFA is disabled for root user"
      description: "MFA is disabled for root user"
    output: "AWS user has disabled MFA, username: {{.user_name}}"
```

## Getting Started

    

Learn AWS,GCP,Azure,and more cloud/Infrastructure's usecase with [Getting Started](https://selefra.io/docs/gettingstared) .

Otherwise,run a demo process through the following setps,in miniutes:

### 1. Install Selefra

To install Selefra on MacOS or [download packages](https://github.com/selefra/selefra/releases)  to install Selefra on other platform.

Install the Selefra tap from our Homebrew packages.

```bash
brew tap selefra/tap
```

Now,install Selefra with selefra/tap/selefra

```bash
brew install selefra/tap/selefra
```

### 2. Initialization project

```bash
selefra init selefra-demo && cd selefra-demo
```

### 3. Build infrastructure Analysis code

```bash
selefra apply 
```



## Community

Selefra is a community driven project,we welcome you to file a bug,suggest an improvement ,or request a new feature.

-  Join [Selefra Community Slack](https://selefra.slack.com) to discuss Selefra and join `Selefra Community Hour`.
-  Follow us on [Twitter](https://twitter.com/SelefraCorp) and shard Selefra messages on Twitter.
-  Have question and feature?Now on [Slack](https://selefra.slack.com) or open a [GitHub Issue](https://github.com/selefra/selefra/issues/new/choose)


## Contributing

To contribute, visit [Contributing.md](https://github.com/selefra/selefra/blob/main/CONTRIBUTING.md) and [Selefra roadmap](https://github.com/orgs/selefra/projects/1)


## License

[Mozilla Public License v2.0](https://github.com/selefra/selefra/blob/main/LICENSE)
