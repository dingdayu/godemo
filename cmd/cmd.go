package cmd

import (
	"demo/pkg/log"
	"fmt"
	"os"

	"demo/cmd/api"
	"demo/pkg/config"

	"github.com/spf13/cobra"
)

// RootCmd RootCmd
var RootCmd = &cobra.Command{
	Use:              "demo",
	Short:            "golang api demo",
	Long:             "demo service demo.",
	TraverseChildren: true,
}

func init() {
	cobra.OnInitialize(onInitialize)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&config.CfgFile, "c", "config/config.yaml", "config file (default is config/config.yaml)")

	RootCmd.AddCommand(api.ServerCmd)

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func onInitialize() {
	//
	config.Init()
	log.Init()
}
