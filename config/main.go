package config

import (
	"errors"
	"fmt"

	"github.com/Aykutfgoktas/orc/cfile"
)

type Service interface {
	// ConfigFile returns the configuration file path.
	ConfigFile() string

	// CheckConfigFile checks if the configuration file is exists.
	CheckConfigFile() bool

	// Create creates the configuration file.
	Create(apikey, org string) (string, error)

	// Read reads the configration file content.
	Read() (*Config, error)

	// UpdateDefaultOrganization updates default organization on the configuration.
	UpdateDefaultOrganization(org string) error

	// AddOrganization adds the new organization to the configuration.
	AddOrganization(org string) (bool, error)

	// DeleteOrganization deletes the selected organization.
	DeleteOrganization(org string) error
}

type Organizations []string

func (a *Organizations) Add(val string) {
	*a = append(*a, val)
}

func (a *Organizations) Remove(val string) {
	for i, v := range *a {
		if v == val {
			*a = append((*a)[:i], (*a)[i+1:]...)
			break
		}
	}
}

func (a *Organizations) Exists(val string) bool {
	for _, v := range *a {
		if v == val {
			return true
		}
	}
	return false
}

type Config struct {
	APIKey              string        `json:"key"`
	DefaultOrganization string        `json:"org"`
	Organizations       Organizations `json:"orgs"`
}

type config struct {
	cfile cfile.IConfigFile
}

func New(cfile cfile.IConfigFile) Service {
	return &config{
		cfile: cfile,
	}
}

func (c *config) ConfigFile() string {
	return c.cfile.ConfigFile()
}

func (c *config) CheckConfigFile() bool {
	return c.cfile.CheckConfigFile()
}

func (c *config) Create(apikey, org string) (string, error) {
	conf := Config{
		APIKey:              apikey,
		DefaultOrganization: org,
		Organizations:       []string{org},
	}

	path, err := c.cfile.Writer(conf)

	if err != nil {
		return "", writerError(err)
	}

	return path, nil
}

func (c *config) Read() (*Config, error) {
	result, err := c.cfile.Reader()

	if err != nil {
		return nil, err
	}

	conf := Config{}

	if err = result.Decode(&conf); err != nil {
		return nil, decodeError(err)
	}

	return &conf, nil
}

func (c *config) UpdateDefaultOrganization(org string) error {
	result, err := c.cfile.Reader()

	if err != nil {
		return readerError(err)
	}

	conf := Config{}

	err = result.Decode(&conf)

	if err != nil {
		return decodeError(err)
	}

	conf.DefaultOrganization = org

	if _, err = c.cfile.Writer(conf); err != nil {
		return writerError(err)
	}

	return nil
}

func (c *config) AddOrganization(org string) (bool, error) {
	flag := false

	result, err := c.cfile.Reader()

	if err != nil {
		return flag, readerError(err)
	}

	conf := Config{}

	if err = result.Decode(&conf); err != nil {
		return flag, decodeError(err)
	}

	flag = conf.Organizations.Exists(org)

	if !flag {
		conf.Organizations.Add(org)

		if _, err := c.cfile.Writer(conf); err != nil {
			return false, writerError(err)
		}
	}

	return flag, nil
}

func (c *config) DeleteOrganization(org string) error {
	result, err := c.cfile.Reader()

	if err != nil {
		return readerError(err)
	}

	conf := Config{}

	if err = result.Decode(&conf); err != nil {
		return decodeError(err)
	}

	conf.Organizations.Remove(org)

	if _, err := c.cfile.Writer(conf); err != nil {
		return writerError(err)
	}

	return nil
}

func writerError(e error) error {
	m := fmt.Sprintf("Error while creating the config file error: %s", e.Error())
	return errors.New(m)
}

func readerError(e error) error {
	m := fmt.Sprintf("Error while reading the config file error: %s", e.Error())
	return errors.New(m)
}

func decodeError(e error) error {
	m := fmt.Sprintf("Error while decoding the config file error: %s", e.Error())
	return errors.New(m)
}
