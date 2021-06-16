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

func (e *echoServer)

func RegisterService(e *echo.Echo, getCredentials auth.GetCredentialsFromHttpRequest, db psql.PGX) {
	repo := registerPsqlRepo(db)
	h := echoServer{}
	handleGetClinicianCalendar := h.handleGetClinicianWeek(getCredentials, registerGetClinicianWeekUsecase(repo))
	handleGetClinicianBookings := h.handleGetClinicianBookings(handleGetClinicianCalendar)

	e.GET("/api/bookings", h.handleGetBookings(getCredentials, handleGetClinicianBookings))
}
