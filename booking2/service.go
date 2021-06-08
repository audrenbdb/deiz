package booking2

import (
	"github.com/audrenbdb/deiz/auth"
	"github.com/audrenbdb/deiz/psql"
	"github.com/labstack/echo"
)

type repo struct {
	getClinicianBookingsInTimeRange getClinicianBookingsInTimeRange
	getClinicianRecurrentBookings   getClinicianRecurrentBookings
	getClinicianOfficeHours         getClinicianOfficeHours
}

func registerGetClinicianWeekUsecase(r *repo) getClinicianWeek {
	getClinicianRecurrentBookingsInTimeRange := createGetClinicianRecurrentBookingsInTimeRange(r.getClinicianRecurrentBookings)
	getClinicianNonRecurrentBookingsInTimeRange := createGetClinicianNonRecurrentBookingsInTimeRangeFunc(r.getClinicianBookingsInTimeRange)
	getClinicianBookingsInWeek := createGetClinicianBookingsInWeek(
		getClinicianNonRecurrentBookingsInTimeRange,
		getClinicianRecurrentBookingsInTimeRange)

	getOfficeHoursInWeek := createGetOfficeHoursInWeekFunc(r.getClinicianOfficeHours)
	return createGetClinicianWeekFunc(getClinicianBookingsInWeek, getOfficeHoursInWeek)
}

func registerPsqlRepo(db psql.PGX) *repo {
	r := newPsqlRepo(db)
	return &repo{
		getClinicianBookingsInTimeRange: r.createGetClinicianBookingsInTimeRangeFunc(),
		getClinicianOfficeHours:         r.createGetClinicianOfficeHoursFunc(),
		getClinicianRecurrentBookings:   r.createGetClinicianRecurrentBookingsFunc(),
	}
}

func RegisterService(e *echo.Echo, auth auth.GetCredentialsFromHttpRequest, db psql.PGX) {
	repo := registerPsqlRepo(db)
	h := echoHandler{}
	handleGetClinicianCalendar := h.handleGetClinicianWeek(auth, registerGetClinicianWeekUsecase(repo))
	handleGetClinicianBookings := h.handleGetClinicianBookings(handleGetClinicianCalendar)

	e.GET("/api/bookings", h.handleGetBookings(auth, handleGetClinicianBookings))
}
