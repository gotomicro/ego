package main

import (
	"context"
	"fmt"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/client/ehttp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	if err := ego.New().Invoker(
		invokerHTTP,
		callHTTP,
	).Run(); err != nil {
		elog.Error("startup", elog.FieldErr(err))
	}
}

var httpComp *ehttp.Component

func invokerHTTP() error {
	httpComp = ehttp.Load("http.test").Build()
	return nil
}

func callHTTP() error {
	tracer := etrace.NewTracer(trace.SpanKindClient)

	req := httpComp.R()
	carrier := propagation.HeaderCarrier(req.Header)
	fmt.Printf("carrier--------------->"+"%+v\n", carrier.Keys())

	ctx, span := tracer.Start(context.Background(), "callHTTP()", carrier)
	defer span.End()

	fmt.Printf("carrier--------------->"+"%+v\n", carrier.Keys())

	// Inject traceId Into Header
	//c1 := etrace.HeaderInjector(ctx, req.Header)
	fmt.Println(span.SpanContext().TraceID())
	info, err := req.SetContext(ctx).Get("http://127.0.0.1:9007/hello")
	if err != nil {
		return err
	}
	fmt.Println(info)
	return nil
}
