package command

import (
	"github.com/spf13/cobra"
	"time"
)

var (
	rootCmd = &cobra.Command{
		Use: "client",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
	addr    string
	timeout int
)

const credRefreshingInterval = 500 * time.Millisecond

func init() {
	rootCmd.PersistentFlags().StringVar(&addr, "addr", "127.0.0.1:13333", "grpc server address")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 2500, "timeout millisecond")
}

// Run runs command.
func Run() error {
	return rootCmd.Execute()
}
