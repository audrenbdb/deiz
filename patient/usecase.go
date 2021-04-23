package patient

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
