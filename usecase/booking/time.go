package booking

import "time"

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
