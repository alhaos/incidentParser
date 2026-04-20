package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var MonthMap = map[string]time.Month{
	"Jan": time.January,
	"Feb": time.February,
	"Mar": time.March,
	"Apr": time.April,
	"May": time.May,
	"Jun": time.June,
	"Jul": time.July,
	"Aug": time.August,
	"Sep": time.September,
	"Oct": time.October,
	"Nov": time.November,
	"Dec": time.December,
}

func ParseDataType1(text string) (time.Time, error) {

	fields := strings.Fields(text)

	month, exists := MonthMap[fields[0]]
	if !exists {
		return time.Time{}, fmt.Errorf("month does not exist")
	}

	day, err := strconv.Atoi(strings.TrimSuffix(fields[1], ","))
	if err != nil {
		return time.Time{}, err
	}

	year, err := strconv.Atoi(fields[2])
	if err != nil {
		return time.Time{}, err
	}

	timeFields := strings.Split(fields[3], ":")
	if len(timeFields) != 3 {
		return time.Time{}, fmt.Errorf("timeFields length is not 3, got %v", timeFields)
	}

	hours, err := strconv.Atoi(timeFields[0])
	if err != nil {
		return time.Time{}, err
	}

	mins, err := strconv.Atoi(timeFields[1])
	if err != nil {
		return time.Time{}, err
	}

	secs, err := strconv.Atoi(timeFields[2])
	if err != nil {
		return time.Time{}, err
	}

	var isPm bool
	switch fields[4] {
	case "PM":
		isPm = true
	case "AM":
		isPm = false
	default:
		return time.Time{}, fmt.Errorf("invalid AM / PM format %v", fields[4])
	}

	hours = convert12hourTo24Compact(hours, isPm)

	return time.Date(year, month, day, hours, mins, secs, 0, time.Local), nil
}

func convert12hourTo24Compact(hour int, isPm bool) int {
	if hour == 12 {
		if isPm {
			return 12
		}
		return 0
	}
	if isPm {
		return hour + 12
	}
	return hour
}
