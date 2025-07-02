package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/randrusky/cash_tracker/logic" // Замініть на ваш шлях до проєкту
)

// CreateMainWindow створює та повертає головне вікно
func CreateMainWindow(app fyne.App) fyne.Window {
	win := app.NewWindow("Облік Готівки")
	win.Resize(fyne.NewSize(400, 600))

	entries := make(map[int]*widget.Entry)
	totalLabel := widget.NewLabel("Загальна сума: 0 грн")

	// Функція для оновлення загальної суми
	updateTotal := func() {
		total := 0
		for _, denom := range logic.Denominations {
			count, err := strconv.Atoi(entries[denom].Text)
			if err == nil {
				total += count * denom
			}
		}
		totalLabel.SetText(fmt.Sprintf("Загальна сума: %d грн", total))
	}

	formItems := []*widget.FormItem{}
	for _, denom := range logic.Denominations {
		entry := widget.NewEntry()
		entry.OnChanged = func(s string) { updateTotal() }
		entries[denom] = entry
		formItems = append(formItems, widget.NewFormItem(fmt.Sprintf("%d грн:", denom), entry))
	}

	form := widget.NewForm(formItems...)

	// Кнопка для збереження запису
	saveButton := widget.NewButton("Додати запис", func() {
		breakdown := make(map[int]int)
		total := 0
		for _, denom := range logic.Denominations {
			count, err := strconv.Atoi(entries[denom].Text)
			if err != nil {
				count = 0
			}
			breakdown[denom] = count
			total += count * denom
		}

		if total > 0 {
			record := logic.Record{
				Date:      time.Now(),
				Total:     total,
				Breakdown: breakdown,
			}
			if err := logic.SaveRecord(record); err != nil {
				// Тут можна показати вікно з помилкою
				fmt.Println("Помилка збереження:", err)
				return
			}
			// Очищення полів після збереження
			for _, entry := range entries {
				entry.SetText("")
			}
			totalLabel.SetText("Загальна сума: 0 грн")
		}
	})

	// Кнопка для відкриття статистики
	statsButton := widget.NewButton("Переглянути статистику", func() {
		statsWin := CreateStatsWindow(app)
		statsWin.Show()
	})

	win.SetContent(container.NewVBox(
		form,
		totalLabel,
		saveButton,
		statsButton,
	))

	return win
}