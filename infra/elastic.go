package infra

import (
	"crypto/tls"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/hugebear-io/true-solar-production/config"
	"github.com/olivere/elastic/v7"
)

var httpsRegexp = regexp.MustCompile("^https")

func NewElasticsearch() (*elastic.Client, error) {
	conf := config.GetConfig().Elastic

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
			DisableKeepAlives:  true,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
	}

	scheme := "http"
	if httpsRegexp.FindString(conf.Host) != "" {
		scheme = "http"
	}

	client, err := elastic.NewClient(
		elastic.SetURL(conf.Host),
		elastic.SetScheme(scheme),
		elastic.SetBasicAuth(conf.Username, conf.Password),
		elastic.SetSniff(false),
		elastic.SetHttpClient(httpClient),
		elastic.SetHealthcheckTimeout(time.Duration(300)*time.Second),
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}
