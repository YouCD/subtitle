package cmd

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	Version   string
	commitID  string
	buildTime string
	goVersion string
	buildUser string
)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print the version info of %s", Name),
	Run: func(cmd *cobra.Command, args []string) {
		green := color.FgGreen.Render

		fmt.Printf("Version:   %s\n", green(Version))
		fmt.Printf("CommitID:  %s\n", green(commitID))
		fmt.Printf("BuildTime: %s\n", green(buildTime))
		fmt.Printf("GoVersion: %s\n", green(goVersion))
		fmt.Printf("BuildUser: %s\n", green(buildUser))
	},
}
