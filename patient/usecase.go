package patient

type repo interface {
	Searcher
	Creater
	Updater
	ClinicianBoundChecker
	AddressCreater
	AddressUpdater
	BookingsGetter
	GetterByEmail
}

type Usecase struct {
	Searcher              Searcher
	Creater               Creater
	Updater               Updater
	GetterByEmail         GetterByEmail
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
		GetterByEmail:         repo,
	}
}
