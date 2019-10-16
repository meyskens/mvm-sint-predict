package main

import (
	"fmt"
	"net"

	"github.com/golang/glog"

	"github.com/meyskens/mvm-sint-predict/pb"
	"github.com/meyskens/mvm-sint-predict/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

type serveCmdOptions struct {
	BindAddr string
	Port     int
}

// NewServeCmd generates the `serve` command
func NewServeCmd() *cobra.Command {
	s := serveCmdOptions{}
	c := &cobra.Command{
		Use:   "serve",
		Short: "Serves the gRPC endpoint",
		Long:  `Serves the gRPC endpoint on the given bind address and port`,
		RunE:  s.RunE,
	}
	c.Flags().StringVarP(&s.BindAddr, "bind-address", "b", "0.0.0.0", "address to bind port to")
	c.Flags().IntVarP(&s.Port, "port", "p", 8080, "Port to listen on")

	viper.BindPFlags(c.Flags())

	return c
}

func (s *serveCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.BindAddr, s.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	glog.Infof("Listening on %s:%d", s.BindAddr, s.Port)

	grpcSrv := grpc.NewServer()
	srv := server.NewSintReplyServer()
	pb.RegisterSintServer(grpcSrv, srv)
	if err := grpcSrv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
