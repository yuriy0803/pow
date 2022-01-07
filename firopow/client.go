//go:generate ../.bin/gen-lookup -package firopow -cacheInit 16777216 -cacheGrowth 131072 -datasetInit 1610612736 -datasetGrowth 8388608

package firopow

import (
	"runtime"

	"github.com/twozopw/pow/internal/common"
	"github.com/twozopw/pow/internal/dag"
)

type Client struct {
	*dag.DAG
}

func New(cfg dag.Config) *Client {
	client := &Client{
		DAG: dag.New(cfg),
	}

	return client
}

func NewFiro() *Client {
	var cfg = dag.Config{
		Name:       "FIRO",
		Revision:   23,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   (1 << 30) + (1 << 29),
		DatasetGrowthBytes: 1 << 23,
		CacheInitBytes:     1 << 24,
		CacheGrowthBytes:   1 << 17,

		CacheSizes:   dag.NewLookupTable(cacheSizes, 2048),
		DatasetSizes: dag.NewLookupTable(datasetSizes, 2048),

		DatasetParents:  512,
		EpochLength:     1300,
		SeedEpochLength: 1300,

		CacheRounds:    3,
		CachesCount:    3,
		CachesLockMmap: false,

		L1Enabled:       true,
		L1CacheSize:     4096 * 4,
		L1CacheNumItems: 4096,
	}

	return New(cfg)
}

func (c *Client) Compute(height, nonce uint64, hash []byte) ([]byte, []byte) {
	epoch := c.CalcEpoch(height)
	datasetSize := c.DatasetSize(epoch)
	cache := c.GetCache(epoch)
	lookup := c.NewLookupFunc2048(cache, epoch)

	mix, digest := firopow(hash, height, nonce, datasetSize, lookup, cache.L1())
	runtime.KeepAlive(cache)

	return mix, digest
}
