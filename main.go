package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eldersoap/filnal-project/pkg/api"
	db "github.com/eldersoap/filnal-project/pkg/db"
)

func main() {
	fmt.Println("Запуск сервера...")
	dbFile := "scheduler.db"
	if err := db.InitDB(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	api.Init()
	defer db.Close()

	webDir := "web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":7540"
	fmt.Printf("Сервер запущен: http://localhost%s/\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
