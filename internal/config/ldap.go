package config

import (
	"strconv"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
)

var Ldap = LdapConfig{}

type LdapConfig struct {
	Host      string
	Port      int
	EnableTls bool
	BindDN    string
	BindPass  string

	BaseDN     string
	UserFilter string

	ready bool
}

func (c *LdapConfig) IsReady() bool {
	return c.ready
}

func (c *LdapConfig) initConfig() {
	c.Host = getEnvOrDefault("LDAP_HOST", "")
	c.Port = getIntEnv("LDAP_PORT", 389)
	c.EnableTls = getBoolEnv("LDAP_ENABLE_TLS", false)
	c.BindDN = getEnvOrDefault("LDAP_BIND_DN", "")
	c.BindPass = getEnvOrDefault("LDAP_BIND_PASS", "")
	c.BaseDN = getEnvOrDefault("LDAP_BASE_DN", "")
	c.UserFilter = getEnvOrDefault("LDAP_USER_FILTER", "")

	go c.InitReady()
}

func (c *LdapConfig) dial() (*ldap.Conn, error) {
	if c.EnableTls {
		return ldap.DialTLS("tcp", c.Host+":"+strconv.Itoa(c.Port), nil)
	}
	return ldap.Dial("tcp", c.Host+":"+strconv.Itoa(c.Port))
}

func (c *LdapConfig) bind(l *ldap.Conn) error {
	if len(c.BindDN) != 0 && len(c.BindPass) != 0 {
		return l.Bind(c.BindDN, c.BindPass)
	}
	return nil
}

func (c *LdapConfig) formatUserFilter(name string) string {
	return strings.ReplaceAll(c.UserFilter, "%s", name)
}

func (c *LdapConfig) FindUserDN(l *ldap.Conn, name string) (string, error) {
	err := c.bind(l)
	if err != nil {
		return "", err
	}
	userFilter := c.formatUserFilter(name)

	search := ldap.NewSearchRequest(c.BaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, userFilter, []string{}, nil)
	sr, err := l.Search(search)
	if err != nil {
		return "", err
	}
	if len(sr.Entries) == 0 {
		return "", nil
	}
	return sr.Entries[0].DN, nil
}

func (c *LdapConfig) InitReady() error {
	if c.ready {
		return nil
	}
	c.ready = false
	l, err := c.dial()
	if err != nil {
		return err
	}
	defer l.Close()
	err = c.bind(l)
	if err != nil {
		return err
	}
	c.ready = true
	return nil
}

func (c *LdapConfig) SearchEntry(name, password string) (string, error) {
	l, err := c.dial()
	if err != nil {
		return "", err
	}
	defer l.Close()
	userDn, err := c.FindUserDN(l, name)
	if err != nil {
		return "", err
	}

	err = l.Bind(userDn, password)
	if err != nil {
		return "", err
	}
	return name, nil
}
