package click

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// ClickDAO represents a single database entry.
type ClickDAO struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	URL       string
}

// ClickDAOCollection represents a collection of Click database model.
type ClickDAOCollection []ClickDAO

// ToDomain maps DAO models into domain models.
func (cc ClickDAOCollection) ToDomain() ClickCollection {
	r := make(ClickCollection, 0, len(cc))

	for _, c := range cc {
		r = append(r, c.ToDomain())
	}

	return r
}

// TableName overrides the table name used by ClickDAO to 'clicks'
func (ClickDAO) TableName() string {
	return "clicks"
}

// NewClickDAO maps Click entity model into database model.
func NewClickDAO(c Click) ClickDAO {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	return ClickDAO{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		URL:       c.URL,
	}
}

// ToModel maps database model into domain model.
func (c *ClickDAO) ToDomain() Click {
	return Click{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		URL:       c.URL,
	}
}

// SQLiteRepository is a SQLite implementation of Click repository.
type SQLiteRepository struct {
	db *gorm.DB
}

// Create persists Click entity.
func (r *SQLiteRepository) Create(ctx context.Context, click Click) (Click, error) {
	dao := NewClickDAO(click)

	result := r.db.WithContext(ctx).Create(&dao)
	if result.Error != nil {
		return Click{}, result.Error
	}

	return dao.ToDomain(), nil
}

// Filter applies provided filters and returns resulting subset.
func (r *SQLiteRepository) Filter(ctx context.Context, filter Filter) (ClickCollection, error) {
	var clicks ClickDAOCollection

	tx := r.db.WithContext(ctx)

	if filter.URL != "" {
		tx = tx.Where("url = ?", filter.URL)
	}
	if !filter.After.IsZero() {
		tx = tx.Where("created_at > ?", filter.After)
	}
	if !filter.Before.IsZero() {
		tx = tx.Where("created_at < ?", filter.Before)
	}

	if tx.Find(&clicks).Error != nil {
		return ClickCollection{}, r.db.Error
	}

	return clicks.ToDomain(), nil
}

// NewSQLiteRepository is a SQLiteRepository constructor.
func NewSQLiteRepository(db *gorm.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}
