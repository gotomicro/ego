package etrace

//func TestCompatibleExtractHTTPTraceID(t *testing.T) {
//	header := http.Header{}
//	header.Set("X-Trace-Id", "123:45:6789:abc")
//	CompatibleExtractHTTPTraceID(header)
//	tp := header.Get("Traceparent")
//	assert.Equal(t, "00-12345-45-0abc", tp)
//}
//
//func TestCompatibleExtractGrpcTraceID(t *testing.T) {
//	md := metadata.Pairs("x-trace-id", "123:45:6789:abc")
//	CompatibleExtractGrpcTraceID(md)
//	exp := "00-12345-45-0abc"
//	traceparent := md.Get("Traceparent")
//	assert.Equal(t, exp, traceparent[0])
//
//	// 测试空的 "x-trace-id"
//	emptyMD := metadata.Pairs("x-trace-id", "")
//	CompatibleExtractGrpcTraceID(emptyMD)
//	tp := emptyMD.Get("Traceparent")
//	assert.Equal(t, "", tp[0])
//}
