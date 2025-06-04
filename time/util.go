/*
 * Copyright (c) 2022 The Mof Authors
 */

package ptime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	LinodeFormat,
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
	LinodeFormat        = "2006-01-02T15:04:05"
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

func OneMonthBeforeTime(in time.Time) time.Time {
	month := in.Month()

	for month == in.Month() {
		in = in.Add(-24 * time.Hour)
	}

	return in
}

func OneMonthBeforeCertainDayTime(in time.Time, day int) time.Time {
	in = OneMonthBeforeTime(in)

	for in.Day() != day {
		in = in.Add(-24 * time.Hour)
	}

	return in
}

func ListDateForNaturalMonthSet(date string, monthCount int) []string {
	res := make([]string, 0)

	ts, err := StringToTime(date)
	if err != nil {
		return res
	}

	day := ts.Day()
	month := ts.Month()
	isLastDay := false
	if ts.Add(24*time.Hour).Month() != month {
		isLastDay = true
	}

	if day < 28 {
		// common day, iterate and return date list
		for i := 0; i < monthCount; i++ {
			ts = OneMonthBeforeCertainDayTime(ts, day)
			res = append(res, TimeToLayoutDay(ts))
		}
	} else {
		// check the months
		for i := 0; i < monthCount; i++ {
			ts = OneMonthBeforeTime(ts)
			newMonth := ts.Month()
			newMonthStr := TimeToLayoutMonth(ts)
			lastDayOfNewTs := LastDayOfMonthTime(ts)

			switch day {
			case 30:
				// if last day:
				// 		1,3,5,7,8,10,12: collect 30,31
				// 		4,6,9,11: collect 30
				// 		2: do nothing
				// if not
				// 		1,3,5,7,8,10,12: collect 30
				// 		4,6,9,11: collect 30
				//		2: do nothing

				if isLastDay {
					switch newMonth {
					case 1, 3, 5, 7, 8, 10, 12:
						res = append(res,
							fmt.Sprintf("%s-30", newMonthStr),
							fmt.Sprintf("%s-31", newMonthStr))
					case 4, 6, 9, 11:
						res = append(res,
							fmt.Sprintf("%s-30", newMonthStr))
					}
				} else {
					switch newMonth {
					case 1, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12:
						res = append(res,
							fmt.Sprintf("%s-30", newMonthStr))
					}
				}
			case 31:
				// 1,3,5,7,8,10,12: collect 31
				// 4,6,9,11: do nothing
				// 2: do nothing
				switch newMonth {
				case 1, 3, 5, 7, 8, 10, 12:
					res = append(res,
						fmt.Sprintf("%s-31", newMonthStr))
				}
			case 29:
				// if month == 2(which means last day):
				// 		1,3,5,7,8,10,12: collect 29,30,31
				// 		4,6,9,11: collect 29,30
				// 		2: if newMonth last day == 29
				//				collect: 29
				//		   else:
				//		        collect: 28
				// else
				//      1,3,4,5,6,7,8,9,10,11,12: collect 29
				//		2: if newMonth last day == 29
				//				collect: 29
				//		   else:
				//		        do nothing
				if month == 2 {
					switch newMonth {
					case 1, 3, 5, 7, 8, 10, 12:
						res = append(res,
							fmt.Sprintf("%s-29", newMonthStr),
							fmt.Sprintf("%s-30", newMonthStr),
							fmt.Sprintf("%s-31", newMonthStr))
					case 4, 6, 9, 11:
						res = append(res,
							fmt.Sprintf("%s-29", newMonthStr),
							fmt.Sprintf("%s-30", newMonthStr))
					case 2:
						if lastDayOfNewTs.Day() == 29 {
							res = append(res,
								fmt.Sprintf("%s-29", newMonthStr))
						} else {
							res = append(res,
								fmt.Sprintf("%s-28", newMonthStr))
						}
					}
				} else {
					switch newMonth {
					case 1, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12:
						res = append(res,
							fmt.Sprintf("%s-29", newMonthStr))
					case 2:
						if lastDayOfNewTs.Day() == 29 {
							res = append(res,
								fmt.Sprintf("%s-29", newMonthStr))
						}
					}
				}
			case 28:
				// if month == 2
				//		if last day
				//			1,3,5,7,8,10,12: collect 28,29,30,31
				//       	4,6,9,11: collect 28,29,30
				// 			2: collect 28 (29 if possible)
				//      else
				//			collect 28
				// else
				//		collect 28
				if month == 2 {
					if isLastDay {
						switch newMonth {
						case 1, 3, 5, 7, 8, 10, 12:
							res = append(res,
								fmt.Sprintf("%s-28", newMonthStr),
								fmt.Sprintf("%s-29", newMonthStr),
								fmt.Sprintf("%s-30", newMonthStr),
								fmt.Sprintf("%s-31", newMonthStr))
						case 4, 6, 9, 11:
							res = append(res,
								fmt.Sprintf("%s-28", newMonthStr),
								fmt.Sprintf("%s-29", newMonthStr),
								fmt.Sprintf("%s-30", newMonthStr))
						case 2:
							res = append(res,
								fmt.Sprintf("%s-28", newMonthStr))
							if LastDayOfMonthTime(ts).Day() == 29 {
								res = append(res,
									fmt.Sprintf("%s-29", newMonthStr))
							}
						}
					} else {
						res = append(res,
							fmt.Sprintf("%s-28", newMonthStr))
					}
				} else {
					res = append(res,
						fmt.Sprintf("%s-28", newMonthStr))
				}
			}
		}
	}

	return res
}

func ListDateForNaturalMonth(date string, monthAgo int) []string {
	res := make([]string, 0)

	set := ListDateForNaturalMonthSet(date, monthAgo)

	ts, err := StringToTime(date)
	if err != nil {
		return res
	}

	// month before
	for i := 0; i < monthAgo; i++ {
		ts = OneMonthBeforeTime(ts)
	}

	monthStr := TimeToLayoutMonth(ts)
	for i := range set {
		e := set[i]

		if strings.HasPrefix(e, monthStr) {
			res = append(res, e)
		}
	}

	return res
}
