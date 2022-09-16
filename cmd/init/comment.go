package init

const ruleComment = `
rules:
  - name: Disabled_MFA
    query: select * from aws_iam_users where user_name = '<root_account>' and mfa_active = 'f'
    labels:
      severity: Critical
    metadata:
      title: "MFA is disabled for root user"
      description: "MFA is disabled for root user"
    output: "AWS user has disabled MFA, username: {{.user_name}}"
`

const moduleComment = `
modules:
  - name: AWS_Security_Demo
    uses: ./rules/iam_mfa.yaml
`
