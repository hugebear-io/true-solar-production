package infra

import (
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/hugebear-io/true-solar-production/config"
)

func NewSnmp() (*gosnmp.GoSNMP, error) {
	conf := config.GetConfig().Snmp

	client := &gosnmp.GoSNMP{
		Target:             conf.TargetHost,
		Port:               conf.TargetPort,
		Transport:          "udp",
		Community:          "public",
		Version:            gosnmp.Version1,
		Timeout:            time.Duration(300) * time.Second,
		Retries:            20,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}

	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}
