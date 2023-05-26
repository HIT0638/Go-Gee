package geecache

import (
	"fmt"
	"log"
	"sync"
)

// A Getter loads data for a key
type Getter interface {
	Get(key string) ([]byte, error)
}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	maincache cache
}

var (
	mu     sync.Mutex
	groups = make(map[string]*Group)
)

// A GetterFunc implements Getter with a function
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// NewGroup creates a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nli Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		maincache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g

	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(key string) *Group {
	mu.Lock()
	g := groups[key]
	mu.Unlock()

	return g

}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.maincache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)

	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)

	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.maincache.add(key, value)
}
