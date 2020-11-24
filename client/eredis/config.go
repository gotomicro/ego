// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eredis

import (
	"github.com/gotomicro/ego/core/util/xtime"
	"time"
)

const (
	//ClusterMode using clusterClient
	ClusterMode string = "cluster"
	//StubMode using reidsClient
	StubMode string = "stub"
)

// Config for redis, contains RedisStubConfig and RedisClusterConfig
type Config struct {
	Addrs         []string      // Addrs 实例配置地址
	Addr          string        // Addr stubConfig 实例配置地址
	Mode          string        // Mode Redis模式 cluster|stub
	Password      string        // Password 密码
	DB            int           // DB，默认为0, 一般应用不推荐使用DB分片
	PoolSize      int           // PoolSize 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	MaxRetries    int           // MaxRetries 网络相关的错误最大重试次数 默认8次
	MinIdleConns  int           // MinIdleConns 最小空闲连接数
	DialTimeout   time.Duration // DialTimeout 拨超时时间
	ReadTimeout   time.Duration // ReadTimeout 读超时 默认3s
	WriteTimeout  time.Duration // WriteTimeout 读超时 默认3s
	IdleTimeout   time.Duration // IdleTimeout 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	Debug         bool          // Debug开关
	ReadOnly      bool          // ReadOnly 集群模式 在从属节点上启用读模式
	EnableTrace   bool          // 是否开启链路追踪，开启以后。使用DoCotext的请求会被trace
	SlowThreshold time.Duration // 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	OnDialError   string        // OnDialError panic|error
}

// DefaultConfig default config ...
func DefaultConfig() *Config {
	return &Config{
		DB:            0,
		PoolSize:      10,
		MaxRetries:    3,
		MinIdleConns:  100,
		DialTimeout:   xtime.Duration("1s"),
		ReadTimeout:   xtime.Duration("1s"),
		WriteTimeout:  xtime.Duration("1s"),
		IdleTimeout:   xtime.Duration("60s"),
		ReadOnly:      false,
		Debug:         false,
		EnableTrace:   false,
		SlowThreshold: xtime.Duration("250ms"),
		OnDialError:   "panic",
	}
}
