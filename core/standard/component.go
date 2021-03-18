package standard

// Component 组件interface
type Component interface {
	Name() string        // 唯一名称，配置key的名称
	PackageName() string // 包名
	Init() error         // 初始化
	Start() error        // 启动
	Stop() error         // 停止
}
