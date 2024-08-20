package util

import "time"

func FirstDayOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 0, -t.Day()+1)
}

func LastDayOfMonth(t time.Time) time.Time {
	firstDayNextMonth := t.AddDate(0, 1, -t.Day()+1)
	return firstDayNextMonth.AddDate(0, 0, -1)
}

func FirstDayOfWeek(t time.Time) time.Time {
	offset := int(t.Weekday()) // Sunday = 0, Saturday = 6
	return t.AddDate(0, 0, -offset)
}

func LastDayOfWeek(t time.Time) time.Time {
	offset := int(t.Weekday())
	return t.AddDate(0, 0, 6-offset)
}
