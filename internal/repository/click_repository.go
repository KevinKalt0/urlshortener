package repository

import (
<<<<<<< Updated upstream
	"fmt"

	"github.com/KevinKalt0/urlshortener/internal/models"
	"gorm.io/gorm"
=======
    "gorm.io/gorm"
    "url-shortener/internal/models"
>>>>>>> Stashed changes
)

type ClickRepository interface {
    Create(click *models.Click) error
    CountByLinkID(linkID uint) (int64, error)
}

type clickRepository struct {
    db *gorm.DB
}

func NewClickRepository(db *gorm.DB) ClickRepository {
    return &clickRepository{db: db}
}

func (r *clickRepository) Create(click *models.Click) error {
    return r.db.Create(click).Error
}

func (r *clickRepository) CountByLinkID(linkID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error
    return count, err
}
