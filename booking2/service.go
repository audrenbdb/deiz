package booking2

import (
	"github.com/audrenbdb/deiz/psql"
)

type repo struct {
	getClinicianNonRecurrentBookingsInTimeRange getClinicianBookingsInTimeRange
	getClinicianRecurrentBookings               getClinicianRecurrentBookings
	getClinicianOfficeHours                     getClinicianOfficeHours
}

func registerGetClinicianCalendarUsecase(r *repo) getClinicianCalendar {
	getClinicianRecurrentBookingsInTimeRange := makeGetClinicianRecurrentBookingsInTimeRange(r.getClinicianRecurrentBookings)
	getClinicianBookingsInTimeRange := makeGetClinicianBookingsInTimeRange(
		r.getClinicianNonRecurrentBookingsInTimeRange,
		getClinicianRecurrentBookingsInTimeRange)

	getOfficeHoursInWeek := makeGetOfficeHoursInWeek(r.getClinicianOfficeHours)
	return makeGetClinicianCalendar(getClinicianBookingsInTimeRange, getOfficeHoursInWeek)
}

func NewPSQLRepo(db psql.PGX) *repo {
	r := psqlRepo{db: db}
	return &repo{
		getClinicianNonRecurrentBookingsInTimeRange: r.makeGetClinicianNonRecurrentBookingsInTimeRange(),
		getClinicianOfficeHours:                     r.makeGetClinicianOfficeHours(),
		getClinicianRecurrentBookings:               r.makeGetClinicianRecurrentBookings(),
	}
}

func (h *echoHandler) NewService(repo *repo) {
	handleGetClinicianCalendar := h.handleGetClinicianCalendar(registerGetClinicianCalendarUsecase(repo))
	handleGetClinicianBookings := h.HandleClinicianGetBookings(handleGetClinicianCalendar, nil)

	h.router.GET("/api/bookings", h.HandleGetBookings(handleGetClinicianBookings, nil))
}
