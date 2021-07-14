package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

var (
	bootCmd = &cobra.Command{
		Use:           "boot",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			stream, err := cc.Boot(ctx, &deepthought.BootRequest{})
			if err != nil {
				return err
			}

			for {
				resp, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					if status.Code(err) == codes.Canceled {
						break
					}
					return fmt.Errorf("receiving boot response: %w", err)
				}
				s, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				fmt.Println(string(s))
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(bootCmd)
}
