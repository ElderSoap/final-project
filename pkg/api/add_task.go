package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/eldersoap/filnal-project/pkg/db"
)

// addTaskHandler обрабатывает POST-запросы на /api/task
func addTaskHandle(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Проверяем дату и правим task.Date при необходимости
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Проверка: обязательное поле Title
	if err := validateTask(&task); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Успешный ответ
	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата не указана → используем сегодня
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	// Проверяем формат даты
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("дата должна быть в формате YYYYMMDD")
	}

	// Если указано правило повторения → проверим его
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Если дата < сегодня
	if afterNow(now, t) {
		if task.Repeat == "" {
			// без повторений → берем сегодня
			task.Date = now.Format("20060102")
		} else {
			// с повторением → берем следующую дату
			task.Date = next
		}
	}
	return nil
}

// writeJSON — вспомогательная функция для ответа JSON
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// validateTask проверяет корректность задачи
func validateTask(t *db.Task) error {
	if t.Title == "" {
		return errors.New("заголовок обязателен")
	}

	if _, err := time.Parse("20060102", t.Date); err != nil {
		return errors.New("некорректная дата, используйте формат YYYYMMDD")
	}

	if t.Repeat != "" {
		matched, _ := regexp.MatchString(`^(?:[dwmy]|[dwmy]\s+\d+)$`, t.Repeat)
		if !matched {
			return errors.New("некорректное поле повторения")
		}
	}

	return nil
}
