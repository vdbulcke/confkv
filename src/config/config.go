package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/vdbulcke/confkv/src/assert"
)

const (
	SourceGit = "git"
	SourceDir = "dir"
)

type Config struct {
	DB      BoltDB            `toml:"db" validate:"required"`
	Buckets map[string]Bucket `toml:"bucket" validate:"required,dive"`

	Cron Scheduler `toml:"scheduler" validate:"required"`
}

type Scheduler struct {
	CronTab string `toml:"cron_tab" validate:"required" default:"0 * * * *"`
}

type BoltDB struct {
	Dir  string `toml:"dir" validate:"required"`
	Name string `toml:"filename"`
}

type Bucket struct {
	// Name string
	Source string `toml:"source" validate:"required,oneof=dir git"`

	LocalDir string `toml:"local_dir" validate:"required"`
	// source Git

	GitRepoUrl        string `toml:"git_repo_url" `
	GitRelativeDir    string `toml:"git_relative_dir"  validate:"required" default:"."`
	GitUsername       string `toml:"git_username" `
	GitPasswordEnvVar string `toml:"git_password_env_var" default:"CONFKV_GIT_PASSWORD"`
}

func MustOpen(filename string) *Config {
	f, err := os.Open(filename)
	assert.ErrNotNil(err, assert.Panic, "error opening file", filename)

	defer f.Close()

	var config Config
	_, err = toml.NewDecoder(f).Decode(&config)
	assert.ErrNotNil(err, assert.Panic, "error parsing toml file", filename)

	err = defaults.Set(&config)
	assert.ErrNotNil(err, assert.Panic, "error setting default")

	return &config
}

func (config *Config) Validate() bool {

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("toml"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	errs := validate.Struct(config)
	if errs != nil {
		for _, e := range errs.(validator.ValidationErrors) {
			fmt.Println(e)
		}
		return false
	}

	valid := true
	for k, v := range config.Buckets {
		if v.Source == SourceGit {
			if v.GitRepoUrl == "" {
				valid = false
				fmt.Printf("Error: missing required 'git_repo_url' for bucket '%s'\n", k)
			}
			if v.GitUsername == "" {
				valid = false
				fmt.Printf("Error: missing required 'git_username' for bucket '%s'\n", k)
			}

		}
	}

	return valid
}

func (config *BoltDB) Validate() bool {

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("toml"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	errs := validate.Struct(config)
	if errs != nil {
		for _, e := range errs.(validator.ValidationErrors) {
			fmt.Println(e)
		}
		return false
	}

	return true
}
