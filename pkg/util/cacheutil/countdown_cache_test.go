package cacheutil

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCountdownCache_SetAndGet(t *testing.T) {
	cache := NewCountdownCache[int](time.Millisecond * 100) // 设置 TTL 为 100 毫秒

	// 设置缓存项
	cache.Set("test1", 123)
	val, found := cache.Get("test1")
	assert.True(t, found, "Expected to find the cache item")
	assert.Equal(t, 123, val, "Expected cache value to be 123")

	// 等待缓存过期
	time.Sleep(time.Millisecond * 150)
	val, found = cache.Get("test1")
	assert.False(t, found, "Expected cache item to be expired")
}

func TestCountdownCache_GetAfterExpiration(t *testing.T) {
	cache := NewCountdownCache[int](time.Millisecond * 100)

	// 设置缓存项
	cache.Set("test2", 456)

	// 等待缓存过期
	time.Sleep(time.Millisecond * 150)
	val, found := cache.Get("test2")
	assert.False(t, found, "Expected cache item to be expired")
	var zeroValue int
	assert.Equal(t, zeroValue, val, "Expected value to be zero for expired cache")
}

func TestCountdownCache_Delete(t *testing.T) {
	cache := NewCountdownCache[int](time.Millisecond * 100)

	// 设置缓存项
	cache.Set("test3", 789)

	// 删除缓存项
	cache.Delete("test3")
	val, found := cache.Get("test3")
	assert.False(t, found, "Expected cache item to be deleted")
	var zeroValue int
	assert.Equal(t, zeroValue, val, "Expected value to be zero for deleted cache")
}

// 定义一个实现了 io.Writer 接口的结构体
type LogBuffer struct {
	buf *bytes.Buffer
}

func (l *LogBuffer) Write(p []byte) (n int, err error) {
	return l.buf.Write(p)
}

func TestCountdownCache_ExpirationLog(t *testing.T) {
	cache := NewCountdownCache[int](time.Millisecond * 100)
	cache.Set("test5", 111)

	// 创建 LogBuffer 实例来捕获日志
	var logBuf LogBuffer
	logBuf.buf = new(bytes.Buffer)

	// 设置 logrus 输出到 logBuf
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(&logBuf)

	// 等待缓存过期
	time.Sleep(time.Millisecond * 150)

	// 确保缓存过期时有日志输出
	assert.Contains(t, logBuf.buf.String(), "cache expired: test5", "Expected log to contain cache expired message")
}
