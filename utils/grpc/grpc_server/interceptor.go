package grpc_server

// func (config *Config) defaultUnaryServerInterceptor() grpc.UnaryServerInterceptor { //nolint
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
// 		// 默认过滤掉该探活日志
// 		if config.EnableSkipHealthLog && info.FullMethod == "/grpc.health.v1.Health/Check" {
// 			return handler(ctx, req)
// 		}

// 		var beg = time.Now()
// 		// 为了性能考虑，如果要加日志字段，需要改变slice大小
// 		var fields = make([]blabel.Label, 0, 20)
// 		var event = "normal"

// 		// try to extract custom keys from request metadata
// 		if md, ok := metadata.FromIncomingContext(ctx); ok {
// 			ctx = transport.CustomKeysMapPropagator.Extract(ctx, transport.GrpcHeaderCarrier(md))
// 			// inject custom keys labels
// 			mp := transport.GetMapFromContext(ctx)
// 			if len(mp) > 0 {
// 				ctx = blog.DecorateLoggerInContext(ctx, func(logger blog.Logger) blog.Logger {
// 					return logger.WithDetails(blabel.Value("customHeaders", mp))
// 				})
// 			}
// 		}

// 		// 此处必须使用defer来recover handler内部可能出现的panic
// 		defer func() {
// 			duration := time.Since(beg)
// 			if rec := recover(); rec != nil {
// 				switch recType := rec.(type) {
// 				case error:
// 					err = recType
// 				default:
// 					err = status.Errorf(codes.Internal, "panic: %v", rec)
// 				}

// 				stack := make([]byte, 4096)
// 				stack = stack[:runtime.Stack(stack, true)]
// 				fields = append(fields, blabel.String("stack", string(stack)))
// 				event = "recover"
// 			}

// 			isSlow := false
// 			if config.SlowLogThreshold > time.Duration(0) && config.SlowLogThreshold < duration {
// 				isSlow = true
// 			}

// 			// 如果没有开启日志组件、并且没有错误，没有慢日志，那么直接返回不记录日志
// 			if err == nil && !config.EnableAccessInterceptor && !isSlow {
// 				return
// 			}

// 			spbStatus := ecodes.Convert(err)
// 			httpStatusCode := ecodes.GrpcToHTTPStatusCode(spbStatus.Code())

// 			fields = append(fields,
// 				zap.String("type", "unary"),
// 				zap.Int64("code", int64(spbStatus.Code())),
// 				zap.Int64("uniformCode", int64(httpStatusCode)),
// 				zap.String("description", spbStatus.Message()),
// 				glog.FieldEvent(event),
// 				glog.FieldMethod(info.FullMethod),
// 				glog.FieldCost(time.Since(beg)),
// 				zap.String("peerIp", getPeerIP(ctx)),
// 			)

// 			span := trace.SpanFromContext(ctx)
// 			// add custom metadata to trace fields
// 			for k, v := range transport.GetMapFromContext(ctx) {
// 				span.SetAttributes(attribute.String(k, v))
// 			}

// 			if config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
// 				fields = append(fields, glog.FieldTrace(etrace.ExtractTraceID(ctx)))
// 			}

// 			if config.EnableAccessInterceptorReq {
// 				fields = append(fields, zap.Any("request", req))
// 				if md, ok := metadata.FromIncomingContext(ctx); ok {
// 					fields = append(fields, zap.Any("metadata", md))
// 				}
// 			}
// 			if config.EnableAccessInterceptorRes {
// 				fields = append(fields, zap.Any("response", res))
// 			}

// 			if isSlow {
// 				blog.Warn(ctx, "slow", fields...)
// 			}

// 			if err != nil {
// 				fields = append(fields, glog.FieldError(err))
// 				// 只记录系统级别错误
// 				if httpStatusCode >= http.StatusInternalServerError {
// 					// 只记录系统级别错误
// 					blog.Error(ctx, "access", fields...)
// 				} else {
// 					// 非核心报错只做warning
// 					blog.Warn(ctx, "access", fields...)
// 				}
// 				return
// 			}

// 			if config.EnableAccessInterceptor {
// 				blog.Info(ctx, "access", fields...)
// 			}
// 		}()

// 		return handler(ctx, req)
// 	}
// }
