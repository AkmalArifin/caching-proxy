/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AkmalArifin/caching-proxy/internal/proxy"
	"github.com/spf13/cobra"
)

// Flags
var port int64
var origin string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "caching-proxy",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	proxy := proxy.NewProxy(origin)

	http.Handle("/", proxy)

	log.Printf("Listen to port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.caching-proxy.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().Int64VarP(&port, "port", "p", 8080, "the port on which the caching proxy server will run")
	rootCmd.PersistentFlags().StringVarP(&origin, "origin", "o", "", "the URL of the server to which the requests will be forwarded")

	rootCmd.MarkPersistentFlagRequired("origin")
}
