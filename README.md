<h1 align="center">
    <a href="https://www.selefra.io" title="Selefra - Infrastructure as Code for Infrastructure Analysis.">
        <img src=".github/images/logo_colorbg.png" width="900">
    </a>
    <p align="center">
    <a href="https://github.com/selefra/selefra/stargazers"><img alt="GitHub stars" src="https://img.shields.io/github/stars/selefra/selefra"/></a>
    <a href="https://github.com/selefra/selefra/releases"><img alt="GitHub releases" src="https://img.shields.io/github/release/teamssix/cf"/></a>
    <a href="https://github.com/selefra/selefra/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/badge/license-MPL%202.0-blue.svg"/></a>
    <a href="https://github.com/selefra/selefra/releases"><img alt="Downloads" src="https://img.shields.io/github/downloads/selefra/selefra/total?color=blue"/></a>
    <a href="https://goreportcard.com/report/github.com/selefra/selefra"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/selefra/selefra"/></a>
    <a href="https://twitter.com/SelefraCorp"><img alt="Twitter" src="https://img.shields.io/twitter/follow/SelefraCorp?style=social" /></a>
    <a href="https://github.com/selefra"><img alt="Github" src="https://img.shields.io/github/followers/selefra?style=social" /></a>
    <a href="https://selefra.io/community/join"><img alt="Slack" src="https://img.shields.io/badge/Join%20Slack-%40Selefra-red" /></a><br></br>    
    </p>
</h1>

### Why Selefra?

Selefra is an open-source cloud governance tool to analyze multi-cloud assets for security, compliance, and policy enforcement. 
* **Detect-to-Remediate**: discover and remediate unnoticed risky problems in one stop.
* **Provider Agnostic**: reduce switching cost between isolated control planes.
* **Ease-of-Use**: simplified usage to write and maintain for quick fixes and long-term usage.

With rules written in YAML and SQL, Selefra automatically pulls data from providers including AWS, GCP, Azure, and [more](https://github.com/selefra/selefra).

For example, a rule to check if *EBS volume are unencrypted*:

```yaml
selefra:
    name: example_project
    cli_version: v0.0.1
    providers:
        - name: aws
          source: selefra/aws
          version: v0.0.3

providers:
    - name: aws
      cache: 1d1h1m1s
      resources:
        - aws_*
      accounts:
         regions:
           - us-east-1
rules:
  - name: example_rule_name
    query: SELECT * FROM aws_ec2_ebs_volumes WHERE encrypted = FALSE
    labels: 
      tag: demo_rule
      author: Selefra
    metadata: 
      severity: Low
      provider: AWS
      resource_type: EC2 
      resource_account_id: '{{.account_id}}'
      resource_id: '{{.id}}'
      resource_region: '{{.availability_zone}}'
      remediation: remediation/ebs_volume_are_unencrypted.md 
      title: EBS volume are unencrypted 
      description: Ensure that EBS volumes are encrypted.
    output: 'EBS volume is unencrypted, EBS id: {{.id}}, availability zone: {{.availability_zone}}'
```

## Getting Started
 
For detailed instructions on prerequisites, check [Get Started](https://selefra.io/docs/get-started/) for more info.

Otherwise, run a demo through the following steps, it should take less than a few miniutes:

### 1. Install Selefra

[download packages](https://github.com/selefra/selefra/releases) to install Selefra.

If you are MacOS user, tap Selefra with Homebrew.

```bash
$ brew tap selefra/tap
```

Now, install Selefra

```bash
$ brew install selefra/tap/selefra
```

### 2. Create a project

```bash
$ selefra init selefra-demo && cd selefra-demo
```

### 3. Build code for the project

```bash
$ selefra apply 
```


## Documentation

See [Docs](https://selefra.io/docs) for best practices and detailed instructions. In docs, you will find info on installation, CLI usage, project workflow and more guides on how to accomplish cloud inspection tasks.

## Community

Selefra is a community-driven project, we welcome you to open a [GitHub Issue](https://github.com/selefra/selefra/issues/new/choose) to report a bug, suggest an improvement, or request new feature.

-  Join [Selefra Community](https://selefra.io/community/join) on Slack. We host `Community Hour` for tutorials and Q&As on regular basis.
-  Follow us on [Twitter](https://twitter.com/SelefraCorp) and share your thoughts！
-  Email us at support@selefra.io

## CONTRIBUTING

For developers interested in building Selefra codebase, read through [Contributing.md](https://github.com/selefra/selefra/blob/main/CONTRIBUTING.md) and [Selefra Roadmap](https://github.com/orgs/selefra/projects/1). 
Let us know what you would like to work on!

## License

[Mozilla Public License v2.0](https://github.com/selefra/selefra/blob/main/LICENSE)


