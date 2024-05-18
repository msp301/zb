package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/msp301/zb/editor"
	"github.com/msp301/zb/graph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

var cfgFile string
var cfgEditor string

var rootCmd = &cobra.Command{
	Use:   "zb",
	Short: "A Zettelkasten notebook utility",
	Run: func(cmd *cobra.Command, args []string) {
		book().Notes.Walk(func(vertex graph.Vertex, depth int) bool {
			jsonStr, err := json.Marshal(vertex)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonStr))
			return true
		}, -1)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zb.toml)")
	rootCmd.PersistentFlags().StringVar(&cfgEditor, "editor", editor.FindEditor(), "program to open notes with")
	rootCmd.PersistentFlags().StringSlice("directory", []string{defaultNotebookDir()}, "notebook directories")
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
	viper.BindPFlag("editor", rootCmd.PersistentFlags().Lookup("editor"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)

		viper.SetConfigType("toml")
		viper.SetConfigName(".zb")
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
