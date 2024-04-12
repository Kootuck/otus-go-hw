package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type KeyValue struct {
	Key   Key // храним в списке так же ключ элемента в словаре для быстрого удаления из словаря
	Value interface{}
}

type lruCache struct {
	capacity int
	mu       sync.Mutex
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Добавить элемент в кэш
// return = флаг, присутствовал ли элемент в кэше.
func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	item, exists := lc.items[key]

	//  если элемента нет в словаре, то добавить в словарь и в начало очереди
	// (если превышена ёмкость кэша, то удалить последний элемент из очереди и из словаря);
	if !exists {
		if lc.queue.Len() == lc.capacity {
			lc.removeLastElement()
		}
		lc.items[key] = lc.queue.PushFront(KeyValue{Key: key, Value: value})
	}

	// если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди
	if exists {
		item.Value = KeyValue{Key: key, Value: value}
		lc.queue.MoveToFront(item)
		return true
	}

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	// если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true;
	item, exists := lc.items[key]

	if exists {
		lc.queue.MoveToFront(item)
		kv, ok := item.Value.(KeyValue)
		if ok {
			return kv.Value, true
		}
	}
	// если элемента нет в словаре, то вернуть nil и false.
	return nil, false
}

func (lc *lruCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.items = make(map[Key]*ListItem, lc.capacity)
	lc.queue = NewList()
}

// Удалить последний элемент из очереди и из словаря.
func (lc *lruCache) removeLastElement() {
	ListItem := lc.queue.Back()
	lc.queue.Remove(ListItem)
	// Найти в словаре ключ для удаления по значению.
	kv, ok := ListItem.Value.(KeyValue)
	if ok {
		delete(lc.items, kv.Key)
	}
}
