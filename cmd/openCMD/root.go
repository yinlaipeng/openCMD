package openCMD

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{

	Use:   "openCMD",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Help for openCMD")
}
