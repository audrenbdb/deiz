package patient

type repo interface {
	Searcher
	Creater
}

type Usecase struct {
	Searcher Searcher
	Creater  Creater
}

func NewUsecase(repo repo) *Usecase {
	return &Usecase{
		Searcher: repo,
		Creater:  repo,
	}
}
