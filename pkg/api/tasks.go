package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eldersoap/filnal-project/pkg/db"
)

// Cтруктура для ответа с задачами
// используется для сериализации в JSON
// и отправки клиенту
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// обработчик для получения списка задач
// принимает GET запрос и возвращает список задач в формате JSON
// если произошла ошибка, то возвращает ошибку в формате JSON
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		writeJSON(w, err)
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}

// обработчик для получения задачи по ID
// принимает GET запрос с параметром id и возвращает задачу в формате JSON
// если произошла ошибка, то возвращает ошибку в формате JSON
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(taskID)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	writeJSON(w, task)
}

// обработчик для обновления задачи
// принимает PUT запрос с JSON телом задачи
// проверяет корректность данных и обновляет задачу в базе данных
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, `{"error":"Некорректный JSON"}`, http.StatusBadRequest)
		return
	}

	if t.ID == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	if err := validateTask(&t); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&t); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	writeJSON(w, map[string]string{})

}

// обработчик для завершения задачи
// принимает POST запрос с параметром id
// проверяет, существует ли задача и завершает её
// если задача с повторением, то обновляет дату
// если нет, то удаляет задачу из базы данных
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}
	t, err := db.GetTask(taskID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}
	if t.Repeat == "" {
		if err := db.DeleteTask(taskID); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{})
		return
	}
	// Если задача с повторением, то просто обновляем дату
	next, err := NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if err := db.UpdateDate(next, taskID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{})
}

// функция для удаления задачи
// принимает DELETE запрос с параметром id
// проверяет, существует ли задача и удаляет её
// если задача не найдена, то возвращает ошибку
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(taskID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{})
}
