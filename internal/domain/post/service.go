package post

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePost(p Post) error {

	err := s.repo.Insert(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetPost(id string) (*Post, error) {
	return &Post{}, nil
}

func (s *Service) GetAll() ([]Post, error) {
	// s.repo.DeleteById(?
	return []Post{}, nil
}
