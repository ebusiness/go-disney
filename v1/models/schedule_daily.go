package models

import (
	"time"
)

// ScheduleDaily -
type ScheduleDaily struct {
	collectionName string     `collectionName:"schedule_daily"`
	StrID          string     `json:"str_id" bson:"str_id"`
	Schedules      []Schedule `json:"schedules,omitempty" bson:"schedules,omitempty"`
}

// Schedule -
type Schedule struct {
	StartTime time.Time `json:"startTime,omitempty" bson:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty" bson:"endTime,omitempty"`
}

// IsConflict -
func (s Schedule) IsConflict(s1 Schedule) bool {

	if s.StartTime.After(s1.EndTime) {
		return false
	}
	if s.EndTime.Before(s1.StartTime) {
		return false
	}
	return true
}
