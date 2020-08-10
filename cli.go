package main

import (
	"errors"
	"github.com/rfizzle/collector-helpers/outputs"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func setupCliFlags() error {
	viper.SetEnvPrefix("GSC")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	flag.String("state-path", "google.state", "state file path")
	flag.Int("schedule", 30, "time in seconds to collect")
	flag.String("google-credentials", "", "google service account creds file path")
	flag.String("impersonated-user", "", "user to impersonate for API access")
	flag.BoolP("verbose", "v", false, "verbose logging")
	outputs.InitCLIParams()
	flag.Parse()
	err := viper.BindPFlags(flag.CommandLine)

	if err != nil {
		log.Fatalf("Failed parsing flags: %v", err.Error())
	}

	// Check parameters
	if err := checkRequiredParams(); err != nil {
		return err
	}

	return nil
}

func checkRequiredParams() error {
	if viper.GetString("state-path") == "" {
		return errors.New("missing State File Path param (--state-path)")
	}

	if viper.GetString("google-credentials") == "" {
		return errors.New("missing Google Credentials param (--google-credentials)")
	}

	if viper.GetString("impersonated-user") == "" {
		return errors.New("missing Impersonate User param (--impersonated-user)")
	}

	if err := outputs.ValidateCLIParams(); err != nil {
		return err
	}

	return nil
}
