package repository

import (
	"github.com/bbconfhq/mycodes/models"
	"gorm.io/gorm"
	"time"
)

type CodeRepo interface {
	Create(data *models.Code) error
	Get(id string) (*models.Code, error)
	GetRecent() (*[]models.Code, error)
	DeleteExpired() (int64, error)
}

type codeRepoImpl struct {
	db *gorm.DB
}

var (
	repo *codeRepoImpl
)

func InitializeCodeRepo(db *gorm.DB) CodeRepo {
	repo = &codeRepoImpl{db: db}
	return repo
}

func (c *codeRepoImpl) Create(data *models.Code) error {
	code := models.Code{
		ID:        data.ID,
		Ip:        data.Ip,
		Name:      data.Name,
		Title:     data.Title,
		Content:   data.Content,
		Language:  data.Language,
		ExpiredAt: data.ExpiredAt,
	}
	result := c.db.Create(&code)
	if result.Error == nil {
		data.CreatedAt = code.CreatedAt
	}
	return result.Error
}

// Get returns an unexpired code. It returns gorm.ErrRecordNotFound when the
// code does not exist or has expired.
func (c *codeRepoImpl) Get(id string) (*models.Code, error) {
	var code models.Code
	result := c.db.
		Where("id = ? AND expired_at >= ?", id, time.Now()).
		First(&code)
	if result.Error != nil {
		return nil, result.Error
	}
	return &code, nil
}

func (c *codeRepoImpl) GetRecent() (*[]models.Code, error) {
	var codes *[]models.Code
	// Everything but content, which the list does not show
	result := c.db.
		Select("id", "ip", "name", "title", "language", "created_at", "expired_at").
		Limit(5).
		Where("expired_at >= ?", time.Now()).
		Order("created_at desc").
		Find(&codes)
	if result.Error != nil {
		return nil, result.Error
	}
	return codes, nil
}

func (c *codeRepoImpl) DeleteExpired() (int64, error) {
	result := c.db.Where("expired_at < ?", time.Now()).Delete(&models.Code{})
	return result.RowsAffected, result.Error
}
