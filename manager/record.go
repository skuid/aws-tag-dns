package manager

import (
	"fmt"

	"github.com/skuid/spec"
)

type RecordType string

const A RecordType = "A"
const CNAME RecordType = "CNAME"

// Config holds configuration values for a Manager
type Config struct {
	HostedZoneID    string
	SubdomainPrefix string
	Region          string
	RecordType      string
	UsePrivateIP    bool
	Debug           bool
	TTL             int64
	ASGTags         spec.SelectorSet
}

// Validate performs validation on a Config
func (c Config) Validate() error {
	if !(c.RecordType == "CNAME" || c.RecordType == "A") {
		return fmt.Errorf("Invalid record type '%s', must be CNAME or A", c.RecordType)
	}
	return nil
}
