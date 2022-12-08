package http_client

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

var (
	Prefix  = "easy-"
	TraceId = Prefix + "traceid"
)

const (
	XRequestId      = "x-request-id"
	XB3TraceId      = "x-b3-traceid"
	XB3SpanId       = "x-b3-spanid"
	XB3ParentSpanId = "x-b3-parentspanid"
	XB3Sampled      = "x-b3-sampled"
	XB3Flags        = "x-b3-flags"
	B3              = "b3"
	XOtSpanContext  = "x-ot-span-context"
)

type TraceHeader struct {
	HttpHeader http.Header
	GrpcMd     metadata.MD
}

func SetHttp(header http.Header) *TraceHeader {

	return &TraceHeader{
		HttpHeader: header,
		GrpcMd:     httpToGrpc(header),
	}
}

func SetGrpc(ctx context.Context) *TraceHeader {
	headersIn, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &TraceHeader{}
	}

	return &TraceHeader{
		GrpcMd:     headersIn,
		HttpHeader: grpcToHttp(headersIn),
	}
}

func SetHeader(header interface{}) *TraceHeader {
	switch header := header.(type) {
	case http.Header:
		return SetHttp(header)
	case context.Context:
		return SetGrpc(header)
	case *TraceHeader:
		return header
	default:
		return &TraceHeader{}
	}
}

func grpcToHttp(headersIn metadata.MD) http.Header {
	httpHeader := http.Header{}

	requestId := headersIn.Get(XRequestId)
	traceId := headersIn.Get(XB3TraceId)
	spanId := headersIn.Get(XB3SpanId)
	panrentSpanId := headersIn.Get(XB3ParentSpanId)
	sampled := headersIn.Get(XB3Sampled)
	flags := headersIn.Get(XB3Flags)
	spanContext := headersIn.Get(XOtSpanContext)
	b3 := headersIn.Get(B3)

	if len(requestId) > 0 {
		httpHeader.Add(XRequestId, requestId[0])
	}
	if len(traceId) > 0 {
		httpHeader.Add(XB3TraceId, traceId[0])
	}
	if len(spanId) > 0 {
		httpHeader.Add(XB3SpanId, spanId[0])
	}
	if len(panrentSpanId) > 0 {
		httpHeader.Add(XB3ParentSpanId, panrentSpanId[0])
	}
	if len(sampled) > 0 {
		httpHeader.Add(XB3Sampled, sampled[0])
	}
	if len(flags) > 0 {
		httpHeader.Add(XB3Flags, flags[0])
	}
	if len(spanContext) > 0 {
		httpHeader.Add(XOtSpanContext, spanContext[0])
	}
	if len(b3) > 0 {
		httpHeader.Add(B3, b3[0])
	}

	return httpHeader
}

func httpToGrpc(header http.Header) metadata.MD {

	medata := map[string]string{}

	if header.Get(XRequestId) != "" {
		medata[XRequestId] = header.Get(XRequestId)
	}
	if header.Get(XB3TraceId) != "" {
		medata[XB3TraceId] = header.Get(XB3TraceId)
	}
	if header.Get(XB3SpanId) != "" {
		medata[XB3SpanId] = header.Get(XB3SpanId)
	}
	if header.Get(XB3ParentSpanId) != "" {
		medata[XB3ParentSpanId] = header.Get(XB3ParentSpanId)
	}
	if header.Get(XB3Sampled) != "" {
		medata[XB3Sampled] = header.Get(XB3Sampled)
	}
	if header.Get(XB3Flags) != "" {
		medata[XB3Flags] = header.Get(XB3Flags)
	}
	if header.Get(XOtSpanContext) != "" {
		medata[XOtSpanContext] = header.Get(XOtSpanContext)
	}
	if header.Get(B3) != "" {
		medata[B3] = header.Get(B3)
	}

	return metadata.New(medata)
}
