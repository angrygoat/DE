package main

import (
	"configurate"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"logcabin"
	"messaging"
	"net"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

var (
	logger     = logcabin.New()
	version    = flag.Bool("version", false, "Print the version information")
	cfgPath    = flag.String("config", "", "The path to the config file")
	dbURI      = flag.String("db", "", "The URI used to connect to the database")
	amqpURI    = flag.String("amqp", "", "The URI used to connect to the amqp broker")
	gitref     string
	appver     string
	builtby    string
	amqpClient *messaging.Client
	db         *sql.DB
)

func init() {
	flag.Parse()
}

// AppVersion prints version information to stdout
func AppVersion() {
	if appver != "" {
		fmt.Printf("App-Version: %s\n", appver)
	}
	if gitref != "" {
		fmt.Printf("Git-Ref: %s\n", gitref)
	}

	if builtby != "" {
		fmt.Printf("Built-By: %s\n", builtby)
	}
}

func insert(state, invID, msg, host, ip string, sentOn int64) (sql.Result, error) {
	insertStr := `
		INSERT INTO job_status_updates (
			external_id,
			message,
			status,
			sent_from,
			sent_from_hostname,
			sent_on
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		) RETURNING id`
	return db.Exec(insertStr, invID, msg, state, ip, host, sentOn)
}

func msg(delivery amqp.Delivery) {
	delivery.Ack(false)
	logger.Println("Message received")
	update := &messaging.UpdateMessage{}
	err := json.Unmarshal(delivery.Body, update)
	if err != nil {
		logger.Print(err)
		return
	}
	if update.State == "" {
		logger.Println("State was unset, dropping update")
		return
	}
	logger.Printf("State is %s\n", update.State)
	if update.Job.InvocationID == "" {
		logger.Println("InvocationID was unset, dropping update")
	}
	logger.Printf("InvocationID is %s\n", update.Job.InvocationID)
	if update.Message == "" {
		logger.Println("Message set to empty string, setting to UNKNOWN")
		update.Message = "UNKNOWN"
	}
	logger.Printf("Message is: %s", update.Message)
	var sentFromAddr string
	if update.Sender == "" {
		logger.Println("Unknown sender, setting from address to 0.0.0.0")
		update.Sender = "0.0.0.0"
	}
	parsedIP := net.ParseIP(update.Sender)
	if parsedIP != nil {
		sentFromAddr = update.Sender
	} else {
		ips, err := net.LookupIP(update.Sender)
		if err != nil {
			logger.Print(err)
		} else {
			if len(ips) > 0 {
				sentFromAddr = ips[0].String()
			}
		}
	}
	logger.Printf("Sent from: %s", sentFromAddr)
	sentOn, err := strconv.ParseInt(update.SentOn, 10, 64)
	if err != nil {
		logger.Printf("Error parsing SentOn field, setting field to 0: %s", err)
		sentOn = 0
	}
	result, err := insert(
		string(update.State),
		update.Job.InvocationID,
		update.Message,
		update.Sender,
		sentFromAddr,
		sentOn,
	)
	if err != nil {
		logger.Print(err)
		return
	}
	rowCount, err := result.RowsAffected()
	if err != nil {
		logger.Print(err)
		return
	}
	logger.Printf("Inserted %d rows\n", rowCount)
}

func main() {
	var err error

	if *version {
		AppVersion()
		os.Exit(0)
	}

	if *dbURI == "" || *amqpURI == "" {
		if *cfgPath == "" {
			logger.Fatal("--config must be set.")
		}
		err := configurate.Init(*cfgPath)
		if err != nil {
			logger.Fatal(err)
		}
		if *dbURI == "" {
			*dbURI, err = configurate.C.String("db.uri")
			if err != nil {
				logger.Fatal(err)
			}
		}
		if *amqpURI == "" {
			*amqpURI, err = configurate.C.String("amqp.uri")
			if err != nil {
				logger.Fatal(err)
			}
		}
	}
	amqpClient, err := messaging.NewClient(*amqpURI, false)
	if err != nil {
		logger.Fatal(err)
	}
	defer amqpClient.Close()
	logger.Println("Connecting to the database...")
	db, err = sql.Open("postgres", *dbURI)
	if err != nil {
		logger.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Connected to the database")
	amqpClient.AddConsumer(messaging.JobsExchange, "job_status_recorder", messaging.UpdatesKey, msg)
	amqpClient.Listen()
}