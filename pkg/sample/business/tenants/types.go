package tenants

import "time"

type TenantID string

type TenantIDGenerator func() TenantID

type Tenant struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	ID        TenantID   `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
}
