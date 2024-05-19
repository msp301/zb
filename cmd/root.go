package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/msp301/zb/config"
	"github.com/msp301/zb/editor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfgEditor string

var rootCmd = &cobra.Command{
	Use:   "zb",
	Short: "A Zettelkasten notebook utility",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", config.GlobalConfigFile))
	rootCmd.PersistentFlags().StringVar(&cfgEditor, "editor", editor.FindEditor(), "program to open notes with")
	rootCmd.PersistentFlags().StringSlice("directory", []string{}, "notebook directories")
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
	viper.BindPFlag("editor", rootCmd.PersistentFlags().Lookup("editor"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(config.GlobalConfigDir)

		viper.SetConfigType(config.CONFIG_TYPE)
		viper.SetConfigName(config.CONFIG_NAME)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// TODO: Use a verbose flag to show this output when wanted
			//fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		} else {
			log.Fatalf("Failed to read config file: %s", err)
		}
	}
}

func defaultNotebookDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic("Failed to get cwd")
	}
	notesDir := path.Join(cwd, "notes")
	if _, err := os.Stat(notesDir); os.IsNotExist(err) {
		notesDir = cwd
	}
	return notesDir
}
