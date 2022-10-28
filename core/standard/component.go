package standard

// Component defines a pluggable and cohesive module. For server Component,
// when Instantiated, it's hole lifecycle will be scheduled by Ego app framework.
//
// User can implement their own Component, if they register their Server Component to Ego app,
// it's starting and stopping will also be controlled by framework.
type Component interface {
	// Name defines component's unique name, such as "my.grpc.server".
	Name() string
	// PackageName presents component's package name, such as "server.grpc".
	PackageName() string
	// Init defines component's Instantiation procedure.
	Init() error
	// Start defines component's start procedure.
	Start() error
	// Stop defines component's stop procedure.
	Stop() error
}
