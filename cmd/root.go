package cmd

import (
	"fmt"
	"os"

	"github.com/cpendery/cq/cq/engine"
	"github.com/cpendery/cq/cq/shell"
	"github.com/cpendery/cq/internal"
	"github.com/cpendery/cq/internal/format"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	_ "github.com/lib/pq"
)

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [flags]", internal.AppName),
	Short: fmt.Sprintf("%s is a universal terminal-based database front-end", internal.AppName),
	Long: format.Tprintf(`cq is a universal terminal-based database front-end.
complete documentation is available at https://github.com/cpendery/cq
supports the following database sources: {{.supportedDBs}}
`, map[string]interface{}{
		"supportedDBs": internal.SupportedDatabases,
	}),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runRootCmd,
}

func setRootFlags(flags *pflag.FlagSet) {
	flags.StringP("username", "u", "", "the username to use when connecting to the server.")
	flags.StringP("port", "p", "", "the TCP/IP port or the local Unix-domain socket file on which the server is listening for connections")
	flags.StringP("dbname", "d", "", "the name of the database to connect to")
	flags.StringP("host", "h", "", "the host name of the machine on which the server is running")
	flags.StringP("passfile", "P", "", "the file to read as password from")
	flags.StringP("type", "t", "",
		fmt.Sprintf("the database type being connected to: %+v. Allows %s to avoid extra network attempts to determine the database type",
			internal.SupportedDatabases, internal.AppName),
	)
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "", false, fmt.Sprintf("help for %s", internal.AppName))
	setRootFlags(rootCmd.Flags())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	dbEngine := engine.NewEngine()
	if err := dbEngine.Connect(
		cmd.Flag("username").Value.String(),
		cmd.Flag("dbname").Value.String(),
		cmd.Flag("host").Value.String(),
		cmd.Flag("port").Value.String(),
		cmd.Flag("passfile").Value.String(),
		cmd.Flag("type").Value.String(),
	); err != nil {
		return err
	}
	appShell := shell.NewShell()
	return appShell.Start()
}
