package entity

import "time"

type ScheduleDay string

const (
	DayToday    ScheduleDay = "today"
	DayTomorrow ScheduleDay = "tomorrow"
)

type Schedule struct {
	Day       ScheduleDay
	UpdatedAt time.Time
	ImageURL  string
	Content   string
}
