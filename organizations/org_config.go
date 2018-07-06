package organizations

import (
	"context"

	"github.com/influxdata/chronograf"
)

// ensure that OrganizationConfig implements chronograf.OrganizationConfigStore
var _ chronograf.OrganizationConfigStore = &OrganizationConfigStore{}

// OrganizationConfigStore facade on a OrganizationConfig that filters OrganizationConfigs by organization.
type OrganizationConfigStore struct {
	store        chronograf.OrganizationConfigStore
	organization string
}

// NewOrganizationConfigStore creates a new OrganizationConfigStore from an existing
// chronograf.OrganizationConfigStore and an organization string
func NewOrganizationConfigStore(s chronograf.OrganizationConfigStore, orgID string) *OrganizationConfigStore {
	return &OrganizationConfigStore{
		store:        s,
		organization: orgID,
	}
}

// FindOrCreate gets an organization's config or creates one if none exists
func (s *OrganizationConfigStore) FindOrCreate(ctx context.Context, orgID string) (*chronograf.OrganizationConfig, error) {
	var err = validOrganization(ctx)
	if err != nil {
		return nil, err
	}

	oc, err := s.store.FindOrCreate(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return oc, nil

}

// Update the OrganizationConfig in OrganizationConfigStore.
func (s *OrganizationConfigStore) Update(ctx context.Context, oc *chronograf.OrganizationConfig) error {
	err := validOrganization(ctx)
	if err != nil {
		return err
	}

	_, err = s.store.FindOrCreate(ctx, oc.OrganizationID)
	if err != nil {
		return err
	}

	return s.store.Update(ctx, oc)
}
