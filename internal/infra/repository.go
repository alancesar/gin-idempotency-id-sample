package infra

import (
	"context"
	"gorm.io/gorm"
	"idempotency/user"
)

type (
	GormUserRepository struct {
		db *gorm.DB
	}

	userEntity struct {
		gorm.Model
		Name  string
		Email string
	}
)

func (e userEntity) ToDomain() user.User {
	return user.User{
		ID:    e.ID,
		Name:  e.Name,
		Email: e.Email,
	}
}

func (e userEntity) TableName() string {
	return "users"
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	_ = db.AutoMigrate(&userEntity{})
	return &GormUserRepository{
		db: db,
	}
}

func (r GormUserRepository) Save(ctx context.Context, u user.User) (user.User, error) {
	entity := newUserFromDomain(u)
	if query := r.db.WithContext(ctx).Save(&entity); query.Error != nil {
		return user.User{}, query.Error
	}

	return entity.ToDomain(), nil
}

func (r GormUserRepository) Get(ctx context.Context, id uint) (user.User, error) {
	entity := userEntity{
		Model: gorm.Model{
			ID: id,
		},
	}

	if query := r.db.WithContext(ctx).Take(&entity); query.Error != nil {
		return user.User{}, query.Error
	}

	return entity.ToDomain(), nil
}

func newUserFromDomain(u user.User) userEntity {
	return userEntity{
		Model: gorm.Model{
			ID: u.ID,
		},
		Name:  u.Name,
		Email: u.Email,
	}
}
