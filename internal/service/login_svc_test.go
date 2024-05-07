package service

import (
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"testing"
	"time"
)

func TestGenLoginToken(t *testing.T) {
	initPFromLocalAppKeyFile()
	token, err := GenLoginToken("admin")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
}

func TestAnalysisLoginToken(t *testing.T) {
	cmd.InitParam()
	initPFromLocalAppKeyFile()
	serverconf.Initialize()
	log.Initialize()
	token := "ad7f4bb4c0f53fe34b5900f26ca1bcffe670fa14cb70ee96f03ca4beaf5a1d188e0d4bfca81864b14a0ee9e206d306b2bfa00b7f99c1bfe35ee0b23f8038f9cf2b20b9cfc102c6e4b1b77b6130e537cc930937ed5659ef6378bd6122c84c11f5124081e6e49fb8325730a1d7e2a960c58b3d5b1c08f0e78f0f7b8229ef275b1db053feb2fe368b2041a1f6402fb6befb2c0eb9b3fe33d2b6110783a618dd8a11"
	validTime := 12 * time.Hour
	m, err := AnalysisLoginToken(token, validTime.Milliseconds())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", m)
}
