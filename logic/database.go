package logic

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // Імпортуємо драйвер
)

var DB *sql.DB

type Record struct {
	ID        int64
	Date      time.Time
	Total     int
	Breakdown map[int]int
}

var Denominations = []int{1000, 500, 200, 100, 50, 20, 10, 5, 2, 1}

func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS records (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		date DATETIME NOT NULL,
		total INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS breakdown (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		record_id INTEGER NOT NULL,
		denomination INTEGER NOT NULL,
		count INTEGER NOT NULL,
		FOREIGN KEY(record_id) REFERENCES records(id) ON DELETE CASCADE
	);
	`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("помилка створення таблиць: %w", err)
	}

	return nil
}

func SaveRecord(record Record) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO records(date, total) VALUES(?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(record.Date, record.Total)
	if err != nil {
		tx.Rollback()
		return err
	}
	stmt.Close()

	recordID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err = tx.Prepare("INSERT INTO breakdown(record_id, denomination, count) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for denom, count := range record.Breakdown {
		if count > 0 {
			if _, err := stmt.Exec(recordID, denom, count); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

// UpdateRecord оновлює існуючий запис
func UpdateRecord(record Record) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// 1. Оновлюємо основний запис
	_, err = tx.Exec("UPDATE records SET date = ?, total = ? WHERE id = ?", record.Date, record.Total, record.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. Видаляємо стару деталізацію
	_, err = tx.Exec("DELETE FROM breakdown WHERE record_id = ?", record.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. Додаємо нову деталізацію
	stmt, err := tx.Prepare("INSERT INTO breakdown(record_id, denomination, count) VALUES(?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for denom, count := range record.Breakdown {
		if count > 0 {
			if _, err := stmt.Exec(record.ID, denom, count); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

// DeleteRecord видаляє запис за його ID
func DeleteRecord(id int64) error {
	// Завдяки "ON DELETE CASCADE" нам достатньо видалити запис з головної таблиці
	_, err := DB.Exec("DELETE FROM records WHERE id = ?", id)
	return err
}

// LoadRecordsByDateRange завантажує записи за певний період
func LoadRecordsByDateRange(start, end time.Time) ([]Record, error) {
	query := "SELECT id, date, total FROM records WHERE date BETWEEN ? AND ? ORDER BY date DESC"
	return loadRecordsWithQuery(query, start, end)
}

// LoadAllRecords завантажує всі записи
func LoadAllRecords() ([]Record, error) {
	query := "SELECT id, date, total FROM records ORDER BY date DESC"
	return loadRecordsWithQuery(query)
}

// Приватна функція, щоб уникнути дублювання коду
func loadRecordsWithQuery(query string, args ...interface{}) ([]Record, error) {
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ID, &rec.Date, &rec.Total); err != nil {
			return nil, err
		}

		rec.Breakdown = make(map[int]int)
		breakdownRows, err := DB.Query("SELECT denomination, count FROM breakdown WHERE record_id = ?", rec.ID)
		if err != nil {
			return nil, err
		}

		for breakdownRows.Next() {
			var denom, count int
			if err := breakdownRows.Scan(&denom, &count); err != nil {
				breakdownRows.Close()
				return nil, err
			}
			rec.Breakdown[denom] = count
		}
		breakdownRows.Close()

		records = append(records, rec)
	}

	return records, nil
}