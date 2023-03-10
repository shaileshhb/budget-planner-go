package general

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains master fields for all entities.
type Base struct {
	ID         uuid.UUID  `gorm:"type:varchar(36);primary_key" json:"id" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `gorm:"index:idx_deleted_at" json:"-"`
	IgnoreHook bool       `gorm:"-" json:"-"`
	// CreatedBy  uuid.UUID  `gorm:"type:varchar(36)" json:"-"`
	// UpdatedBy  uuid.UUID  `gorm:"type:varchar(36)" json:"-"`
	// DeletedBy  uuid.UUID  `gorm:"type:varchar(36)" json:"-"`
}

// BeforeCreate will be called before the entity is added to db.
func (b *Base) BeforeCreate(scope *gorm.DB) error {
	if b.IgnoreHook {
		return nil
	}

	b.ID = uuid.New()
	return nil
}

// BaseDTO contains master fields for DTO specifically.
// Should only be used for reading operations.
type BaseDTO struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
