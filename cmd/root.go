/*
Copyright © 2022 Patrick Falk Nielsen <git@patricknielsen.dk>
*/
package cmd

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/patrickfnielsen/nsoctl/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile   string
	debugFlag bool
	cfg       config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nsoctl",
	Short: "nsoctl helps in day to day operations of Cisco NSO",
	Long: `nsoctl is used in day to day operations of Cisco NSO.
It gives the user a quick way to performe actions on services and devices.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "enabled to show debug output")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nsoctl.toml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".nsoctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".nsoctl")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if debugFlag {
			log.Printf("Using config file: %s", viper.ConfigFileUsed())
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		if debugFlag {
			log.Printf("Couldn't read config: %s", err)
		}
	}

	// if needed disable tls verification
	if cfg.Nso.InsecureSkipVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}
