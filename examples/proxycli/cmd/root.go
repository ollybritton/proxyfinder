// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"gitlab.com/ollybritton/proxyfinder"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "proxycli",
	Short: "Find and get proxies on the command line",
	Long:  `proxycli is a tool to get proxies on the command line.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		count, _ := cmd.Flags().GetInt("count")
		noColor, _ := cmd.Flags().GetBool("no-color")
		raw, _ := cmd.Flags().GetBool("raw")
		printProxies(count, !noColor, raw)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntP("count", "c", 20, "The amount of proxies to get")
	rootCmd.Flags().BoolP("no-color", "", false, "Don't print with colors")
	rootCmd.Flags().BoolP("raw", "", false, "Only print proxy")
}

func colorProxy(proxy proxyfinder.Proxy) {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Println(yellow(proxy.URL.String()), "from", red(proxy.Provider))
}

func noColorProxy(proxy proxyfinder.Proxy) {
	fmt.Println(proxy.URL.String(), "from", proxy.Provider)
}

func printProxies(count int, wantColor bool, raw bool) {
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	if count > len(proxies.All()) {
		fmt.Printf("Error: cannot print %d proxies, only %d avaliable.\n", count, len(proxies.All()))
		return
	}

	switch {
	case raw:
		for i := 0; i < count; i++ {
			proxy := proxies.New()
			fmt.Println(proxy.URL.String())
		}
	case wantColor:
		for i := 0; i < count; i++ {
			colorProxy(proxies.New())
		}
	case !wantColor:
		for i := 0; i < count; i++ {
			noColorProxy(proxies.New())
		}
	}

}
