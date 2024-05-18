package view

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// ViewDAO represents a single database entry.
type ViewDAO struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	URL       string
}

// TableName overrides the table name used by ViewDAO to 'views'
func (ViewDAO) TableName() string {
	return "views"
}

// ToModel maps database model into domain model.
func (c *ViewDAO) ToDomain() View {
	return View{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		URL:       c.URL,
	}
}

// ViewDAOCollection represents a collection of View database model.
type ViewDAOCollection []ViewDAO

func (cc ViewDAOCollection) ToDomain() ViewCollection {
	r := make(ViewCollection, 0, len(cc))

	for _, c := range cc {
		r = append(r, c.ToDomain())
	}

	return r
}

// NewViewDAO maps Click entity model into database model.
func NewViewDAO(c View) ViewDAO {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	return ViewDAO{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		URL:       c.URL,
	}
}

// SQLiteRepository is a SQLite implementation of View repository.
type SQLiteRepository struct {
	db *gorm.DB
}

// Create persists Click entity.
func (r *SQLiteRepository) Create(ctx context.Context, view View) (View, error) {
	dao := NewViewDAO(view)

	result := r.db.WithContext(ctx).Create(&dao)
	if result.Error != nil {
		return View{}, result.Error
	}

	return dao.ToDomain(), nil
}

func (r *SQLiteRepository) Filter(ctx context.Context, filter Filter) (ViewCollection, error) {
	var views ViewDAOCollection

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

	if tx.Find(&views).Error != nil {
		return ViewCollection{}, r.db.Error
	}

	return views.ToDomain(), nil
}

// NewSQLiteRepository is a SQLiteRepository constructor.
func NewSQLiteRepository(db *gorm.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}
