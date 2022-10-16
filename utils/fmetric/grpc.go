package fmetric

import "strings"

const unknown = "unknown"

// SplitGrpcMethodName split grpc full method into service and method.
func SplitGrpcMethodName(fullMethodName string) (service string, method string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return unknown, unknown
}
