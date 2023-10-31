package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ipCmd)
}

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "output the local ip",
	Run: func(cmd *cobra.Command, args []string) {
		var ip = getLocalIp()
		fmt.Println(ip)
	},
}

func getLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("can not get local ip")
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
