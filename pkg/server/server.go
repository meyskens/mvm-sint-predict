package server

// SintReplyServer implements the Sint gRPC service for processing data
type SintReplyServer struct{}

// NewSintReplyServer gives a new SintReplyServer instance
func NewSintReplyServer() *SintReplyServer {
	return &SintReplyServer{}
}
