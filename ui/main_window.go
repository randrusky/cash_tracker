package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"cash_tracker/v2/logic" // <-- Не забудьте замінити на ваш шлях
)

// CreateMainWindow створює та повертає головне вікно
func CreateMainWindow(app fyne.App) fyne.Window {
	win := app.NewWindow("Облік Готівки")
	win.Resize(fyne.NewSize(450, 600))

	entries := make(map[int]*widget.Entry)
	subtotalLabels := make(map[int]*widget.Label) // Мапа для зберігання лейблів проміжних підсумків
	grandTotalLabel := widget.NewLabel("Загальна сума: 0 грн")
	grandTotalLabel.TextStyle.Bold = true

	// Контейнер для всіх полів вводу
	rowsContainer := container.NewVBox()

	// Функція для оновлення всіх сум
	updateTotals := func() {
		grandTotal := 0
		for _, denom := range logic.Denominations {
			count, err := strconv.Atoi(entries[denom].Text)
			if err != nil {
				count = 0 // Якщо ввід некоректний, вважаємо його за нуль
			}

			subtotal := count * denom
			grandTotal += subtotal

			// Оновлюємо лейбл проміжного підсумку для конкретної купюри
			subtotalLabels[denom].SetText(fmt.Sprintf("= %d грн", subtotal))
		}
		// Оновлюємо загальну суму
		grandTotalLabel.SetText(fmt.Sprintf("Загальна сума: %d грн", grandTotal))
	}

	// Створення рядків для кожного номіналу
	for _, denom := range logic.Denominations {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("0") // Підказка
		entries[denom] = entry

		subtotalLabel := widget.NewLabel("= 0 грн")
		subtotalLabels[denom] = subtotalLabel

		// Встановлюємо функцію, яка буде викликатись при зміні тексту
		entry.OnChanged = func(s string) { updateTotals() }

		// Створюємо горизонтальний контейнер для одного рядка
		row := container.NewHBox(
			widget.NewLabel(fmt.Sprintf("%d грн:", denom)),
			entry,
			subtotalLabel,
		)
		rowsContainer.Add(row)
	}

	// Кнопка для збереження запису
	saveButton := widget.NewButton("Додати запис", func() {
		breakdown := make(map[int]int)
		total := 0
		hasInput := false

		for _, denom := range logic.Denominations {
			count, err := strconv.Atoi(entries[denom].Text)
			if err != nil {
				count = 0
			}
			if count > 0 {
				hasInput = true
			}
			breakdown[denom] = count
			total += count * denom
		}

		// Зберігаємо тільки якщо є якісь дані
		if hasInput {
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
				entry.SetText("") // Це автоматично викличе updateTotals()
			}
		}
	})

	// Кнопка для відкриття статистики
	statsButton := widget.NewButton("Переглянути статистику", func() {
    CreateStatsWindow(app, win) // Новий код, передаємо win
})

	// Фінальне компонування вікна
	win.SetContent(container.NewVBox(
		rowsContainer,
		layout.NewSpacer(), // Розпірка, щоб кнопки були внизу
		grandTotalLabel,
		saveButton,
		statsButton,
	))

	return win
}