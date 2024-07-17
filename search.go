package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	_ "github.com/go-sql-driver/mysql"
)

func search(db *sql.DB) {

	reader := bufio.NewReader(os.Stdin)

	// Спросить пользователя значение для поиска
	fmt.Print("Введите значение для поиска: ")
	searchValue, _ := reader.ReadString('\n')
	searchValue = strings.TrimSpace(searchValue)

	// Получение списка всех таблиц в базе данных
	tables, err := getTables(db)
	if err != nil {
		log.Fatal(err)
	}

	// Поиск значения по всем таблицам и столбцам
	for _, table := range tables {
		columns, err := getColumns(db, table)
		if err != nil {
			log.Fatal(err)
		}

		for _, column := range columns {
			rows, err := searchInTable(db, table, column, searchValue)
			if err != nil {
				log.Fatal(err)
			}

			for rows.Next() {
				columns, err := rows.Columns()
				if err != nil {
					log.Fatal(err)
				}

				values := make([]sql.NullString, len(columns))
				scanArgs := make([]interface{}, len(values))
				for i := range values {
					scanArgs[i] = &values[i]
				}

				err = rows.Scan(scanArgs...)
				if err != nil {
					log.Fatal(err)
				}

				result := make(map[string]string)
				for i, col := range values {
					if col.Valid {
						result[columns[i]] = col.String
					} else {
						result[columns[i]] = "NULL"
					}
				}

				fmt.Printf("Найдено в таблице %s, столбец %s:\n", table, column)
				printResult(result)
			}
		}
	}
}

// getTables получает список всех таблиц в базе данных
func getTables(db *sql.DB) ([]string, error) {
	var tables []string
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// getColumns получает список всех столбцов для заданной таблицы
func getColumns(db *sql.DB, table string) ([]string, error) {
	var columns []string
	query := fmt.Sprintf("SHOW COLUMNS FROM `%s`", table)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var field, colType, null, key string
		var defaultValue sql.NullString
		var extra string
		err := rows.Scan(&field, &colType, &null, &key, &defaultValue, &extra)
		if err != nil {
			return nil, err
		}
		columns = append(columns, field)
	}

	return columns, nil
}

// searchInTable выполняет поиск значения в указанной таблице и столбце
func searchInTable(db *sql.DB, table, column, value string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT * FROM `%s` WHERE `%s` LIKE ?", table, column)
	return db.Query(query, "%"+value+"%")
}

func printResult(result map[string]string) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(writer, "СТОЛБЕЦ\tЗНАЧЕНИЕ")
	fmt.Fprintln(writer, "-------\t--------")
	for key, value := range result {
		fmt.Fprintf(writer, "%s\t%s\n", key, value)
	}
	writer.Flush()
}
