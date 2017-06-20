[![Build Status](https://travis-ci.org/skuid/aws-tag-dns.svg)](https://travis-ci.org/skuid/aws-tag-dns)
[![https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/skuid/aws-tag-dns/)
[![aws-tag-dns](https://quay.io/repository/skuid/aws-tag-dns/status "aws-tag-dns")](https://quay.io/repository/skuid/aws-tag-dns)

# aws-tag-dns

Generate DNS records for a given set of AWS instance tags.

Example:

```bash
aws-tag-dns --prefix etcd --private --record A --tags Name=MyEtcdAsg --zone Z2SNGMHS3A6Z7I
{
  "level": "info",
  "timestamp": "2017-06-19T13:37:40.827-0400",
  "caller": "manager/manager.go:216",
  "message": "The following routes will be assigned",
  "etcd0.skuid.com.": "172.20.67.72",
  "etcd1.skuid.com.": "172.20.88.113"
  "etcd2.skuid.com.": "172.20.112.3"
}
```

## Usage

```
Usage of aws-tag-dns:
  --dry-run
    	Don't actually update records
  -p, --period duration
    	The interval for the update to run (default 1m0s)
  --prefix string
    	The DNS subdomain prefix to use
  --private
    	Use the instance's private IP or DNS (default true)
  --record string
    	The record type to use. Must be "A" or "CNAME" (default "A")
  --region AWS_DEFAULT_REGION
    	AWS region to use, looks for AWS_DEFAULT_REGION env var
  -t, --tags value
    	The tags to match on. Each tag should have the format "k=v"
  --ttl int
    	The TTL to set (default 60)
  -z, --zone string
    	The Route53 hosted zone ID to use

```

## AWS Permissions Required

Below is an example IAM policy that outlines the permissions needed for this
application. You may optionally further restrict the `ec2:DescribeInstances`
call by specifying the region in the `Resource` section.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstances"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "route53:ChangeResourceRecordSets",
                "route53:GetHostedZone"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/<HostedZoneID>"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "route53:ListResourceRecordSets"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```

## Developing

### Dependencies

Use [`dep`](https://github.com/golang/dep).

```bash
dep status
dep ensure
```

## License
MIT License. See [License](/LICENSE) for full text
