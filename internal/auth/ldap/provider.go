package ldap

import (
	"adb-backup/internal/config"
)

var Provider = provider{
	config: &config.Ldap,
}

type provider struct {
	config *config.LdapConfig
}

func (p *provider) IsReady() bool {
	p.config.InitReady()
	return p.config.IsReady()
}

func (p *provider) Authenticate(username, password string) (bool, error) {
	_, err := p.config.SearchEntry(username, password)
	if err != nil {
		return false, err
	}
	return true, nil
}
