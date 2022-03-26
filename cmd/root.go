/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	isDebug    bool
	isHeadless bool
	remoteWs   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bli",
	Short: "bli is a bilibili crawlers command line tool",
	Long:  `bli is a bilibili crawlers command line tool`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "Output verbose debug information")
	rootCmd.PersistentFlags().BoolVar(&isHeadless, "headless", false, "set headless mode, default value is true that means no gui")
	rootCmd.PersistentFlags().StringVar(&remoteWs, "remote-ws", "", "remote websocket address from docker container")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// set log
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		DisableQuote:  true,
		FullTimestamp: true,
		// FieldMap: log.FieldMap{
		// 	log.FieldKeyTime:  "@time",
		// 	log.FieldKeyLevel: "@level",
		// 	log.FieldKeyMsg:   "@msg",
		// 	log.FieldKeyFunc:  "@caller",
		// },
	})
	log.SetOutput(os.Stdout)
	switch isDebug {
	case true:
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		log.Debugln("Enabling debug output")
	default:
		log.SetLevel(log.InfoLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".bli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	log.WithFields(log.Fields{
		"event": "start",
		"key":   "root",
	}).Debugf("Start it!\n")
}
