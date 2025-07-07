package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"cash_tracker/v2/logic" // Замініть на ваш шлях
)

const timeFormat = "02-01-2006"

// CreateStatsWindow створює вікно статистики (спрощена версія)
func CreateStatsWindow(app fyne.App, parent fyne.Window) {
	statsWin := app.NewWindow("Статистика надходжень")
	statsWin.Resize(fyne.NewSize(800, 600))

	var records []logic.Record
	var table *widget.Table

	// Функція для оновлення таблиці
	refreshTable := func(newRecords []logic.Record) {
		records = newRecords
		if table != nil {
			table.Refresh()
		}
	}

	// --- Компоненти фільтрації ---
	startDateEntry := widget.NewEntry()
	startDateEntry.SetPlaceHolder(timeFormat)

	endDateEntry := widget.NewEntry()
	endDateEntry.SetPlaceHolder(timeFormat)
	endDateEntry.SetText(time.Now().Format(timeFormat))

	filterButton := widget.NewButton("Фільтрувати", func() {
		start, err1 := time.Parse(timeFormat, startDateEntry.Text)
		end, err2 := time.Parse(timeFormat, endDateEntry.Text)

		if err1 != nil || err2 != nil {
			dialog.ShowError(fmt.Errorf("неправильний формат дати. Використовуйте ДД-ММ-РРРР"), statsWin)
			return
		}
		end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

		loadedRecords, err := logic.LoadRecordsByDateRange(start, end)
		if err != nil {
			dialog.ShowError(err, statsWin)
			return
		}
		refreshTable(loadedRecords)
	})

	resetButton := widget.NewButton("Скинути", func() {
		startDateEntry.SetText("")
		endDateEntry.SetText(time.Now().Format(timeFormat))
		allRecords, _ := logic.LoadAllRecords()
		refreshTable(allRecords)
	})

	filterBox := container.NewHBox(
		widget.NewLabel("З:"), startDateEntry,
		widget.NewLabel("По:"), endDateEntry,
		filterButton, resetButton,
	)

	// --- Створення таблиці ---
	table = widget.NewTable(
		func() (int, int) { return len(records) + 1, 5 }, // 5 стовпців
		// **ЗМІНА 1: Створюємо універсальний контейнер для комірки**
		func() fyne.CanvasObject {
			// Кожна комірка буде містити і текст, і кнопку.
			// container.Max показує лише один видимий віджет зі стопки.
			return container.NewMax(widget.NewLabel("..."), widget.NewButton("", nil))
		},
		// **ЗМІНА 2: Оновлюємо вміст контейнера, а не замінюємо його**
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			// Отримуємо доступ до віджетів всередині контейнера
			c := cell.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			button := c.Objects[1].(*widget.Button)

			if id.Row == 0 { // Заголовки
				button.Hide()
				label.Show()
				headers := []string{"ID", "Дата", "Сума", "Редагувати", "Видалити"}
				label.SetText(headers[id.Col])
				label.TextStyle.Bold = true
				return
			}
			label.TextStyle.Bold = false
			
			rec := records[id.Row-1]

			switch id.Col {
			case 0, 1, 2: // Текстові стовпці
				button.Hide()
				label.Show()
				var text string
				if id.Col == 0 {
					text = strconv.FormatInt(rec.ID, 10)
				} else if id.Col == 1 {
					text = rec.Date.Format("02-01-2006 15:04")
				} else {
					text = fmt.Sprintf("%d грн", rec.Total)
				}
				label.SetText(text)
			case 3: // Кнопка "Редагувати"
				label.Hide()
				button.Show()
				button.SetText("✏️")
				button.OnTapped = func() {
					CreateEditWindow(app, rec, func() {
						resetButton.OnTapped()
					})
				}
			case 4: // Кнопка "Видалити"
				label.Hide()
				button.Show()
				button.SetText("🗑️")
				button.OnTapped = func() {
					dialog.ShowConfirm("Підтвердження", fmt.Sprintf("Ви впевнені, що хочете видалити запис #%d?", rec.ID), func(ok bool) {
						if !ok { return }
						if err := logic.DeleteRecord(rec.ID); err != nil {
							dialog.ShowError(err, statsWin)
							return
						}
						resetButton.OnTapped()
					}, statsWin)
				}
			}
		},
	)
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 160)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 100)

	content := container.NewBorder(filterBox, nil, nil, nil, table)

	statsWin.SetContent(content)
	statsWin.Show()

	resetButton.OnTapped()
}