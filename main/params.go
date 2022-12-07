package main

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	versionKey = "version"
)

func flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("coreth", flag.ContinueOnError)

	fs.Bool(versionKey, false, "If true, print version and quit")

	return fs
}

func getViper() (*viper.Viper, error) {
	v := viper.New()

	fs := flagSet()
	pflag.CommandLine.AddGoFlagSet(fs)
	pflag.Parse()

	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		return nil, err
	}

	return v, nil
}

func PrintVersion() (bool, error) {
	v, err := getViper()
	if err != nil {
		return false, err
	}

	if v.GetBool(versionKey) {
		return true, nil
	}

	return false, nil
}
