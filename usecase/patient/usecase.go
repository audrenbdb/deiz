package patient

type repo interface {
	Searcher
	Creater
	Updater
	ClinicianBoundChecker
	AddressCreater
	AddressUpdater
	BookingsGetter
}

type Usecase struct {
	Searcher              Searcher
	Creater               Creater
	Updater               Updater
	ClinicianBoundChecker ClinicianBoundChecker
	AddressCreater        AddressCreater
	AddressUpdater        AddressUpdater
	BookingsGetter        BookingsGetter
}

func NewUsecase(repo repo) *Usecase {
	return &Usecase{
		Searcher:              repo,
		Creater:               repo,
		ClinicianBoundChecker: repo,
		AddressCreater:        repo,
		AddressUpdater:        repo,
		Updater:               repo,
		BookingsGetter:        repo,
	}
}
