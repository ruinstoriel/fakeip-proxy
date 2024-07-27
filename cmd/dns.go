/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	d "fakeip-proxy/dns"
	"fakeip-proxy/listener"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logfile *os.File

// dnsCmd represents the dns command
var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dns called", viper.Get("log"))
		go d.Start()
		go listener.Start()
		end := make(chan bool, 1)
		wait(func() {
			err := logfile.Close()
			if err != nil {
				log.Fatalf("Failed to close log file %s: %v", *logfile, err)
			}
		}, end)
	},
}

func wait(cancel func(), end chan bool) {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	defer signal.Stop(sigs)
	<-sigs
	fmt.Println("收到退出信号")
	cancel()
	<-end
	fmt.Println("清理结束")

}

func init() {
	rootCmd.AddCommand(dnsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	logFileString := dnsCmd.PersistentFlags().String("log", "/var/log/me", "log file")
	var err error
	logfile, err = os.OpenFile(*logFileString, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", *logfile, err)
	}
	log.SetOutput(logfile)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dnsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
