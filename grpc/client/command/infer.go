package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/care0717/deepthought-api/grpc/client/service"
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

			clientTLSCreds, err := getCredentials()
			if err != nil {
				return err
			}
			cc1, err := grpc.Dial(addr, grpc.WithTransportCredentials(clientTLSCreds))
			if err != nil {
				return err
			}
			authClient := service.NewAuthClient(cc1, username, password)
			interceptor, err := service.NewAuthInterceptor(authClient, authMethods())
			if err != nil {
				return err
			}
			conn, err := grpc.Dial(
				addr,
				grpc.WithTransportCredentials(clientTLSCreds),
				grpc.WithUnaryInterceptor(interceptor.Unary()),
				grpc.WithStreamInterceptor(interceptor.Stream()),
			)
			if err != nil {
				return err
			}
			defer conn.Close()

			cc := deepthought.NewComputeClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
			defer cancel()

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
