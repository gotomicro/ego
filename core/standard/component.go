package standard

type Component interface {
	Name() string // 唯一名称
	Init() error  // 初始化
	Start() error // 启动
	Stop() error  // 停止
}
