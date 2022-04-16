package meshdoc

import "sort"

type Context struct {
	kv               map[interface{}]interface{}
	config           *MeshdocConfig
	sortedInputFiles []GenericPath
}

func NewContext() *Context {
	return &Context{
		kv: map[interface{}]interface{}{},
	}
}

func (c *Context) Config() *MeshdocConfig {
	return c.config
}

func (c *Context) SetConfig(config *MeshdocConfig) {
	c.config = config
}

// InputFiles returns the sorted list of the generic paths of input files.
func (c *Context) InputFiles() []GenericPath {
	return c.sortedInputFiles
}

// SetInputFiles sorts and stores the input file paths.
func (c *Context) SetInputFiles(paths []GenericPath) {
	sort.Sort(GenericPathList(paths))
	c.sortedInputFiles = paths
}

func (c *Context) Get(key interface{}) (value interface{}, ok bool) {
	value, ok = c.kv[key]
	return
}

func (c *Context) Set(key, value interface{}) {
	c.kv[key] = value
}
