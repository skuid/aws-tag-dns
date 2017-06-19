package manager

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/skuid/spec"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

var NoRecord error = errors.New("No record")

func init() {
	var err error
	logger, err = spec.NewStandardLogger()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}

// Manager performs Route53 updates
type Manager struct {
	session *session.Session
	config  Config
	Domain  string
}

// New returns a new manager with the given Config attached and initializes an AWS session
func New(config Config) (*Manager, error) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(config.Region)}))
	m := Manager{
		session: sess,
		config:  config,
	}
	svc := route53.New(m.session)

	resp, err := svc.GetHostedZone(&route53.GetHostedZoneInput{
		Id: aws.String(config.HostedZoneID),
	})
	if err != nil {
		return nil, err
	}
	m.Domain = *resp.HostedZone.Name

	return &m, nil
}

// getTargets returns a list of targets that the new DNS records will point to.
func (m Manager) getTargets() (targets []string, err error) {
	svc := ec2.New(m.session)

	filters := []*ec2.Filter{}
	for k, v := range m.config.ASGTags.ToMap() {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", k)),
			Values: []*string{aws.String(v)},
		})
	}
	filters = append(filters, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []*string{aws.String("running")},
	})

	response, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: filters,
	})
	if err != nil {
		return targets, err
	}

	for _, res := range response.Reservations {
		for _, instance := range res.Instances {
			switch m.config.UsePrivateIP {
			case false:
				if *instance.PublicIpAddress == "" {
					logger.Warn(
						"Public IP's requested, but instance doesn't have a public IP! Skipping instance",
						zap.String("instance-id", *instance.InstanceId),
					)
					continue
				}
				if m.config.RecordType == "A" {
					targets = append(targets, *instance.PublicIpAddress)
				}
				if m.config.RecordType == "CNAME" {
					targets = append(targets, *instance.PublicDnsName)
				}
			default:
				if m.config.RecordType == "A" {
					targets = append(targets, *instance.PrivateIpAddress)
				}
				if m.config.RecordType == "CNAME" {
					targets = append(targets, *instance.PrivateDnsName)
				}
			}
		}
	}
	return targets, err
}

// buildRoutes returns a map of the DNS names for the given IPs
func (m Manager) buildRoutes(ips []string) map[string]string {
	response := make(map[string]string)

	sort.Strings(ips) // Sort ips for reproducible results

	for i, ip := range ips {
		response[fmt.Sprintf("%s%d.%s", m.config.SubdomainPrefix, i, m.Domain)] = ip
	}
	return response
}

// getRoute returns the value of a given DNS name. Returns NoRoute if
// the route doesn't exist
func (m Manager) getRoute(route string) (string, error) {
	svc := route53.New(m.session)
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(m.config.HostedZoneID), // Required
		StartRecordName: aws.String(route),
	}
	resp, err := svc.ListResourceRecordSets(params)
	if err != nil {
		return "", err
	}
	if len(resp.ResourceRecordSets) == 0 {
		return "", NoRecord
	}
	if len(resp.ResourceRecordSets[0].ResourceRecords) == 0 {
		return "", NoRecord
	}
	return *resp.ResourceRecordSets[0].ResourceRecords[0].Value, nil
}

// defineUpdates creates a map of any route changes to be made
func (m Manager) defineUpdates() (map[string]string, error) {
	targets, err := m.getTargets()
	if err != nil {
		return nil, err
	}
	desiredRoutes := m.buildRoutes(targets)

	updates := map[string]string{}

	for k, v := range desiredRoutes {
		existingRecord, err := m.getRoute(k)
		if err != nil && err != NoRecord {
			logger.Error(
				"Couldn't determine the status of route",
				zap.String("route", k),
				zap.Error(err),
			)
			continue
		}
		if v != existingRecord {
			updates[k] = v
		}
	}
	return updates, nil
}

// assignRoutes performs the given route assignments
func (m Manager) assignRoutes(assignments map[string]string) error {
	svc := route53.New(m.session)

	changes := []*route53.Change{}

	for k, v := range assignments {
		changes = append(changes, &route53.Change{
			Action: aws.String("UPSERT"),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name:            aws.String(k),
				Type:            aws.String(m.config.RecordType),
				ResourceRecords: []*route53.ResourceRecord{{Value: aws.String(v)}},
				TTL:             aws.Int64(m.config.TTL),
			},
		})
	}

	_, err := svc.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(m.config.HostedZoneID),
	})
	return err
}

// AssignRoutes creates the DNS entries in Route53
func (m Manager) AssignRoutes() error {
	updates, err := m.defineUpdates()
	if err != nil {
		return err
	}
	if len(updates) == 0 {
		logger.Info("No updates required.")
		return nil
	}
	fields := []zapcore.Field{}
	for k, v := range updates {
		fields = append(fields, zap.String(k, v))
	}
	logger.Info(
		"The following routes will be assigned",
		fields...,
	)
	if m.config.Debug {
		return nil
	}
	err = m.assignRoutes(updates)
	if err != nil {
		return err
	}
	logger.Info(
		"Routes successfully updated!",
		fields...,
	)
	return nil
}
