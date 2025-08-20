package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// проверяет, что дата date > now (без учёта времени)
// возвращает true, если дата больше текущей
func afterNow(date, now time.Time) bool {
	date = date.Truncate(24 * time.Hour) // обрезаем время
	now = now.Truncate(24 * time.Hour)   // обрезаем время

	return date.After(now)
}

// NextDate вычисляет следующую дату по правилу repeat
// принимает текущую дату now, начальную дату dstart и правило повторения repeat
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("повторение не указано")
	}

	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неккорректная дата начала: %w", err)
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неккорректный формат правила для дней")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", errors.New("неккорректное значение дней")
		}

		date := start
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неккорректный формат правила для лет")
		}
		date := start
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	default:
		return "", errors.New("неподдерживаемое правило повторения")
	}
}

// обработчик для обновления задачи, по правилу повторения
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}
