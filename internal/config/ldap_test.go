package config

import (
	"testing"
)

func TestDial(t *testing.T) {
	InitConfig()
	conf := Ldap

	con, err := conf.dial()
	if err != nil {
		t.Fatalf("Failed to dial LDAP server: %v", err)
	}
	defer con.Close()

	err = conf.bind(con)
	if err != nil {
		t.Fatalf("Failed to bind to LDAP server: %v", err)
	}

	userDN, err := conf.FindUserDN(con, "test")
	if err != nil {
		t.Fatalf("Failed to find user DN: %v", err)
	}
	t.Logf("Found user DN: %s", userDN)
}

func TestSearchEntry(t *testing.T) {
	InitConfig()
	conf := Ldap
	name, err := conf.SearchEntry("test", "test123")
	if err != nil {
		t.Fatalf("Failed to find user entry: %v", err)
	}
	t.Logf("Found user name: %s", name)
}

func TestIsReady(t *testing.T) {
	InitConfig()
	conf := Ldap
	conf.InitReady()
	if !conf.IsReady() {
		t.Fatalf("LDAP server is not ready")
	}
}
