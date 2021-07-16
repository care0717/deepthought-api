package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/tls/certprovider/pemfile"
	"google.golang.org/grpc/security/advancedtls"
	"google.golang.org/grpc/security/advancedtls/testdata"
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

			identityOptions := pemfile.Options{
				CertFile:        testdata.Path("client_cert_1.pem"),
				KeyFile:         testdata.Path("client_key_1.pem"),
				RefreshDuration: credRefreshingInterval,
			}
			identityProvider, err := pemfile.NewProvider(identityOptions)
			if err != nil {
				return err
			}
			rootOptions := pemfile.Options{
				RootFile:        testdata.Path("client_trust_cert_1.pem"),
				RefreshDuration: credRefreshingInterval,
			}
			rootProvider, err := pemfile.NewProvider(rootOptions)
			if err != nil {
				return err
			}
			options := &advancedtls.ClientOptions{
				IdentityOptions: advancedtls.IdentityCertificateOptions{
					IdentityProvider: identityProvider,
				},
				VerifyPeer: func(params *advancedtls.VerificationFuncParams) (*advancedtls.VerificationResults, error) {
					return &advancedtls.VerificationResults{}, nil
				},
				RootOptions: advancedtls.RootCertificateOptions{
					RootProvider: rootProvider,
				},
				VType: advancedtls.CertVerification,
			}
			clientTLSCreds, err := advancedtls.NewClientCreds(options)
			if err != nil {
				return err
			}
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(clientTLSCreds))
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
