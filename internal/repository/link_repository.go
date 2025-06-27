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

type LinkRepository interface {
    Create(link *models.Link) error
    FindByCode(code string) (*models.Link, error)
}

type linkRepository struct {
    db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) LinkRepository {
    return &linkRepository{db: db}
}

func (r *linkRepository) Create(link *models.Link) error {
    return r.db.Create(link).Error
}

func (r *linkRepository) FindByCode(code string) (*models.Link, error) {
    var link models.Link
    err := r.db.Where("short_code = ?", code).First(&link).Error
    return &link, err
}
