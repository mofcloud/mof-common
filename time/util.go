/*
 * Copyright (c) 2022 The Mof Authors
 */

package ptime

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var layoutList = []string{
	"2006-01-02",
	"2006-1-2",
	"2006-01",
	"2006-1",
	"2006-01-02T15:04Z",
	ProphetFormat,
	AwsCloudWatchFormat,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

const (
	ProphetFormat       = "2006-01-02 15:04:05"
	AwsCloudWatchFormat = "2006-01-02 15:04:05 -0700 MST"
)

// TimeToLayoutMonth converts time to YYYY-MM layout
func TimeToLayoutMonth(ts time.Time) string {
	year, month, _ := ts.Date()

	monthStr := fmt.Sprintf("%d", month)
	// convert month and day it MM or DD format
	if month < 10 {
		monthStr = fmt.Sprintf("0%d", month)
	}

	return fmt.Sprintf("%d-%s", year, monthStr)
}

// TimeToLayoutDay converts time to YYYY-MM-DD layout
func TimeToLayoutDay(ts time.Time) string {
	year, month, day := ts.Date()

	monthStr := fmt.Sprintf("%d", month)
	dayStr := fmt.Sprintf("%d", day)
	// convert month and day it MM or DD format
	if month < 10 {
		monthStr = fmt.Sprintf("0%d", month)
	}

	if day < 10 {
		dayStr = fmt.Sprintf("0%d", day)
	}

	return fmt.Sprintf("%d-%s-%s", year, monthStr, dayStr)
}

// StringToLayoutMonth try our best to convert incoming string to YYYY-MM
func StringToLayoutMonth(str string) (string, error) {
	// 1: try to convert string to time
	ts, err := StringToTime(str)
	if err == nil {
		return TimeToLayoutMonth(ts), nil
	}

	return "", errors.New(fmt.Sprintf("failed to convert to YYYY-MM layout for %s", str))
}

// StringToLayoutDaily try our best to convert incoming string to YYYY-MM-DD
func StringToLayoutDaily(str string) (string, error) {
	// 1: try to convert string to time
	ts, err := StringToTime(str)
	if err == nil {
		return TimeToLayoutDay(ts), nil
	}

	return "", errors.New(fmt.Sprintf("failed to convert to YYYY-MM-DD layout for %s", str))
}

// FirstDayOfMonthString Convert current time to YYYY-MM-01 layout
func FirstDayOfMonthString(ts time.Time) string {
	year, month, _ := ts.Date()

	monthStr := fmt.Sprintf("%d", month)
	// convert month and day it MM or DD format
	if month < 10 {
		monthStr = fmt.Sprintf("0%d", month)
	}

	return fmt.Sprintf("%d-%s-01", year, monthStr)
}

// FirstDayOfMonthTime Convert current time to YYYY-MM-01 layout
func FirstDayOfMonthTime(ts time.Time) time.Time {
	year, month, _ := ts.Date()

	monthStr := fmt.Sprintf("%d", month)
	// convert month and day it MM or DD format
	if month < 10 {
		monthStr = fmt.Sprintf("0%d", month)
	}

	res, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-%s-01", year, monthStr))
	return res
}

// LastDayOfMonthString Convert current time to YYYY-MM-30 layout
func LastDayOfMonthString(ts time.Time) string {
	currMonth := ts.Month()

	for currMonth == ts.Month() {
		ts = ts.Add(24 * time.Hour)
	}

	// we are already at next month
	// now shift left for 24 hours
	ts = ts.Add(-24 * time.Hour)

	return TimeToLayoutDay(ts)
}

func LastMonthString(str string) (string, error) {
	ts, err := StringToTime(str)
	if err != nil {
		return "", err
	}

	month := ts.Month()
	for month == ts.Month() {
		ts = ts.Add(-24 * time.Hour)
	}

	return TimeToLayoutMonth(ts), nil
}

func YesterdayString(str string) (string, error) {
	ts, err := StringToTime(str)
	if err != nil {
		return "", err
	}

	ts = ts.Add(-24 * time.Hour)

	return TimeToLayoutDay(ts), nil
}

// LastDayOfMonthTime Convert current time to YYYY-MM-30 layout
func LastDayOfMonthTime(ts time.Time) time.Time {
	currMonth := ts.Month()

	yesterday := time.Now().Add(-24 * time.Hour)

	for currMonth == ts.Month() && ts.Before(yesterday) {
		ts = ts.Add(24 * time.Hour)
	}

	if ts.Month() == currMonth {
		return ts
	}

	// we are already at next month
	// now shift left for 24 hours
	ts = ts.Add(-24 * time.Hour)

	return ts
}

func LastDayOfMonthTimeActual(ts time.Time) time.Time {
	currMonth := ts.Month()

	for currMonth == ts.Month() {
		ts = ts.Add(24 * time.Hour)
	}

	//// we are already at next month
	//// now shift left for 24 hours
	//ts = ts.Add(-24 * time.Hour)

	// move time
	newTs, _ := StringToTime(fmt.Sprintf("%d-%d-01", ts.Year(), ts.Month()))
	newTs = newTs.Add(-1 * time.Second)

	return newTs
}

// StringToTime try our best to parse string to time
func StringToTime(str string) (time.Time, error) {
	var err error
	var ts time.Time

	for _, l := range layoutList {
		ts, err = time.Parse(l, str)
		if err == nil {
			break
		}
	}

	return ts, err
}

func ToProphetFormatFromString(str string) (string, error) {
	var err error
	var ts time.Time

	for _, l := range layoutList {
		ts, err = time.Parse(l, str)
		if err == nil {
			break
		}
	}

	if err == nil {
		return ts.Format(ProphetFormat), nil
	}

	// parse as epoch time
	if raw, err := strconv.ParseInt(str, 10, 64); err != nil {
		return "", err
	} else {
		ts = time.UnixMilli(raw)
	}

	return ts.Format(ProphetFormat), nil
}

func ToStdFormatFromString(str string) (string, error) {
	var ts time.Time

	// is the incoming ts is epoch?
	if num, err := strconv.ParseInt(str, 10, 64); err == nil {
		ts = time.Unix(num, 0)
		return ts.Format(time.RFC3339Nano), nil
	} else {
		for _, l := range layoutList {
			ts, err = time.Parse(l, str)
			if err == nil {
				break
			}
		}

		if err == nil {
			return ts.Format(time.RFC3339Nano), nil
		}
	}

	return ts.Format(time.RFC3339Nano), errors.New("invalid timestamp format")
}

func ToAlibabaFormatFromString(str string) (string, error) {
	var ts time.Time

	// is the incoming ts is epoch?
	if num, err := strconv.ParseInt(str, 10, 64); err == nil {
		ts = time.UnixMilli(num)
		return ts.Format(time.RFC3339Nano), nil
	} else {
		for _, l := range layoutList {
			ts, err = time.Parse(l, str)
			if err == nil {
				break
			}
		}

		if err == nil {
			return ts.Format(time.RFC3339Nano), nil
		}
	}

	return ts.Format(time.RFC3339Nano), errors.New("invalid timestamp format")
}

func ToUCloudFormatFromString(str string) (string, error) {
	var ts time.Time

	// is the incoming ts is epoch?
	if num, err := strconv.ParseInt(str, 10, 64); err == nil {
		ts = time.Unix(num, 0)
		return ts.Format(time.RFC3339Nano), nil
	} else {
		for _, l := range layoutList {
			ts, err = time.Parse(l, str)
			if err == nil {
				break
			}
		}

		if err == nil {
			return ts.Format(time.RFC3339Nano), nil
		}
	}

	return ts.Format(time.RFC3339Nano), errors.New("invalid timestamp format")
}

// NextMonthLayoutMonthTime get next month
func NextMonthLayoutMonthTime(tsStart time.Time) time.Time {
	ts, _ := time.Parse("2006-01", NextMonthLayoutMonthString(tsStart))
	return ts
}

// NextMonthLayoutMonthString get next month and return YYYY-MM-DD layout
func NextMonthLayoutMonthString(tsStart time.Time) string {
	year := tsStart.Year()
	month := tsStart.Month()

	if month < time.December {
		month++
	} else {
		year++
		month = time.January
	}

	monthStr := fmt.Sprintf("%d", month)
	// convert month and day it MM
	if month < 10 {
		monthStr = fmt.Sprintf("0%d", month)
	}

	return fmt.Sprintf("%d-%s", year, monthStr)
}

func NextMonthLayoutFromString(str string) (string, error) {
	ts, err := StringToTime(str)
	if err != nil {
		return "", err
	}

	return NextMonthLayoutMonthString(ts), nil
}

// NextDayLayoutDayTime get next day
func NextDayLayoutDayTime(tsStart time.Time) time.Time {
	return tsStart.Add(24 * time.Hour)
}

// NextDayLayoutDayString get next day and return YYYY-MM-DD layout
func NextDayLayoutDayString(tsStart time.Time) string {
	return TimeToLayoutDay(tsStart.Add(24 * time.Hour))
}

// IsStdMonthLayout checks whether incoming string is in format of YYYY-MM
func IsStdMonthLayout(str string) bool {
	_, err := time.Parse("2006-01", str)
	if err != nil {
		return false
	}

	return true
}

// IsXStdMonthLayout checks whether incoming string is in format of YYYY-M
func IsXStdMonthLayout(str string) bool {
	_, err := time.Parse("2006-1", str)
	if err != nil {
		return false
	}

	return true
}

// IsStdDayLayout checks whether incoming string is in format of YYYY-MM-DD
func IsStdDayLayout(str string) bool {
	_, err := time.Parse("2006-01-02", str)
	if err != nil {
		return false
	}

	return true
}

// IsXStdDayLayout checks whether incoming string is in format of YYYY-M-D
func IsXStdDayLayout(str string) bool {
	_, err := time.Parse("2006-1-2", str)
	if err != nil {
		return false
	}

	return true
}

// ToStdDayLayout convert any YYYY-MM-DD or YYYY-M-D layout to YYYY-MM-DD
//
// If str is not YYYY-MM-DD or YYYY-M-D, then return false
func ToStdDayLayout(str string) (string, bool) {
	// YYYY-MM-DD
	if IsStdDayLayout(str) {
		return str, true
	}

	// YYYY-M-D
	if IsXStdDayLayout(str) {
		res, _ := StringToLayoutDaily(str)
		return res, true
	}

	return "", false
}

// ToStdMonthLayout convert any YYYY-MM or YYYY-M layout to YYYY-MM
//
// If str is not YYYY-MM or YYYY-M-D, then return false
func ToStdMonthLayout(str string) (string, bool) {
	// YYYY-MM
	if IsStdMonthLayout(str) {
		return str, true
	}

	// YYYY-M
	if IsXStdMonthLayout(str) {
		res, _ := StringToLayoutMonth(str)
		return res, true
	}

	return "", false
}

// CalcStartTimeAndEndTime
//
// # Please make sure parameters follows bellow format
//
// timestamp: YYYY-MM
// start: YYYY-MM-DD
// end: YYYY-MM-DD
//
// case 1: [startDay, endDay, firstDay, lastDay] => error
// case 2: [startDay, firstDay, endDay, lastDay] => [firstDay, endDay]
// case 3: [startDay, firstDay, lastDay, endDay] => [firstDay, lastDay]
// case 4: [firstDay, startDay, endDay, lastDay] => [startDay, endDay]
// case 5: [firstDay, startDay, lastDay, endDay] => [startDay, lastDay]
// case 6: [firstDay, lastDay, startDay, endDay] => error
func CalcStartDayAndEndDay(timestamp, start, end string) (time.Time, time.Time, error) {
	resStart, resEnd := time.Time{}, time.Time{}

	// get currMonth as time.Time
	currMonth, err := StringToTime(timestamp)
	if err != nil {
		return resStart, resEnd, err
	}

	// get firstDay of currMonth
	firstDay := FirstDayOfMonthTime(currMonth)
	// get lastDay of currMonth
	lastDay := LastDayOfMonthTime(currMonth)
	// get startDay as time.Time
	startDay, err := StringToTime(start)
	if err != nil {
		return resStart, resEnd, err
	}
	// get endDay as time.Time
	endDay, err := StringToTime(end)
	if err != nil {
		return resStart, resEnd, err
	}

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

// DaysInMonth
// Calculate how many days in current month.
//
// month should follow layout of YYYY-MM
func DaysInMonth(month string) (int, error) {
	currMonth, err := StringToTime(month)
	if err != nil {
		return 0, err
	}

	// get first day of time.Time
	firstDayOfCurrMonth := FirstDayOfMonthTime(currMonth)

	// get last day of time.Time
	lastDayOfCurrMonth := LastDayOfMonthTime(currMonth)

	daysInCurrMonth := 1
	// calculate number of days in current month
	for lastDayOfCurrMonth.After(firstDayOfCurrMonth) {
		daysInCurrMonth++
		firstDayOfCurrMonth = firstDayOfCurrMonth.Add(24 * time.Hour)
	}

	return daysInCurrMonth, nil
}

func MonthToTimePeriodDaily(month string) *TimePeriod {
	tp := &TimePeriod{
		Start: fmt.Sprintf("%s-01", month),
	}

	ts, _ := StringToTime(tp.Start)
	ts = LastDayOfMonthTimeActual(ts)
	tp.End = fmt.Sprintf("%s-%d", month, ts.Day())

	return tp
}
