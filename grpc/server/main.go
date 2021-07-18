package main

import (
	"fmt"
	"github.com/care0717/deepthought-api/grpc/proto/auth"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	model2 "github.com/care0717/deepthought-api/grpc/server/model"
	repository2 "github.com/care0717/deepthought-api/grpc/server/repository"
	service2 "github.com/care0717/deepthought-api/grpc/server/service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/tls/certprovider/pemfile"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/security/advancedtls"
	"google.golang.org/grpc/security/advancedtls/testdata"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	portNumber             = 13333
	promPort               = 18888
	credRefreshingInterval = 1 * time.Minute
	secretKey              = "secret"
	tokenDuration          = 15 * time.Minute
)

func main() {
	kep := keepalive.EnforcementPolicy{
		MinTime: 10 * time.Second,
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	identityOptions := pemfile.Options{
		CertFile:        testdata.Path("server_cert_1.pem"),
		KeyFile:         testdata.Path("server_key_1.pem"),
		RefreshDuration: credRefreshingInterval,
	}
	identityProvider, err := pemfile.NewProvider(identityOptions)
	if err != nil {
		logger.Fatal(fmt.Sprintf("pemfile.NewProvider(%v) failed: %v", identityOptions, err))
	}
	defer identityProvider.Close()
	rootOptions := pemfile.Options{
		RootFile:        testdata.Path("server_trust_cert_1.pem"),
		RefreshDuration: credRefreshingInterval,
	}
	rootProvider, err := pemfile.NewProvider(rootOptions)
	if err != nil {
		logger.Fatal(fmt.Sprintf("pemfile.NewProvider(%v) failed: %v", rootOptions, err))
	}
	defer rootProvider.Close()

	// Start a server and create a client using advancedtls API with Provider.
	options := &advancedtls.ServerOptions{
		IdentityOptions: advancedtls.IdentityCertificateOptions{
			IdentityProvider: identityProvider,
		},
		RootOptions: advancedtls.RootCertificateOptions{
			RootProvider: rootProvider,
		},
		RequireClientCert: true,
		VerifyPeer: func(params *advancedtls.VerificationFuncParams) (*advancedtls.VerificationResults, error) {
			// This message is to show the certificate under the hood is actually reloaded.
			logger.Info(fmt.Sprintf("Client common name: %s.", params.Leaf.Subject.CommonName))
			return &advancedtls.VerificationResults{}, nil
		},
		VType: advancedtls.CertVerification,
	}
	serverTLSCreds, err := advancedtls.NewServerCreds(options)
	if err != nil {
		logger.Fatal(fmt.Sprintf("advancedtls.NewServerCreds(%v) failed: %v", options, err))
	}

	userStore := repository2.NewInMemoryUserStore()
	if err = seedUsers(userStore); err != nil {
		logger.Fatal(fmt.Sprintf("advancedtls.NewServerCreds(%v) failed: %v", options, err))
	}
	jwtManager := service2.NewJWTManager(secretKey, tokenDuration)
	authServer := NewAuthServer(userStore, jwtManager)
	interceptor := service2.NewAuthInterceptor(jwtManager, accessibleRoles())
	serv := grpc.NewServer(
		grpc.Creds(serverTLSCreds),
		grpc.KeepaliveEnforcementPolicy(kep),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			interceptor.Stream(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			interceptor.Unary(),
		)),
	)
	auth.RegisterAuthServer(serv, authServer)
	deepthought.RegisterComputeServer(serv, &DeepthoughtServer{})
	grpc_prometheus.Register(serv)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		logger.Info(fmt.Sprintf("prometheus metrics listen port: %d", promPort))
		logger.Fatal(fmt.Sprintf("prometheus serve failed: %v", http.ListenAndServe(fmt.Sprintf(":%d", promPort), mux)))
	}()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	logger.Info(fmt.Sprintf("listen port :%d", portNumber))
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to listen: %v", err))
	}
	logger.Fatal(fmt.Sprintf("serve failed: %v", serv.Serve(l)))
}

func seedUsers(userStore repository2.UserStore) error {
	user, err := model2.NewUser("admin1", "secret", "admin")
	if err != nil {
		return err
	}
	err = userStore.Save(user)
	if err != nil {
		return err
	}
	user, err = model2.NewUser("user1", "secret", "user")
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

func accessibleRoles() map[string][]string {
	const deepthoughtServicePath = "/deepthought.Compute/"

	return map[string][]string{
		deepthoughtServicePath + "Boot":  {"admin"},
		deepthoughtServicePath + "Infer": {"admin", "user"},
	}
}
