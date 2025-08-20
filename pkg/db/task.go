package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// функция для добавления задачи в базу данных
// возвращает ID добавленной задачи
func AddTask(t *Task) (int64, error) {

	query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// функция для получения списка задач из базы данных
// принимает лимит на количество задач, возвращает срез задач
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date DESC LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

// функция для получения задачи по ID
// возвращает указатель на задачу
func GetTask(id string) (*Task, error) {
	t := &Task{}

	var tid int

	err := DB.QueryRow(`
        SELECT id, date, title, comment, repeat
        FROM scheduler WHERE id = ?`, id).Scan(&tid, &t.Date, &t.Title, &t.Comment, &t.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, err
	}

	t.ID = fmt.Sprintf("%d", tid)

	return t, nil
}

// функция для обновления задачи в базе данных
// принимает указатель на задачу и обновляет её данные
func UpdateTask(t *Task) error {
	query := `
        UPDATE scheduler 
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?`

	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// функция для удаления задачи из базы данных
// принимает ID задачи и удаляет её
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := DB.Exec(query, id)

	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// функция для обновления даты задачи
// принимает новую дату и ID задачи, обновляет дату задачи в базе данных
func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
