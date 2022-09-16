package init

const ruleComment = `
rules:
  - name: HostHighCpuLoad
    input:
      account_id:
        type: string
        description: "cpu rate"
        default: "587534146112"
    query: SELECT * FROM aws_s3_account_config WHERE "account_id" = '{{.account_id}}';
    interval: 0s
    labels:
      severity: warning
      team: ops
      author: leon
    metadata:
      id: SF001
      summary: Host high CPU load (instance {{ .labels.instance }})
      description: "Test desc "
    output: "{{.account_id}} is warning, block_public_acls: {{.block_public_acls}}"

`

const moduleComment = `
modules:
  - name: CpuMonitor
    uses: ./rules/default.yaml
    input:
      name: uzju
      labels:
        instance: warning
`
