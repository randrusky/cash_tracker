package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/randrusky/cash_tracker/logic" // Замініть на ваш шлях
)

// CreateStatsWindow створює вікно статистики
func CreateStatsWindow(app fyne.App) fyne.Window {
	win := app.NewWindow("Статистика надходжень")
	win.Resize(fyne.NewSize(800, 500))

	records, err := logic.LoadRecords()
	if err != nil {
		win.SetContent(widget.NewLabel("Помилка завантаження даних: " + err.Error()))
		return win
	}

	totalIncome := 0
	for _, rec := range records {
		totalIncome += rec.Total
	}

	totalLabel := widget.NewLabel(fmt.Sprintf("Загальний дохід за весь час: %d грн", totalIncome))

	// Створення таблиці для історії
	table := widget.NewTable(
		func() (int, int) {
			// Кількість рядків і стовпців
			return len(records) + 1, 3 // +1 для заголовка
		},
		func() fyne.CanvasObject {
			// Шаблон для комірки
			return widget.NewLabel("Шаблон")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			// Оновлення даних в комірці
			label := cell.(*widget.Label)
			if id.Row == 0 { // Заголовок
				switch id.Col {
				case 0:
					label.SetText("Дата")
				case 1:
					label.SetText("Сума (грн)")
				case 2:
					label.SetText("Деталізація")
				}
				label.TextStyle.Bold = true
				return
			}

			// Рядки з даними
			rec := records[len(records)-id.Row] // Показуємо новіші зверху
			switch id.Col {
			case 0:
				label.SetText(rec.Date.Format("02-01-2006 15:04"))
			case 1:
				label.SetText(strconv.Itoa(rec.Total))
			case 2:
				details := ""
				for _, denom := range logic.Denominations {
					if count := rec.Breakdown[denom]; count > 0 {
						details += fmt.Sprintf("%dгрнx%d; ", denom, count)
					}
				}
				label.SetText(details)
			}
		},
	)
    table.SetColumnWidth(0, 150)
    table.SetColumnWidth(1, 100)
    table.SetColumnWidth(2, 450)


	win.SetContent(container.NewBorder(totalLabel, nil, nil, nil, table))

	return win
}