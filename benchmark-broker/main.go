package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// dbUser := os.Getenv("DB_USER")
    // dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := "root"
    dbPassword := "1234"
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    // dbName := os.Getenv("DB_NAME")
	dbName := "Application"
    kafkaServers := os.Getenv("KAFKA_SERVERS")
	tableNumber := os.Getenv("TABLE_NUMBER")

	tableName := fmt.Sprintf("throughput_%s", tableNumber)
    columnName := fmt.Sprintf("throughput%s", tableNumber)

    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    _, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        serial_num INT AUTO_INCREMENT PRIMARY KEY,
        %s FLOAT,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )`, tableName, columnName))
	if err != nil {
		log.Fatal(err)
	}

    cmd := exec.Command("sh", "-c", fmt.Sprintf("bin/kafka-producer-perf-test.sh --producer.config config/producer.properties --print-metrics --throughput -1 --num-records 1000000 --topic test --record-size 100000 --producer-props bootstrap.servers=%s", kafkaServers))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	re := regexp.MustCompile(`(\d+\.\d+) MB/sec`)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) 
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			throughput, err := strconv.ParseFloat(matches[1], 64)
			if err != nil {
				log.Fatal(err)
			}

			_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (%s, timestamp) VALUES (?, NOW())", tableName, columnName), throughput)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
