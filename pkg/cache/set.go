package cache

type Set[K comparable] struct {
	cache *Hashed[K, struct{}]
}

func NewSet[K comparable](keyFn func(K) K) *Set[K] {
	return &Set[K]{
		cache: NewHashed[K, struct{}](keyFn, nil),
	}
}

func (s *Set[K]) Add(key K) {
	s.cache.Add(key, struct{}{})
}

func (s *Set[K]) Has(key K) bool {
	return s.cache.Has(key)
}
