package main

import (
	"cash_tracker/v2/logic"
	"cash_tracker/v2/ui" // Замініть на ваш шлях
	"log"

	"fyne.io/fyne/v2/app"
)

const dbPath = "cash_tracker.db"

func main() {
	// Ініціалізуємо базу даних при старті
	if err := logic.InitDB(dbPath); err != nil {
		log.Fatalf("Помилка ініціалізації БД: %v", err)
	}
	defer logic.DB.Close() // Закриваємо з'єднання при виході з програми

	a := app.New()
	win := ui.CreateMainWindow(a)

	win.ShowAndRun()
}