package cache

import "time"

const (
	MEMORY_CACHE = "memory"
	FREE_CACHE   = "FreeCache"
	REDIS_CACHE  = "Redis"
	MEM_CACHE    = "Memcache"
	BIG_CACHE    = "BigCache"
)

type RedisConf struct {
	Master          string
	Addr            []string
	Password        string
	Db              int
	MaxRetries      int
	MaxRetryBackoff int
}

type MemoryConf struct {
	Interval int
}

type MemCacheConf struct {
	Addr []string
}

type FreeCacheConf struct {
	Size int
}

type BigConf struct {
	Shards             int
	LifeWindow         time.Duration
	CleanWindow        time.Duration
	MaxEntriesInWindow int
	MaxEntrySize       int
	Verbose            bool
}
