package repo

import "github.com/hugebear-io/true-solar-production/model"

type SolarmanCredentialRepo interface {
	GetCredentials() ([]model.SolarmanCredential, error)
}
