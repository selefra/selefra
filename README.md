<h1 align="center">
    <a href="https://www.selefra.io" title="Selefra - Infrastructure as Code for Infrastructure Analysis.">
        <img src=".github/images/logo_colorbg.png" width="900">
    </a>
</h1>

# Introduction

Selefra is an open-source data integration and analysis tool for developers. You can use Selefra to extract, load, and analyze infrastructure data anywhere from Public Cloud, SaaS platform, development platform, and more.

Simply write code in  YAML and SQL, Selefra automatically  pull control plane data and analysis your AWS, GCP, Azure, and [other hosted data source providers](https://github.com/selefra/registry).

For example, here's sample usage of test item `ebs_volume_are_unencrypted`:
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
Read detailed documentation for how to [get started](https://selefra.io/docs/get-started/) with Selefra.

For quick start, run this demo, it should take less than a few miniutes:

### 1. Install Selefra

For non macOS users, [download packages](https://github.com/selefra/selefra/releases) to install Selefra.

On macOS, tap Selefra with Homebrew:

```bash
brew tap selefra/tap
```

Next, install Selefra:

```bash
brew install selefra/tap/selefra
```

### 2. Initialization project

```bash
selefra init selefra-demo && cd selefra-demo
```

### 3. Build code

```bash
selefra apply 
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
