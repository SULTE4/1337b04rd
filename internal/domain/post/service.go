package post

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePost(p Post) error {

	return s.repo.Insert(p)
}

func (s *Service) GetPost(id string) (*Post, error) {
	return &Post{}, nil
}
