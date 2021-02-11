package patient

type repo interface {
	Searcher
}

type Usecase struct {
	Searcher Searcher
}

func NewUsecase(repo repo) *Usecase {
	return &Usecase{
		Searcher: repo,
	}
}
