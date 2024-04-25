package repository

type Repository[T any] interface {
	Store(*T) error
	Update(uint, *T) error
	Find(*T) (*T, error)
	Delete(uint) error
}
