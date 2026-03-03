package svg

import (
	"errors"
	"fmt"
	"io/fs"
	"sync"

	"github.com/vapstack/htm"
)

func New(fs fs.FS, hotReload bool) (*Comp, error) {
	return &Comp{
		fs:  fs,
		hot: hotReload,
	}, nil
}

type Comp struct {
	fs    fs.FS
	hot   bool
	cache sync.Map // map[string]*htm.Node
}

func (c *Comp) Get(name string, mods ...htm.Mod) *htm.Node {
	if !c.hot {
		if v, ok := c.cache.Load(name); ok {
			return v.(*htm.Node).Clone().Apply(mods)
		}
	}
	return c.get(name, mods)
}

func (c *Comp) get(name string, mods []htm.Mod) *htm.Node {
	n, err := c.load(name)
	if err != nil {
		return htm.Span().Text(err.Error())
	}
	if c.hot {
		return n.Apply(mods)
	}
	return n.Clone().Apply(mods)
}

func (c *Comp) load(name string) (*htm.Node, error) {
	if name == "" {
		return nil, errors.New("svg.Get: name is empty")
	}
	var (
		b []byte
		e error
	)
	if hasExt(name) {
		b, e = fs.ReadFile(c.fs, name)
	} else {
		b, e = fs.ReadFile(c.fs, name+".svg")
	}
	if e != nil {
		return nil, e
	}
	n, err := htm.Parse(b, htm.ParseTopLevelRawContent, htm.ParseReuseBuffer)
	if err != nil {
		return nil, err
	}
	if c.hot {
		return n, nil
	}
	v, loaded := c.cache.LoadOrStore(name, n)
	if loaded {
		n.Release()
		return v.(*htm.Node), nil
	}
	return n, nil
}

func hasExt(name string) bool {
	if len(name) < 5 {
		return false
	}
	i := len(name) - 4
	return name[i] == '.' && (name[i+1]|0x20) == 's' && (name[i+2]|0x20) == 'v' && (name[i+3]|0x20) == 'g'
}

// Init initializes a package-level instance to allow direct usage of package function Get.
func Init(fs fs.FS, hotReload bool) error {
	c, err := New(fs, hotReload)
	if err != nil {
		return err
	}
	comp = c
	return nil
}

var comp *Comp

func Get(name string, mods ...htm.Mod) *htm.Node {
	if comp == nil {
		panic(fmt.Errorf("svg.Get: package was not initialized"))
	}
	return comp.Get(name, mods...)
}
