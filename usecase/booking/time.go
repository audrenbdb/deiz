package booking

import (
	"github.com/audrenbdb/deiz"
	"sort"
	"time"
)

//timeRangesOverlaps checks if two time ranges overlaps
func TimeRangesOverlaps(startA, endA, startB, endB time.Time) bool {
	return startA.Before(endB) && startB.Before(endA)
}

//LimitTimeRange merge two time range such as second one cannot overlaps first one
func LimitTimeRange(lowerLimit, upperLimit, rangeStart, rangeEnd time.Time) (time.Time, time.Time) {
	if !TimeRangesOverlaps(lowerLimit, upperLimit, rangeStart, rangeEnd) {
		return time.Time{}, time.Time{}
	}
	if rangeStart.Before(lowerLimit) {
		rangeStart = lowerLimit
	}
	if rangeEnd.After(upperLimit) {
		rangeEnd = upperLimit
	}
	return rangeStart, rangeEnd
}

func SortBookingsByStart(bookings []deiz.Booking) []deiz.Booking {
	sort.SliceStable(bookings, func(i, j int) bool {
		return bookings[i].Start.Before(bookings[j].Start)
	})
	return bookings
}

//extract time ranges from bookings
func GetTimeRangesFromBookings(bookings []deiz.Booking, ranges [][2]time.Time) [][2]time.Time {
	if bookings == nil || len(bookings) == 0 {
		return ranges
	}
	ranges = append(ranges, [2]time.Time{bookings[0].Start, bookings[0].End})
	return GetTimeRangesFromBookings(bookings[1:], ranges)
}

func GetTimeRangesNotOverLapping(duration int, anchor, upperLimit time.Time, rangesToNotOverlap [][2]time.Time, notOverlappingRanges [][2]time.Time) [][2]time.Time {
	nextAnchor := anchor.Add(time.Minute * time.Duration(duration))
	for _, r := range rangesToNotOverlap {
		if TimeRangesOverlaps(anchor, nextAnchor, r[0], r[1]) {
			anchor = r[1]
			nextAnchor = anchor.Add(time.Minute * time.Duration(duration))
		}
	}
	if nextAnchor.After(upperLimit) {
		return notOverlappingRanges
	}
	return GetTimeRangesNotOverLapping(duration, nextAnchor,
		upperLimit, rangesToNotOverlap,
		append(notOverlappingRanges, [2]time.Time{anchor, nextAnchor}))
}
