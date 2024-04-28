package errors

import "errors"

// Organization errors
var (
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrOrganizationExists   = errors.New("organization already exists")
	ErrOrganizationCreate   = errors.New("failed to create organization")
	ErrOrganizationUpdate   = errors.New("failed to update organization")
	ErrOrganizationDelete   = errors.New("failed to delete organization")
	ErrOrganizationGet      = errors.New("failed to get organization")
	ErrOrganizationGetAll   = errors.New("failed to get all organizations")
)

// Oranization Member errors
var (
	ErrNotMemberOfOrganization    = errors.New("user is not a member of the organization")
	ErrOrganizationMemberNotFound = errors.New("organization member not found")
	ErrOrganizationMemberExists   = errors.New("organization member already exists")
	ErrOrganizationMemberCreate   = errors.New("failed to create organization member")
	ErrOrganizationMemberUpdate   = errors.New("failed to update organization member")
	ErrOrganizationMemberDelete   = errors.New("failed to delete organization member")
	ErrOrganizationMemberGet      = errors.New("failed to get organization member")
	ErrOrganizationMemberGetAll   = errors.New("failed to get all organization members")
)
