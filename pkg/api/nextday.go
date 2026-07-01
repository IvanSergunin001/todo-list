package api

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		err := errors.New("method not allowed: must be GET")
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	now := req.URL.Query().Get("now")
	var targetDate time.Time
	if now == "" {
		targetDate = time.Now()
		now = targetDate.Format("20060102")
	}

	date := req.URL.Query().Get("date")
	repeatLetter := req.URL.Query().Get("repeat")
	if repeatLetter == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(""))
		return
	}

	newDate, err := NextDate(date, now, repeatLetter)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(""))
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(newDate))
}

func NextDate(dstart, now, repeat string) (string, error) {
	letterCheck := []string{"d", "y", "w"}

	daysOnWeek := map[string]int{
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
		"Sunday":    7,
	}

	if repeat == "" {
		return "", nil
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}
	timeNow, err := time.Parse("20060102", now)
	if err != nil {
		return "", err
	}

	var interval int

	splitRepeat := strings.Split(repeat, " ")
	if slices.Contains(letterCheck, splitRepeat[0]) {
		if splitRepeat[0] == "d" {
			if len(splitRepeat) == 2 {
				number, err := strconv.Atoi(splitRepeat[1])
				if err != nil {
					return "", err
				}
				if number <= 400 && number > 0 {
					interval = number
					for {
						date = date.AddDate(0, 0, interval)
						if afterNow(date, timeNow) {
							break
						}
					}
				} else {
					return "", nil
				}
			} else {
				return "", nil
			}
		} else if splitRepeat[0] == "y" {
			if len(splitRepeat) == 1 {
				interval = 1
				for {
					date = date.AddDate(1, 0, 0)
					if afterNow(date, timeNow) {
						break
					}
				}
			} else {
				return "", nil
			}
		} else if splitRepeat[0] == "w" {
			interval, err := weekInterval(date, daysOnWeek, splitRepeat)
			if err != nil {
				return "", err
			}
			for {
				date = date.AddDate(0, 0, interval)
				if afterNow(date, timeNow) {
					break
				} else {
					interval, err = weekInterval(date, daysOnWeek, splitRepeat)
					if err != nil {
						return "", err
					}
				}
			}
		} else {
			return "", nil
		}
	} else {
		return "", nil
	}
	return date.Format("20060102"), nil
}

func afterNow(date, timeNow time.Time) bool {
	return date.After(timeNow)
}

func weekInterval(date time.Time, daysOnWeek map[string]int, splitRepeat []string) (int, error) {
	if len(splitRepeat) != 2 {
		err := errors.New("method not allowed: must be POST")
		return 0, err
	}
	weekDay := date.Weekday()
	weekDayNum := daysOnWeek[weekDay.String()]
	days := strings.Split(splitRepeat[1], ",") // проблема
	numsStartWeek := make([]int, 0, len(days))
	numsEndWeek := make([]int, 0, len(days))
	var numsWeek int

	for _, s := range days {
		num, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		if num < 1 || num > 7 {
			err = errors.New("method not allowed: must be POST")
			return 0, err
		}

		if num < weekDayNum {
			numsStartWeek = append(numsStartWeek, num)
			continue
		} else if num > weekDayNum {
			numsEndWeek = append(numsEndWeek, num)
			continue
		} else if num == weekDayNum {
			numsWeek = num
			continue
		}
	}

	var totalDay int

	if len(numsEndWeek) != 0 {
		totalDay = slices.Min(numsEndWeek)
	} else if len(numsStartWeek) != 0 {
		totalDay = slices.Min(numsStartWeek)
	} else if len(numsStartWeek) == 0 && len(numsEndWeek) == 0 {
		totalDay = numsWeek
	} else {
		// ошибка: нет корректных дней недели
	}

	interval := (totalDay - weekDayNum + 7) % 7

	if interval == 0 {
		interval = 7
	}
	return interval, nil
}
