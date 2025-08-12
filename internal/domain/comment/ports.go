package comment

type Repository interface {
	Insert(com Comment) error
	GetCommentsByPost(id int) ([]Comment, error)
	UpdatePostExpire(id int) error
}
