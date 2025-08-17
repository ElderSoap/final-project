package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler(
	id		INTEGER 	 NOT NULL PRIMARY KEY AUTOINCREMENT DEFAULT 0,
	date 	CHAR(8) 	 NOT NULL DEFAULT '',
	title 	VARCHAR(256) NOT NULL DEFAULT ('""') ,
	comment TEXT 		 NOT NULL DEFAULT ('""') ,
	repeat 	VARCHAR(128) NOT NULL DEFAULT ('""') 
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

// Открывает или создвет БД
func InitDB(dbFile string) error {
	//проверяем, что файл БД существует
	var install bool
	if _, err := os.Stat(dbFile); err != nil {
		install = true
		if dir := filepath.Dir(dbFile); dir != "" && dir != "." {
			if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
				return fmt.Errorf("mkdir %s: %w", dir, mkErr)
			}
		}
	}
	d, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}

	if err := d.Ping(); err != nil {
		_ = d.Close()
		return fmt.Errorf("ping sqlite: %w", err)
	}

	if install {
		if _, err := d.Exec(schema); err != nil {
			_ = d.Close()
			return fmt.Errorf("apply schema: %w", err)
		}
	}

	// Сохраняем соединение в глобальную переменную
	DB = d
	return nil
}

func Handle() *sql.DB { return DB }

func Close() error {
	if DB != nil {
		err := DB.Close()
		DB = nil // обнуляем, чтобы не использовать закрытое подключение
		return err
	}
	return nil
}
