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

import "github.com/go-redis/redis"

//TODO 引入redis统一错误码

//Component client (cmdable and config)
type Component struct {
	Config *Config
	Client redis.Cmdable
}

// Cluster try to get a redis.ClusterClient
func (r *Component) Cluster() *redis.ClusterClient {
	if c, ok := r.Client.(*redis.ClusterClient); ok {
		return c
	}
	return nil
}

//Stub try to get a redis.Client
func (r *Component) Stub() *redis.Client {
	if c, ok := r.Client.(*redis.Client); ok {
		return c
	}
	return nil
}
