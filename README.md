# Cache

缓存模块

#### 简述

通过适配器模式实现不同类型的缓存，定义了

#### 已接入缓存

- memory
- redis
- freecache
- bigcache
- memcache

#### 基础接口

- Init
- Get
- Set
- Del
- IsExist

#### RoadMap

- 了解每个类型缓存的适用场景
- 修改init时用的config，现在的config需要用map类型读取，比较麻烦
- 接入GroupCache，ristretto
- 完善基础接口

