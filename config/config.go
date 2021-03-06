// Copyright 2016 IBM Corporation
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package config

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/amalgam8/controller/util"
	"github.com/codegangsta/cli"
)

// Database config
type Database struct {
	Type     string
	Username string
	Password string
	Host     string
	//URL string
}

// Config for the controller
type Config struct {
	Database     Database
	APIPort      int
	ControlToken string
	SecretKey    string
	StatsdHost   string
	PollInterval time.Duration
	LogLevel     logrus.Level
}

// New config instance
func New(context *cli.Context) *Config {
	// TODO: add HTTPS to database host if not present?

	// TODO: parse this more gracefully
	loggingLevel := logrus.DebugLevel
	logLevelArg := context.String(logLevel)
	var err error
	loggingLevel, err = logrus.ParseLevel(logLevelArg)
	if err != nil {
		loggingLevel = logrus.DebugLevel
	}

	return &Config{
		Database: Database{
			Type:     context.String(dbType),
			Username: context.String(dbUser),
			Password: context.String(dbPassword),
			Host:     "https://" + context.String(dbHost), // FIXME: conditionally add HTTPS
		},
		APIPort:      context.Int(apiPort),
		ControlToken: context.String(controlToken),
		SecretKey:    context.String(secretKey),
		StatsdHost:   context.String(statsdHost),
		PollInterval: context.Duration(pollInterval),
		LogLevel:     loggingLevel,
	}
}

// Validate the config
func (c *Config) Validate() error {
	// Create list of validation checks
	validators := []util.ValidatorFunc{
		util.IsNotEmpty("Database type", c.Database.Type),
		util.IsInRange("API port", c.APIPort, 1, 65535),
		util.IsNotEmpty("Control token", c.ControlToken),
		util.IsNotEmpty("Secret key", c.SecretKey),
		util.IsNotEmpty("Statsd host", c.StatsdHost),
	}

	if c.Database.Type == "cloudant" {
		additionalValidators := []util.ValidatorFunc{
			util.IsNotEmpty("Database password", c.Database.Password),
			util.IsNotEmpty("Database username", c.Database.Username),
			util.IsNotEmpty("Database host name", c.Database.Host),
		}
		validators = append(validators, additionalValidators...)
	} else if c.Database.Type != "memory" {
		return fmt.Errorf("Invalid database type %v", c.Database.Type)
	}

	if len(c.SecretKey) != 16 {
		return fmt.Errorf("Secret must have a length of 16 characters")
	}

	return util.Validate(validators)
}
