package etrace

import (
	"go.opentelemetry.io/otel/attribute"
)

// CustomTag ...
func CustomTag(key string, val string) attribute.KeyValue {
	return attribute.String(key, val)

}

// TagComponent ...
func TagComponent(component string) attribute.KeyValue {
	return attribute.String("component", component)
}

// TagSpanKind ...
func TagSpanKind(kind string) attribute.KeyValue {
	return attribute.String("span.kind", kind)
}

// TagSpanURL ...
func TagSpanURL(url string) attribute.KeyValue {
	return attribute.String("span.url", url)
}
