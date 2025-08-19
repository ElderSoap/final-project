package db

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(t *Task) (int64, error) {
	query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
