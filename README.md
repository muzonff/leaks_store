# leaks_store
Simple Golang app to convert CSV/EXCEL data leaks to one MYSQL base

## Functionality
-  Parse CSV and excel files
-  Asking if column names isn't defined
-  Writing data to MYSQL database instance
-  Searching across all converted recodrs

## Technologies 
- [SQL driver for Golang](github.com/go-sql-driver/mysql)
- [Envconfig](github.com/kelseyhightower/envconfig)
- [Excelize (lib for working with Excel files on Go)](github.com/xuri/excelize)

## Usage
1) Clone repository and install depends
   ```sh
   git clone https://github.com/muzonff/leaks_store/
   cd leaks_store
   go mod tidy
   ```
2) Add Database credentials to env
   ```sh
    export STORE_PORT=8889
    export STORE_USER=root
    export STORE_PASSWORD=root
    export STORE_HOST=localhost
    export STORE_DB_NAME=digger
   ```
3) Just run app
   ```sh
    go run .
   ```
## TODO
- [x] Rewrite README
- [x] Optimize search
- [ ] Add TXT base support
- [ ] Add telegram bot functionality

## Support
If you need help or get an error - welcome to issue's
