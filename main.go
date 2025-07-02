package main

import (
	"fyne.io/fyne/v2/app"
	"os"
	"github.com/randrusky/cash_tracker/ui" // Замініть на ваш шлях
)

func main() {
	// Створюємо папку для даних, якщо її не існує
	_ = os.Mkdir("data", 0755)

	a := app.New()
	win := ui.CreateMainWindow(a)

	win.ShowAndRun()
}