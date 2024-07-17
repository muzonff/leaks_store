package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strings"
)

type DB_CFG struct {
	Port     int
	User     string
	Password string
	Host     string
	db_name  string
}

func main() {
	var conf DB_CFG
	err := envconfig.Process("store", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	conn_string := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.db_name)
	db, err := sql.Open("mysql", conn_string) //root:root@tcp(localhost:8889)/digger
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Выберите режим работы (1-поиск, 2-добавление): ")
	work_type, _ := reader.ReadString('\n')
	work_type = strings.TrimSpace(work_type)
	if work_type == "1" {
		search(db)
	} else {

		// Спросить пользователя, какой файл он хочет парсить
		fmt.Print("Введите путь к файлу (Excel или CSV): ")
		filePath, _ := reader.ReadString('\n')
		filePath = strings.TrimSpace(filePath)

		// Определить тип файла по расширению
		var fileType string
		if strings.HasSuffix(filePath, ".xlsx") {
			fileType = "excel"
		} else if strings.HasSuffix(filePath, ".csv") {
			fileType = "csv"
		} else {
			log.Fatal("Неподдерживаемый формат файла")
		}

		// Спросить пользователя название таблицы
		fmt.Print("Введите название таблицы: ")
		tableName, _ := reader.ReadString('\n')
		tableName = strings.TrimSpace(tableName)

		var columns []string
		var records [][]string

		// Парсинг файла
		if fileType == "excel" {
			f, err := excelize.OpenFile(filePath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			// Спросить пользователя название листа
			fmt.Print("Введите название листа (по умолчанию 'Sheet1'): ")
			sheetName, _ := reader.ReadString('\n')
			sheetName = strings.TrimSpace(sheetName)
			if sheetName == "" {
				sheetName = "Sheet1"
			}

			// Получение всех строк из указанного листа
			rows, err := f.GetRows(sheetName)
			if err != nil {
				log.Fatal(err)
			}

			if len(rows) > 0 {
				columns = nil
				records = rows[1:]
			}

		} else if fileType == "csv" {
			f, err := os.Open(filePath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			reader := csv.NewReader(f)
			rows, err := reader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}

			if len(rows) > 0 {
				columns = rows[0]
				records = rows[1:]
			}
		}

		// Если заголовки не указаны, спросить у пользователя
		if len(columns) == 0 {
			fmt.Print("Введите названия столбцов через запятую: ")
			columnsInput, _ := reader.ReadString('\n')
			columns = strings.Split(strings.TrimSpace(columnsInput), ",")
		}

		// Создание таблицы
		createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)
		for i, column := range columns {
			createTableSQL += fmt.Sprintf("`%s` TEXT", column)
			if i < len(columns)-1 {
				createTableSQL += ", "
			}
		}
		createTableSQL += ");"
		_, err = db.Exec(createTableSQL)
		if err != nil {
			log.Fatal(err)
		}

		// Вставка данных в таблицу
		for _, record := range records {
			values := make([]interface{}, len(record))
			placeholders := make([]string, len(record))
			for i, value := range record {
				values[i] = value
				placeholders[i] = "?"
			}
			insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ","), strings.Join(placeholders, ","))
			_, err = db.Exec(insertSQL, values...)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Данные успешно добавлены в таблицу", tableName)
	}
}
