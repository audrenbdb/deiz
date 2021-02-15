package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	ClinicianOfficeHoursGetter interface {
		GetClinicianOfficeHours(ctx context.Context, clinicianID int) ([]deiz.OfficeHours, error)
	}
)

//GetAllOfficeHoursTimeRanges gets a list of time ranges from a list of office hours
func GetAllOfficeHoursTimeRange(start, end time.Time, hours []deiz.OfficeHours, loc *time.Location) [][2]time.Time {
	var timeRanges [][2]time.Time
	for _, h := range hours {
		rangeStart, rangeEnd := GetOfficeHoursTimeRange(start, end, h, loc)
		timeRanges = append(timeRanges, [2]time.Time{rangeStart, rangeEnd})
	}
	return timeRanges
}

//GetOfficeHoursTimeRange converts generic office hours into time range within a given time range
func GetOfficeHoursTimeRange(anchor, end time.Time, h deiz.OfficeHours, loc *time.Location) (time.Time, time.Time) {
	anchorInLoc := anchor.In(loc)
	y, m, d := anchorInLoc.Date()
	if int(anchorInLoc.Weekday()) == h.WeekDay {
		officeOpensAt := time.Date(y, m, d, h.StartMn/60, h.StartMn%60, 0, 0, loc).UTC()
		officeClosesAt := time.Date(y, m, d, h.EndMn/60, h.EndMn%60, 0, 0, loc).UTC()
		return LimitTimeRange(anchor, end, officeOpensAt, officeClosesAt)
	}
	//abort if above given time range
	if anchor.After(end) {
		return time.Time{}, time.Time{}
	}
	nextAnchor := time.Date(y, m, d+1, 0, 0, 0, 0, loc).UTC()
	return GetOfficeHoursTimeRange(nextAnchor, end, h, loc)
}
