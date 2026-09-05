package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func lastDayOfMonth(date time.Time) int {
	firstOfNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfMonth.Day()
}

func parseWeekdays(s string) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, item := range strings.Split(s, ",") {
		wd, err := strconv.Atoi(item)
		if err != nil || wd < 1 || wd > 7 {
			return nil, errors.New("недопустимый день недели")
		}
		result[wd] = true
	}
	return result, nil
}

func parseMonthDays(s string) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, item := range strings.Split(s, ",") {
		d, err := strconv.Atoi(item)
		if err != nil {
			return nil, errors.New("недопустимый день месяца")
		}
		if d == -1 || d == -2 {
			result[d] = true
			continue
		}
		if d < 1 || d > 31 {
			return nil, errors.New("недопустимый день месяца")
		}
		result[d] = true
	}
	return result, nil
}

func parseMonths(s string) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, item := range strings.Split(s, ",") {
		m, err := strconv.Atoi(item)
		if err != nil || m < 1 || m > 12 {
			return nil, errors.New("недопустимый номер месяца")
		}
		result[m] = true
	}
	return result, nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("некорректная дата")
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат правила y")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "d":
		if len(parts) != 2 {
			return "", errors.New("не указан интервал в днях")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("интервал в днях должен быть числом")
		}
		if days < 1 || days > 400 {
			return "", errors.New("превышен допустимый интервал в днях")
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", errors.New("не указаны дни недели")
		}
		weekdays, err := parseWeekdays(parts[1])
		if err != nil {
			return "", err
		}
		for {
			date = date.AddDate(0, 0, 1)
			if !afterNow(date, now) {
				continue
			}
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}
			if weekdays[wd] {
				break
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("неверный формат правила m")
		}
		days, err := parseMonthDays(parts[1])
		if err != nil {
			return "", err
		}
		var months map[int]bool
		if len(parts) == 3 {
			months, err = parseMonths(parts[2])
			if err != nil {
				return "", err
			}
		}
		for {
			date = date.AddDate(0, 0, 1)
			if !afterNow(date, now) {
				continue
			}
			if months != nil && !months[int(date.Month())] {
				continue
			}
			day := date.Day()
			last := lastDayOfMonth(date)
			if days[day] || (days[-1] && day == last) || (days[-2] && day == last-1) {
				break
			}
		}

	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}

	return date.Format(dateFormat), nil
}
