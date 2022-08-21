package component

type Cache struct {
	Named
}

func NewCache() *Cache {
	return &Cache{Named: componentCache}
}
