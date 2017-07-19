package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

var PollIntervals = []int{
	//1, 2, 2, // polls after 1sec, 2secs, 2secs, then stops
	//1, -2, // polls after 1sec, then infinitely every 2 seconds
	0, -5, // polls immediately, then infinitely every 5 seconds

	// 1, 2, 2, 5, 5, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	// 1, 2, 2,

}

func main() {

	if len(os.Args) == 1 {
		return
	}

	// first check if all urls are valid
	err := isAllProtocolsValid(os.Args[1:])
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// dispatch polling in order received
	for _, v := range os.Args[1:] {
		smallv := strings.ToLower(v)
		if strings.HasPrefix(smallv, "http") {
			isUp, totSecs := pollHTTP200(smallv, PollIntervals)
			if !isUp {
				log.Fatalf("\tNot ok. Service %s polling timedout after %d seconds.\n", smallv, totSecs)
			} else {
				log.Printf("\tOk. Service %s is up after about %d seconds.\n", smallv, totSecs)
			}
		} else if strings.HasPrefix(smallv, "postgres") {
			useDbPing := false
			isUp, totSecs := pollPostgres(smallv, PollIntervals, useDbPing)
			if !isUp {
				log.Fatalf("\tNot ok. Service %s polling timedout after %d seconds.\n", smallv, totSecs)
			} else {
				log.Printf("\tOk. Service %s is up after about %d seconds.\n", smallv, totSecs)
			}
		}
	}
}

func isAllProtocolsValid(args []string) error {
	for _, v := range args {
		smallv := strings.ToLower(v)
		if !strings.HasPrefix(smallv, "http") && !strings.HasPrefix(smallv, "postgres") {
			return fmt.Errorf("Unrecognized protocol for %s. Allowed list is [http, postgres]", v)
		}
	}
	return nil
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

func pollPostgres(url string, pollIntervals []int, useDbPing bool) (bool, int) {

	pg_isready := "pg_isready"

	// check if pg_isready exists
	err := isInPath(pg_isready)
	if err != nil {
		log.Printf("Cannot continue, required command %s not found in path.\n", pg_isready)
		return false, -1
	}

	// split url into host and port
	host, port, err := splitPostgresURL(url)
	if err != nil {
		log.Printf("Invalid postgres url. Should be in the form: postgres://host:port")
		return false, -1
	}

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

		if !useDbPing {
			out, _ := captureOutput(exec.Command(pg_isready, "-h", host, "-p", port))
			log.Print("\tpg_isready response: " + string(out))
			if bytes.Index(out, []byte("accepting connections")) >= 0 {
				isUp = true
				break
			}
		} else {
			if postgresDBPing(url) {
				isUp = true
				break
			}
		}
	}

	if isUp {
		// log.Printf("\t%s is Up\n", url)
		return true, totSecs
	}

	// log.Printf("\t%s timeout\n", url)
	return false, totSecs
}

// url format: fmt.Sprintf("postgres://user:password@host:port/db"
func postgresDBPing(url string) bool {

	db, err := sql.Open("pgx", url)
	if err != nil {
		return false
	}

	err = db.Ping()
	if err != nil {
		return false
	}
	return true
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

func isInPath(cmd string) error {
	_, err := exec.LookPath(cmd)
	return err
}

// returns the host and port of the postgres url
func splitPostgresURL(pgURL string) (string, string, error) {
	u, err := url.Parse(pgURL)
	if err != nil {
		return "", "", err
	}

	if u.Scheme != "postgres" {
		return "", "", fmt.Errorf("Expected a postgres scheme. Received: %s", u.Scheme)
	}

	var host, port string
	h := strings.Split(u.Host, ":")
	if len(h) >= 1 {
		host = h[0]
	}
	if len(h) >= 2 {
		port = h[1]
	}

	return host, port, nil
}
