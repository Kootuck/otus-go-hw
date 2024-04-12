package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic capacity", func(t *testing.T) {
		// на логику выталкивания элементов из-за размера очереди
		// (например: n = 3, добавили 4 элемента - 1й из кэша вытолкнулся);
		lc := NewCache(5)
		for i := 0; i < 10; i++ { // добавили 10 элементов
			lc.Set(Key(strconv.Itoa(i)), i)
		}
		// должен содержать 5..9 элементы
		for i := 5; i < 10; i++ { // добавили 10 элементов
			_, ok := lc.Get(Key(strconv.Itoa(i)))
			require.True(t, ok)
		}
		// не должен содержать 0..4 элементы
		for i := 0; i < 4; i++ { // добавили 10 элементов
			_, ok := lc.Get(Key(strconv.Itoa(i)))
			require.False(t, ok)
		}
	})

	t.Run("purge logic not used", func(t *testing.T) {
		// на логику выталкивания давно используемых элементов
		// (например: n = 3, добавили 3 элемента, обратились несколько раз к разным элементам:

		lc := NewCache(3)
		for i := 0; i < 3; i++ { // добавили 3
			lc.Set(Key(strconv.Itoa(i)), i)
		}
		_, ok := lc.Get("2")
		require.True(t, ok)

		_, ok = lc.Get("1")
		require.True(t, ok)

		// изменили значение, получили значение и пр. - добавили 4й элемент,
		// из первой тройки вытолкнется тот элемент, что был затронут наиболее давно).
		ok = lc.Set("3", 3)
		require.False(t, ok)

		_, ok = lc.Get("0")
		require.False(t, ok)
	})
}

// число итераций изменено 1kk -> 30k, чтобы тест не падал по таймауту.
func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 30_000; i++ {
			require.NotPanics(t, func() { c.Set(Key(strconv.Itoa(i)), i) })
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 30_000; i++ {
			require.NotPanics(t, func() { c.Get(Key(strconv.Itoa(rand.Intn(30_000)))) })
		}
	}()

	wg.Wait()
}

func BenchmarkCache(b *testing.B) {
	c := NewCache(20)
	for i := 0; i < b.N; i++ {
		c.Set(Key(strconv.Itoa(i)), i)
	}
	for i := 0; i < b.N; i++ {
		c.Get(Key(strconv.Itoa(rand.Intn(b.N))))
	}
}
