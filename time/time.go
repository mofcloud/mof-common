/*
 * Copyright (c) 2022 The Mof Authors
 */

package ptime

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	Monthly = "monthly"
	Daily   = "daily"
)

// ************************************
// ************ TimePeriod ************
// ************************************

// GetCurrMonthTimePeriod return current month time period
func GetCurrMonthTimePeriod() *TimePeriod {
	now := time.Now()

	return &TimePeriod{
		Start: FirstDayOfMonthString(now),
		End:   TimeToLayoutDay(now),
	}
}

// GetThisYearAndLastDecTimePeriod Get the time period from last December 1st to now
func GetThisYearAndLastDecTimePeriod(now time.Time) *TimePeriod {
	return &TimePeriod{
		Start: fmt.Sprintf("%d-12-01", now.Year()-1),
		End:   TimeToLayoutDay(now),
	}
}

func GetLastXMonthTimePeriod(now time.Time, x int) *TimePeriod {
	year := now.Year()
	month := now.Month() - time.Month(x)

	for month < 1 {
		month = 12 + month
		year--
	}

	format := "%d-"
	if month < 10 {
		format = format + "0%d-"
	} else {
		format = format + "%d-"
	}
	format = format + "01"

	return &TimePeriod{
		Start: fmt.Sprintf(format, year, month),
		End:   TimeToLayoutDay(now),
	}
}

// TimePeriod
// Start: should be YYYY-MM-DD
// End: should be YYYY-MM-DD
type TimePeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// Validate valid time period should be format of YYYY-MM-DD
func (t *TimePeriod) Validate() error {
	// 1: invalid time format of start time
	if _, err := StringToLayoutDaily(t.Start); err != nil {
		return errors.New(fmt.Sprintf("invalid startTime: %s, should be format YYYY-MM-DD", t.Start))
	}

	// 2: invalid format of end time
	if _, err := StringToLayoutDaily(t.End); err != nil {
		return errors.New(fmt.Sprintf("invalid endTime: %s, should be format of YYYY-MM-DD", t.End))
	}

	startTime, _ := StringToTime(t.Start)
	endTime, _ := StringToTime(t.End)

	// 3: start time is after today
	if startTime.After(time.Now()) {
		return errors.New("start time is out of range, should be before today")
	}

	// 4: start time is after end time
	if startTime.Equal(endTime) || startTime.After(endTime) {
		return errors.New(fmt.Sprintf("invalid startTime: %s endTime: %s, startTime is not before endTime",
			t.Start, t.End))

	}

	return nil
}

// ToMonthList parse month list
func (t *TimePeriod) ToMonthList() []string {
	res := make([]string, 0)

	if t.Start == "" || t.End == "" {
		return res
	}

	tsStart, _ := StringToTime(t.Start)
	tsEnd, _ := StringToTime(t.End)

	for !tsEnd.Before(tsStart) {
		res = append(res, TimeToLayoutMonth(tsStart))
		tsStart = NextMonthLayoutMonthTime(tsStart)
	}

	return res
}

func (t *TimePeriod) ToMonthListForGCP() []string {
	res := t.ToMonthList()

	for i := range res {
		res[i] = strings.Replace(res[i], "-", "", -1)
	}

	return res
}

func (t *TimePeriod) ToDayList() []string {
	res := make([]string, 0)

	if t.Start == "" || t.End == "" {
		return res
	}

	startTime := t.Start
	endTime := t.End

	// check element
	if IsStdMonthLayout(t.Start) {
		startTime = fmt.Sprintf("%s-01", startTime)
	}

	if IsStdMonthLayout(endTime) {
		now := time.Now()
		currMonth := TimeToLayoutMonth(now)

		if endTime == currMonth {
			// set to current date
			endTime = TimeToLayoutDay(now)
		} else {
			ts, _ := StringToTime(endTime)
			endTime = LastDayOfMonthString(ts)
		}
	}

	tsStart, _ := StringToTime(startTime)
	tsEnd, _ := StringToTime(endTime)

	for !tsEnd.Before(tsStart) {
		res = append(res, TimeToLayoutDay(tsStart))
		tsStart = NextDayLayoutDayTime(tsStart)
	}

	return res
}

// ToStartTime convert to time.Time
func (t *TimePeriod) ToStartTime() time.Time {
	start, _ := StringToTime(t.Start)
	return start
}

// ToEndTime convert to time.Time
func (t *TimePeriod) ToEndTime() time.Time {
	end, _ := StringToTime(t.End)
	return end
}

// NumOfDays how may days in current month
func (t *TimePeriod) NumOfDays(month string) (int, error) {
	startDay, endDay, err := t.StartAndEndInMonth(month)
	if err != nil {
		return 0, err
	}

	return int(endDay.Sub(startDay).Hours() / 24), nil
}

// InRange is current timestamp in range?
func (t *TimePeriod) InRange(ts string) bool {
	if v, err := StringToTime(ts); err != nil {
		return false
	} else {
		return !v.Before(t.ToStartTime()) && !v.After(t.ToEndTime())
	}
}

// StartAndEndInMonth
// Please make sure parameters follows bellow format
//
// month: YYYY-MM
// start: YYYY-MM-DD
// end: YYYY-MM-DD
//
// case 1: [startDay, endDay, firstDay, lastDay] => error
// case 2: [startDay, firstDay, endDay, lastDay] => [firstDay, endDay]
// case 3: [startDay, firstDay, lastDay, endDay] => [firstDay, lastDay]
// case 4: [firstDay, startDay, endDay, lastDay] => [startDay, endDay]
// case 5: [firstDay, startDay, lastDay, endDay] => [startDay, lastDay]
// case 6: [firstDay, lastDay, startDay, endDay] => error
func (t *TimePeriod) StartAndEndInMonth(month string) (time.Time, time.Time, error) {
	// calculate startTime and endTime
	resStart, resEnd := time.Time{}, time.Time{}

	// get currMonth as time.Time
	currMonth, err := StringToTime(month)
	if err != nil {
		return resStart, resEnd, err
	}

	firstDay := FirstDayOfMonthTime(currMonth)
	lastDay := LastDayOfMonthTime(currMonth)
	startDay := t.ToStartTime()
	endDay := t.ToEndTime()

	firstDayNano := firstDay.UnixNano()
	lastDayNano := lastDay.UnixNano()
	startDayNano := startDay.UnixNano()
	endDayNano := endDay.UnixNano()

	errMsg := fmt.Sprintf("failed to calculate startDay and EndDay, firstDay:%s, lastDay:%s, startDay:%s, endDay:%s",
		firstDay.String(), lastDay.String(), startDay.String(), endDay.String())

	if startDayNano > endDayNano {
		return resStart, resEnd, errors.New(errMsg)
	}

	if endDayNano <= firstDayNano {
		// case 1: [startDay, endDay, firstDay, lastDay] => error
		return resStart, resEnd, errors.New(errMsg)
	} else if startDayNano <= firstDayNano && firstDayNano <= endDayNano && endDayNano <= lastDayNano {
		// case 2: [startDay, firstDay, endDay, lastDay] => [firstDay, endDay]
		return firstDay, endDay, nil
	} else if startDayNano <= firstDayNano && lastDayNano <= endDayNano {
		// case 3: [startDay, firstDay, lastDay, endDay] => [firstDay, lastDay]
		return firstDay, lastDay, nil
	} else if firstDayNano <= startDayNano && endDayNano <= lastDayNano {
		// case 4: [firstDay, startDay, endDay, lastDay] => [startDay, endDay]
		return startDay, endDay, nil
	} else if firstDayNano <= startDayNano && lastDayNano <= endDayNano {
		// case 5: [firstDay, startDay, lastDay, endDay] => [startDay, lastDay]
		return startDay, lastDay, nil
	} else if lastDayNano <= startDayNano {
		// case 6: [firstDay, lastDay, startDay, endDay] => error
		return resStart, resEnd, errors.New(errMsg)
	}

	return resStart, resEnd, errors.New(errMsg)
}

func (t *TimePeriod) OneMonthBefore() error {
	lastMonthStr, err := LastMonthString(t.Start)
	if err != nil {
		return err
	}
	t.Start = lastMonthStr
	return nil
}

func (t *TimePeriod) String() string {
	return fmt.Sprintf("%s->%s", t.Start, t.End)
}
