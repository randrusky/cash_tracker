package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"cash_tracker/v2/logic" // Замініть на ваш шлях
)

// CreateEditWindow створює та показує вікно для редагування запису
func CreateEditWindow(app fyne.App, record logic.Record, onSave func()) {
	editWin := app.NewWindow(fmt.Sprintf("Редагування запису #%d", record.ID))
	editWin.Resize(fyne.NewSize(450, 600))

	entries := make(map[int]*widget.Entry)
	grandTotalLabel := widget.NewLabel("")
	grandTotalLabel.TextStyle.Bold = true

	rowsContainer := container.NewVBox()

	updateTotals := func() {
		grandTotal := 0
		for _, denom := range logic.Denominations {
			count, err := strconv.Atoi(entries[denom].Text)
			if err != nil {
				count = 0
			}
			grandTotal += count * denom
		}
		grandTotalLabel.SetText(fmt.Sprintf("Загальна сума: %d грн", grandTotal))
	}

	for _, denom := range logic.Denominations {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("0")
		// Встановлюємо початкові значення з запису, що редагується
		if count, ok := record.Breakdown[denom]; ok {
			entry.SetText(strconv.Itoa(count))
		}
		entries[denom] = entry
		entry.OnChanged = func(s string) { updateTotals() }

		rowsContainer.Add(container.NewHBox(
			widget.NewLabel(fmt.Sprintf("%d грн:", denom)),
			entry,
		))
	}
	
	updateTotals() // Перший розрахунок сум

	saveButton := widget.NewButton("Зберегти зміни", func() {
		newBreakdown := make(map[int]int)
		newTotal := 0
		for denom, entry := range entries {
			count, _ := strconv.Atoi(entry.Text)
			newBreakdown[denom] = count
			newTotal += count * denom
		}
		record.Breakdown = newBreakdown
		record.Total = newTotal

		if err := logic.UpdateRecord(record); err != nil {
			// Тут можна показати діалог помилки
			fmt.Printf("Помилка оновлення запису #%d: %v\n", record.ID, err)
			return
		}
		onSave()       // Викликаємо callback, щоб оновити вікно статистики
		editWin.Close() // Закриваємо вікно редагування
	})

	editWin.SetContent(container.NewVBox(
		rowsContainer,
		grandTotalLabel,
		saveButton,
	))

	editWin.Show()
}