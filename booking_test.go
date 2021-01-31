package deiz

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestFillTimeRangeWithBookings(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour * time.Duration(1))

	motive := BookingMotive{
		Duration: 15,
	}

	//should fit 4 slot
	b := fillTimeRangeWithFreeSlots(end, start, Clinician{}, []Booking{}, Address{}, motive)
	assert.Equal(t, 4, len(b))
}

func TestGetOfficeHoursTimeRange(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Paris")
	start := time.Date(2021, 1, 1, 12, 0, 0, 0, loc)
	end := start.Add(time.Hour * time.Duration(24*6))

	h := OfficeHours{
		StartMn: 800,
		EndMn:   900,
		WeekDay: 2,
	}
	expectedOfficeHoursStart := time.Date(2021, 1, 5, 800/60, 800%60, 0, 0, loc)
	expectedOfficeHoursEnd := time.Date(2021, 1, 5, 900/60, 900%60, 0, 0, loc)

	r1, r2 := getOfficeHoursTimeRange(start, end, h, loc)
	assert.Equal(t, expectedOfficeHoursStart, r1)
	assert.Equal(t, expectedOfficeHoursEnd, r2)
}

func TestFillOfficeHoursWithFreeSlots(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Paris")
	start := time.Date(2021, 1, 1, 12, 0, 0, 0, loc)
	end := start.Add(time.Hour * time.Duration(24*6))

	h := []OfficeHours{{
		StartMn: 800,
		EndMn:   860,
		WeekDay: 2,
	}}

	b := fillOfficeHoursWithFreeSlots(start, end, Clinician{}, h, []Booking{}, BookingMotive{
		Duration: 15,
	}, loc)
	assert.Len(t, b, 4)
}

func TestLimitRange(t *testing.T) {
	lowerLimit := time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)
	upperLimit := time.Date(2021, 1, 1, 18, 0, 0, 0, time.UTC)

	rangeStart := time.Date(2021, 1, 1, 11, 0, 0, 0, time.UTC)
	rangeEnd := time.Date(2021, 1, 1, 19, 0, 0, 0, time.UTC)

	//should bound to limits
	start, end := limitTimeRange(lowerLimit, upperLimit, rangeStart, rangeEnd)
	assert.Equal(t, lowerLimit, start)
	assert.Equal(t, upperLimit, end)

	//should bound to upper limits
	rangeStart = time.Date(2021, 1, 1, 13, 0, 0, 0, time.UTC)
	start, end = limitTimeRange(lowerLimit, upperLimit, rangeStart, rangeEnd)
	assert.Equal(t, rangeStart, start)
	assert.Equal(t, upperLimit, end)

	//should not change range
	rangeEnd = time.Date(2021, 1, 1, 15, 0, 0, 0, time.UTC)
	start, end = limitTimeRange(lowerLimit, upperLimit, rangeStart, rangeEnd)
	assert.Equal(t, rangeStart, start)
	assert.Equal(t, rangeEnd, end)
}

func TestRemoveOverlappingFreeSlots(t *testing.T) {
	b := Booking{
		Start: time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
		End:   time.Date(2021, 1, 1, 13, 0, 0, 0, time.UTC),
	}

	freeSlots := []Booking{
		{
			Start: time.Date(2021, 1, 1, 11, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 1, 1, 13, 0, 0, 0, time.UTC),
		},
	}
	bookings := removeOverlappingFreeSlots(b, []Booking{b}, freeSlots)
	//should remove booking slot
	assert.Len(t, bookings, 1)

	//should remove only one free slot
	freeSlots = append(freeSlots, Booking{
		Start: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
		End:   time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
	})
	bookings = removeOverlappingFreeSlots(b, []Booking{b}, freeSlots)
	assert.Len(t, bookings, 2)
}

type fakeBookingGetter struct{}

func (r *fakeBookingGetter) GetBookingsInTimeRange(ctx context.Context, from, to time.Time, clinicianID int) ([]Booking, error) {
	return []Booking{
		{
			Start: time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 1, 1, 13, 0, 0, 0, time.UTC),
		},
	}, nil
}

func (r *fakeBookingGetter) GetOfficeHours(ctx context.Context, clinicianID int) ([]OfficeHours, error) {
	return []OfficeHours{
		{
			WeekDay: 2,
			StartMn: 200,
			EndMn:   1200,
		},
	}, nil
}

func (r *fakeBookingGetter) GetCalendarSettings(ctx context.Context, clinicianID int) (CalendarSettings, error) {
	return CalendarSettings{
		Timezone:      Timezone{Name: "Europe/Paris"},
		DefaultMotive: BookingMotive{Duration: 30},
	}, nil
}

func TestGetAllBookingSlotsFromWeekFunc(t *testing.T) {
	r := &fakeBookingGetter{}
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	getBookingsFunc := getAllBookingSlotsFromWeekFunc(r, r, r)
	_, err := getBookingsFunc(context.Background(), start, 1)
	assert.NoError(t, err)
}

func TestGetFreeBookingSlotsFromWeekFunc(t *testing.T) {
	r := &fakeBookingGetter{}
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	getBookingsFunc := getFreeBookingSlotsFromWeekFunc(r, r, r)
	b, err := getBookingsFunc(context.Background(), start, 1)
	assert.NoError(t, err)
	log.Println(b)
}
