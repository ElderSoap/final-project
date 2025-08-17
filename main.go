package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Запуск сервера...")
	webDir := "web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX

	addr := ":7540"
	fmt.Printf("Сервер запущен: http://localhost%s/\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
