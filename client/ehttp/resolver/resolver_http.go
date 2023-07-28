package resolver

type baseHttpBuilder struct {
	name string
}

type baseHttpResolver struct {
}

func init() {
	m["http"] = &baseHttpBuilder{}
	m["https"] = &baseHttpBuilder{}
	// 用户是没有传协议
	m[""] = &baseHttpBuilder{}
}

// Build ...
func (b *baseHttpBuilder) Build(addr string) (Resolver, error) {
	return &baseHttpResolver{}, nil
}

// Scheme ...
func (b *baseHttpBuilder) Scheme() string {
	return b.name
}

func (b *baseHttpResolver) GetAddr() string {
	return ""
}
