package xgorm

import "gorm.io/gorm"

type Repository[T, Key any] struct {
	db         *gorm.DB
	expands    map[string]string
	autoExpand []string
}

func (r *Repository[T, Key]) Query() *gorm.DB {
	return r.db
}

func (r *Repository[T, Key]) Save(t *T) error {
	return r.db.Save(t).Error
}

func (r *Repository[T, Key]) Create(t *T) error {
	return r.db.Create(t).Error
}

func (r *Repository[T, Key]) Count() (int, error) {
	var count int64
	var model T
	err := r.db.Model(&model).Count(&count).Error
	return int(count), err
}

func (r *Repository[T, Key]) One(id Key, expand ...string) (T, error) {
	var model T
	q := r.preload(expand...)
	err := q.First(&model, "id = ?", id).Error
	return model, err
}

func (r *Repository[T, Key]) All(expand ...string) ([]T, error) {
	var models []T
	q := r.preload()
	err := q.Find(&models).Error
	return models, err
}

func (r *Repository[T, Key]) Delete(t T) error {
	return r.db.Delete(t).Error
}

func (r *Repository[T, Key]) DeleteById(id Key) error {
	var model T
	return r.db.Delete(model, "id = ?", id).Error
}

func (r *Repository[T, Key]) preload(expand ...string) *gorm.DB {

	query := r.db

	for _, e := range r.autoExpand {
		query = query.Preload(e)
	}

	for _, e := range expand {
		v, ok := r.expands[e]
		if ok {
			query = query.Preload(v)
		}
	}

	return query
}

type User struct {
	ID string
}

type Users struct {
	*Repository[User, string]
}

type GormRepository[T struct{}] interface {
	Query() *gorm.DB

	Count() (int, error)

	Find(id string) (T, error)

	All() ([]T, error)

	Delete(t T) error

	DeleteById(id string) error

	Save(t T) error

	Create(t T) error
}
