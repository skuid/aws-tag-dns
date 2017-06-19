package main

import (
	"fmt"
	"log"
	"os"
	"time"

	flag "github.com/ogier/pflag"
	"github.com/skuid/aws-tag-dns/manager"
	"github.com/skuid/spec"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = spec.NewStandardLogger()
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}

	err = os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	var config manager.Config
	flag.StringVarP(&config.HostedZoneID, "zone", "z", "", "The Route53 hosted zone ID to use")
	flag.StringVar(&config.SubdomainPrefix, "prefix", "", "The DNS subdomain prefix to use")
	flag.VarP(&config.ASGTags, "tags", "t", `The tags to match on. Each tag should have the format "k=v"`)
	flag.StringVar(&config.Region, "region", os.Getenv("AWS_DEFAULT_REGION"), "AWS region to use, looks for `AWS_DEFAULT_REGION` env var")
	flag.StringVar(&config.RecordType, "record", "A", "The record type to use. Must be \"A\" or \"CNAME\"")
	flag.BoolVar(&config.UsePrivateIP, "private", true, "Use the instance's private IP or DNS")
	flag.Int64Var(&config.TTL, "ttl", 60, "The TTL to set")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Don't actually update records")

	var period = flag.DurationP("period", "p", time.Duration(time.Minute), "The interval for the update to run")

	flag.Parse()

	timerMinimum := time.Duration(time.Second * 10)
	if *period < timerMinimum {
		logger.Warn(
			"Timer period set too low, increasing to 10 seconds",
			zap.Duration("specified", *period),
			zap.Duration("using", timerMinimum),
		)
		*period = timerMinimum
	}

	err := config.Validate()
	if err != nil {
		logger.Fatal(err.Error())
	}

	m, err := manager.New(config)
	if err != nil {
		logger.Fatal(err.Error())
	}

	assignFunc := func() {
		err = m.AssignRoutes()
		if err != nil {
			logger.Error(err.Error())
		}
	}

	go assignFunc()
	c := time.Tick(*period)
	for range c {
		go assignFunc()
	}
}
