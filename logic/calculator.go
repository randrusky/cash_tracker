package logic

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// Номінали українських гривень
var Denominations = []int{1000, 500, 200, 100, 50, 20, 10, 5, 2, 1}

// Record представляє один запис про надходження готівки
type Record struct {
	Date      time.Time `json:"date"`
	Total     int       `json:"total"`
	Breakdown map[int]int `json:"breakdown"` // Кількість купюр по номіналах
}

// SaveRecord додає новий запис у файл JSON
func SaveRecord(record Record) error {
	records, err := LoadRecords()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	records = append(records, record)

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("data/records.json", data, 0644)
}

// LoadRecords завантажує всі записи з файлу JSON
func LoadRecords() ([]Record, error) {
	file, err := ioutil.ReadFile("data/records.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []Record{}, nil // Якщо файлу немає, повертаємо порожній список
		}
		return nil, err
	}

	var records []Record
	err = json.Unmarshal(file, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}