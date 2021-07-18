package command

import (
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/tls/certprovider/pemfile"
	"google.golang.org/grpc/security/advancedtls"
	"google.golang.org/grpc/security/advancedtls/testdata"
	"time"
)

const credRefreshingInterval = 500 * time.Millisecond

func getCredentials() (credentials.TransportCredentials, error) {
	identityOptions := pemfile.Options{
		CertFile:        testdata.Path("client_cert_1.pem"),
		KeyFile:         testdata.Path("client_key_1.pem"),
		RefreshDuration: credRefreshingInterval,
	}
	identityProvider, err := pemfile.NewProvider(identityOptions)
	if err != nil {
		return nil, err
	}
	rootOptions := pemfile.Options{
		RootFile:        testdata.Path("client_trust_cert_1.pem"),
		RefreshDuration: credRefreshingInterval,
	}
	rootProvider, err := pemfile.NewProvider(rootOptions)
	if err != nil {
		return nil, err
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
	return advancedtls.NewClientCreds(options)
}
