package command

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "client",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
	addr     string
	timeout  int
	username string
	password string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&addr, "addr", "127.0.0.1:13333", "grpc server address")
	rootCmd.PersistentFlags().StringVar(&username, "user", "admin1", "user name")
	rootCmd.PersistentFlags().StringVar(&password, "password", "secret", "password")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 2500, "timeout millisecond")
}

func authMethods() map[string]bool {
	const deepthoughtServicePath = "/deepthought.Compute/"

	return map[string]bool{
		deepthoughtServicePath + "Boot":  true,
		deepthoughtServicePath + "Infer": true,
	}
}

// Run runs command.
func Run() error {
	return rootCmd.Execute()
}
