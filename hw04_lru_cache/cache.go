package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	myCacheItem := &cacheItem{ // соберём структуру
		key:   key,
		value: value,
	}
	v, ok := l.items[key]
	if ok { //ключ есть в словаре
		if value == v.Value { // то-же значение
		} else { //новое значение
			v.Value = *myCacheItem
		}
		l.queue.MoveToFront(v)
	} else if l.queue.Len() < l.capacity { // не достигли предела. только записываем
		l.items[key] = l.queue.PushFront(*myCacheItem)
	} else { // достигли предела. сперва удаляем
		delete(l.items, l.queue.Back().Value.(cacheItem).key)
		l.queue.Remove(l.queue.Back())
		l.items[key] = l.queue.PushFront(*myCacheItem)
	}
	return ok
}
func (l *lruCache) Get(key Key) (interface{}, bool) {
	v, ok := l.items[key]
	if ok { // ключ есть в словаре
		l.queue.MoveToFront(v)
	}
	if ok {
		return v.Value.(cacheItem).value, ok
	} else {
		return nil, false
	}
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = nil
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
