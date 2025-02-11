package product

type Service struct {
	repository *Repository
}

func NewService(r *Repository) *Service {
	return &Service{
		repository: r,
	}
}
