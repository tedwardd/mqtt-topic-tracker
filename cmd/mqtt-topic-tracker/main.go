package main

import "github.com/spf13/cobra"

func main() {
	basics := commandBasics{}
	rootCmd := newRootCmd(&basics)
	cobra.CheckErr(rootCmd.Execute())
}
