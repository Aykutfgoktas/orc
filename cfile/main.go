package cfile

import (
	"encoding/json"
	"io/fs"
	"os"
)

var permission fs.FileMode = 0600

type IConfigFile interface {
	// ConfigFile returns the configuration file path.
	ConfigFile() string

	// CheckConfigFile checks if the configuration file is exists.
	CheckConfigFile() bool

	// Reader reads the configuration file.
	Reader() (IReader, error)

	// Writer writes the given data to the configuration file.
	Writer(data interface{}) (string, error)
}

type IReader interface {
	// Decode decodes the configuration file content to given interface.
	Decode(d interface{}) error
}

type ReaderResult struct {
	b []byte
}

type cfile struct {
	file string
}

func New(file string) IConfigFile {
	return &cfile{
		file: file,
	}
}

func (r *cfile) ConfigFile() string {
	return r.file
}

func (r *cfile) CheckConfigFile() bool {
	if _, err := os.Stat(r.file); os.IsNotExist(err) {
		return false
	}

	return true
}

func (r *cfile) Reader() (IReader, error) {
	data, err := os.ReadFile(r.file)
	if err != nil {
		return nil, err
	}

	return &ReaderResult{
		b: data,
	}, nil
}

func (r *cfile) Writer(data interface{}) (string, error) {
	b, err := json.Marshal(&data)

	if err != nil {
		return "", err
	}

	err = os.WriteFile(r.file, b, permission)

	if err != nil {
		return "", err
	}

	cur, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return cur, nil
}

func (rr *ReaderResult) Decode(d interface{}) error {
	err := json.Unmarshal(rr.b, d)

	if err != nil {
		return err
	}

	return nil
}
