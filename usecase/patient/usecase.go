package patient

type repo interface {
	Searcher
	Creater
	ClinicianBoundChecker
	AddressCreater
	AddressUpdater
}

type Usecase struct {
	Searcher              Searcher
	Creater               Creater
	ClinicianBoundChecker ClinicianBoundChecker
	AddressCreater        AddressCreater
	AddressUpdater        AddressUpdater
}

func NewUsecase(repo repo) *Usecase {
	return &Usecase{
		Searcher:              repo,
		Creater:               repo,
		ClinicianBoundChecker: repo,
		AddressCreater:        repo,
		AddressUpdater:        repo,
	}
}
