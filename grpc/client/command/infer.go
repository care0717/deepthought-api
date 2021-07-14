package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

var (
	inferCmd = &cobra.Command{
		Use:           "infer",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("usage: infer QUERY")
			}
			query := args[0]

			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			cc := deepthought.NewComputeClient(conn)

			ctx, cancel := context.WithCancel(context.Background())
			go func(cancel func()) {
				time.Sleep(time.Duration(timeout) * time.Millisecond)
				cancel()
			}(cancel)

			resp, err := cc.Infer(ctx, &deepthought.InferRequest{Query: query})
			if err != nil {
				return err
			}
			s, err := json.Marshal(resp)
			if err != nil {
				return err
			}
			fmt.Println(string(s))
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(inferCmd)
}
