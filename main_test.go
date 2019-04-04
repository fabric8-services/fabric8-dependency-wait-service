package main

import (
	"net/url"
	"path/filepath"
	"testing"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func Test_isIn(t *testing.T) {
	httpGood := []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}
	if !isIn(httpGood, 200) {
		t.Errorf("Expected 200 to be in httpGood")
	}
	if isIn(httpGood, 20) {
		t.Errorf("Expected 20 to not be in httpGood")
	}
}

func Test_ConfigMysqlFormat(t *testing.T) {
	expected := "root:admin@tcp(localhost:3307)/airlock"
	u, _ := url.Parse("mariadb://root:admin@localhost:3307/airlock")

	cfg := mysql.NewConfig()
	cfg.DBName = filepath.Base(u.Path)
	cfg.User = u.User.Username()
	cfg.Passwd, _ = u.User.Password()
	cfg.Addr = u.Host
	cfg.Net = "tcp"

	found := cfg.FormatDSN()
	if expected != found {
		t.Errorf("expected: %s found: %s", expected, found)
	}
}

func Test_generateConnString(t *testing.T) {
	input, _ := url.Parse("mariadb://root:admin@localhost:3307/airlock")
	want := "root:admin@tcp(localhost:3307)/airlock"
	got := generateConnString(input)
	if got != want {
		t.Errorf("generateConnString() = %v, want %v", got, want)
	}
}
