package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
)

var defaultPollIntervals = []int{
	//1, 2, 2, // polls after 1sec, 2secs, 2secs, then stops
	//1, -2, // polls after 1sec, then infinitely every 2 seconds
	//0, -5, // polls immediately, then infinitely every 5 seconds
	0, -1, // polls immediately, then infinitely every 1 seconds

	// 1, 2, 2, 5, 5, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	// 1, 2, 2,

}

var gVerbose bool

func main() {

	gVerbose = getVerbosity()

	if len(os.Args) == 1 {
		return
	}

	pollIntervals := getPollIntervals()

	// dispatch polling in order received
	v := os.Args[1]
	smallv := strings.ToLower(v)

	if strings.HasPrefix(smallv, "http") {
		isUp, totSecs := pollHTTP200(smallv, pollIntervals)
		if !isUp {
			log.Fatalf("\tNot ok. Service %s polling timedout after %d seconds.\n", smallv, totSecs)
		} else {
			log.Printf("\tOk. Service %s is up after about %d seconds.\n", smallv, totSecs)
		}
	} else {
		isUp, totSecs := pollDB(smallv, pollIntervals)
		if !isUp {
			log.Fatalf("\tNot ok. Service %s polling timedout after %d seconds.\n", smallv, totSecs)
		} else {
			log.Printf("\tOk. Service %s is up after about %d seconds.\n", smallv, totSecs)
		}
	}
}

func pollHTTP200(url string, pollIntervals []int) (bool, int) {
	// intervals in seconds.

	log.Printf("Checking if %s is up.\n", url)
	isUp := false
	totSecs := 0
	for i := 0; i < len(pollIntervals); i++ {
		pollInt := pollIntervals[i]
		if pollInt < 0 {
			// infinitely keep adding the same entry
			pollIntervals = append(pollIntervals, pollInt)
			pollInt = -pollInt
		}
		log.Printf("\tNext poll after %d seconds.\n", pollInt)
		totSecs += pollInt
		time.Sleep(time.Second * time.Duration(pollInt))

		if gVerbose {
			log.Printf("\tGoing to check %s\n", url)
		}
		ok := httpPoll(url)
		if ok {
			isUp = true
			break
		}
	}

	if isUp {
		// log.Printf("\t%s is Up\n", url)
		return true, totSecs
	}

	// log.Printf("\t%s timeout\n", url)
	return false, totSecs
}

func httpPoll(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}

	//return isIn([]int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}, resp.StatusCode)
	return isIn([]int{200}, resp.StatusCode)
}

func getConnArgs(dbURL string) (string, string, error) {
	u, err := url.Parse(dbURL)

	if err != nil {
		return "", "", err
	}

	if u.Scheme == "postgres" {
		return "pgx", dbURL, nil
	} else if u.Scheme == "mariadb" || u.Scheme == "mysql" {
		return "mysql", generateConnString(u), nil
	} else {
		return "", "", errors.New("unsupported db protocol")
	}
}

func generateConnString(u *url.URL) string {
	cfg := mysql.NewConfig()
	if u.Path != "" {
		cfg.DBName = filepath.Base(u.Path)
	}
	cfg.User = u.User.Username()
	cfg.Passwd, _ = u.User.Password()
	cfg.Addr = u.Host
	cfg.Net = "tcp"
	return cfg.FormatDSN()
}

func pollDB(dbURL string, pollIntervals []int) (bool, int) {
	log.Printf("Checking if %s is up.\n", dbURL)
	totSecs := 0

	driverName, dataSourceName, err := getConnArgs(dbURL)
	if err != nil {
		log.Fatal(err)
		return false, totSecs
	}
	db, err := sql.Open(driverName, dataSourceName)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err != nil {
		log.Fatal(err)
		defer db.Close()
		return false, totSecs
	}

	for i := 0; i < len(pollIntervals); i++ {
		pollInt := pollIntervals[i]

		if pollInt < 0 {
			// infinitely keep adding the same entry
			pollIntervals = append(pollIntervals, pollInt)
			pollInt = -pollInt
		}
		totSecs += pollInt
		time.Sleep(time.Duration(pollInt) * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollInt)*time.Second)
		defer cancel()
		log.Printf("\tchecking \n")
		err = db.PingContext(ctx)
		if err == nil {
			defer db.Close()
			return true, totSecs
		}
		log.Printf("\tNext poll after %d seconds.\n", pollInt)
	}
	// log.Printf("\t%s timeout\n", url)
	return false, totSecs
}

func captureOutput(cmd *exec.Cmd) ([]byte, error) {
	// out, err := cmd.Output()
	return cmd.Output()
}

func isIn(list []int, val int) bool {
	for _, v := range list {
		if val == v {
			return true
		}
	}
	return false
}

func getPollIntervals() []int {
	key := "DEPENDENCY_POLL_INTERVAL"
	intervalStr := strings.TrimSpace(os.Getenv(key))
	if len(intervalStr) == 0 {
		return defaultPollIntervals
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		log.Printf("Error: Dependency service key %s has invalid value: %s. Should be a positive integer.\n", key, intervalStr)
		os.Exit(1)
	}

	if gVerbose {
		log.Printf("Got env value for poll interval. Will poll every %d seconds.\n", interval)
	}

	return []int{0, -interval}
}

func getVerbosity() bool {
	key := "DEPENDENCY_LOG_VERBOSE"
	verbosityStr := strings.TrimSpace(os.Getenv(key))
	if len(verbosityStr) == 0 {
		return false
	}

	verbosity, err := strconv.ParseBool(verbosityStr)
	if err != nil {
		log.Printf("Error: Dependency service key %s has invalid value: %s. Should be true/false.\n", key, verbosityStr)
		os.Exit(1)
	}

	if verbosity {
		log.Printf("Got env value verbose: %t.\n", verbosity)
	}

	return verbosity
}
