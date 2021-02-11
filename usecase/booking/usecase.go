package booking

type repo interface {
	ClinicianOfficeHoursGetter
	BookingsInTimeRangeGetter
	Creater
	Deleter
}

type Usecase struct {
	OfficeHoursGetter         ClinicianOfficeHoursGetter
	BookingsInTimeRangeGetter BookingsInTimeRangeGetter
	Creater                   Creater
	Deleter                   Deleter
}

func NewUsecase(repo repo) *Usecase {
	return &Usecase{
		OfficeHoursGetter:         repo,
		BookingsInTimeRangeGetter: repo,
		Creater:                   repo,
		Deleter:                   repo,
	}
}
