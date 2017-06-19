[![Build Status](https://travis-ci.org/skuid/aws-tag-dns.svg)](https://travis-ci.org/skuid/aws-tag-dns)
[![https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/skuid/aws-tag-dns/)

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
  "etcd.skuid.com.": "172.20.67.72",
  "etcd.skuid.com.": "172.20.88.113"
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

## Developing

### Dependencies

Use [`dep`](https://github.com/golang/dep). 

```bash
dep status
dep ensure
```

## License
MIT License. See [License](/LICENSE) for full text
