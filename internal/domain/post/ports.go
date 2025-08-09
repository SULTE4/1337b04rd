package post

type Repository interface {
	Insert(p Post) (int, error)
	GetByID(id int) (Post, error)
	GetAll() ([]Post, error)
	DeleteById(id int) (int, error)
}
