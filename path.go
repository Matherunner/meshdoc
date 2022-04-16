package meshdoc

import (
	"path"
	"strings"
)

type GenericPath struct {
	path       string
	withoutExt string
}

func NewGenericPath(path string) *GenericPath {
	return &GenericPath{path: path}
}

func (p *GenericPath) Path() string {
	return p.path
}

func (p *GenericPath) WebPath() string {
	// TODO: is this correct or reliable?
	return strings.ReplaceAll(p.path, "\\", "/")
}

func (p *GenericPath) WithoutExt() string {
	if p.withoutExt == "" {
		ext := path.Ext(p.path)
		p.withoutExt = strings.TrimSuffix(p.path, ext)
	}
	return p.withoutExt
}

func (p *GenericPath) SetExt(ext string) *GenericPath {
	p.path = p.WithoutExt() + ext
	return p
}

func (p *GenericPath) Clone() *GenericPath {
	newPath := *p
	return &newPath
}

type GenericPathList []*GenericPath

func (p GenericPathList) Len() int {
	return len(p)
}

func (p GenericPathList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p GenericPathList) Less(i, j int) bool {
	return p[i].Path() < p[j].Path()
}
