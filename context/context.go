package context

type Context interface {
	Set(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
}

type DefaultContext struct {
	kv map[interface{}]interface{}
}

func NewDefaultContext() Context {
	return &DefaultContext{
		kv: map[interface{}]interface{}{},
	}
}

func (c *DefaultContext) Get(key interface{}) (value interface{}, ok bool) {
	value, ok = c.kv[key]
	return
}

func (c *DefaultContext) Set(key, value interface{}) {
	c.kv[key] = value
}
