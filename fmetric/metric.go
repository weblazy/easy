package fmetric

var (
	// TypeHTTP ...
	TypeHTTP = "http"
	// TypeGRPCClient ...
	TypeGRPCClient = "grpc_client"
	// TypeGRPCServer ...
	TypeGRPCServer = "grpc_server"
	// TypeRedis ...
	TypeRedis = "redis"
	// TypeGorm ...
	TypeGorm = "gorm"
	// TypeMySQL ...
	TypeMySQL = "mysql"

	// DefaultNamespace ...
	DefaultNamespace = ""
)

const (
	CodeOK    = "OK"
	CodeError = "Error"
)
