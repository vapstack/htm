package svg

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unsafe"

	"github.com/vapstack/htm"
)

type Setup struct {
	IconFS        fs.FS
	ImageFS       fs.FS
	HotReload     bool
	DefaultHeight string
}

func New(setup Setup) (*Comp, error) {
	c := &Comp{
		iconFS:  setup.IconFS,
		imageFS: setup.ImageFS,
		icons:   make(map[string][]string),
		images:  make(map[string][]string),
		hot:     setup.HotReload,
		height:  setup.DefaultHeight,
	}
	if c.height == "" {
		c.height = "24px"
	}
	return c, c.Reload()
}

type Comp struct {
	iconFS  fs.FS
	imageFS fs.FS

	icons  map[string][]string
	images map[string][]string

	hot bool
	mu  sync.RWMutex

	height string
}

var viewBoxRx = regexp.MustCompile(`viewBox="(.*?)"`)
var contentRx = regexp.MustCompile(`(?s)<svg[^>]*>(.*?)</svg>`)

func (c *Comp) Reload() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := parseFiles(c.iconFS, ".", c.icons); err != nil {
		return err
	}
	if err := parseFiles(c.imageFS, ".", c.images); err != nil {
		return err
	}
	return nil
}

func (c *Comp) hotLoadIcon(name string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := parseFile(c.iconFS, name, c.icons); err != nil {
		return nil, err
	}
	parts := c.icons[strings.TrimSuffix(name, ".svg")]
	return parts, nil
}

func (c *Comp) hotLoadImage(name string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := parseFile(c.imageFS, name, c.images); err != nil {
		return nil, err
	}
	parts := c.images[strings.TrimSuffix(name, ".svg")]
	return parts, nil
}

func (c *Comp) Icon(name string, mods ...htm.Mod) *htm.Node {
	var parts []string
	if c.hot {
		var err error
		if parts, err = c.hotLoadIcon(name + ".svg"); err != nil {
			return htm.Span().Text(err.Error())
		}
	} else {
		parts = c.icons[name]
	}
	if len(parts) == 2 {
		return htm.Build("svg").
			Attr("viewBox", parts[0]).
			Attr("height", c.height).
			Apply(mods).
			Content(htm.RawString(parts[1]))
	}
	return nil
}

func (c *Comp) Image(name string, mods ...htm.Mod) *htm.Node {
	var parts []string
	if c.hot {
		var err error
		if parts, err = c.hotLoadImage(name + ".svg"); err != nil {
			return htm.Span().Text(err.Error())
		}
	} else {
		parts = c.images[name]
	}
	if len(parts) == 2 {
		return htm.Build("svg").
			Attr("viewBox", parts[0]).
			Attr("height", c.height).
			Apply(mods).
			Content(htm.RawString(parts[1]))
	}
	return nil
}

func parseFiles(sfs fs.FS, dir string, target map[string][]string) error {
	entries, _ := fs.ReadDir(sfs, dir)
	for _, entry := range entries {
		name := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if err := parseFiles(sfs, name, target); err != nil {
				return err
			}
			continue
		}
		if filepath.Ext(entry.Name()) != ".svg" {
			continue
		}
		if err := parseFile(sfs, name, target); err != nil {
			return err
		}
	}
	return nil
}

func parseFile(sfs fs.FS, name string, target map[string][]string) error {
	b, err := fs.ReadFile(sfs, name)
	if err != nil {
		return err
	}
	str := unsafe.String(unsafe.SliceData(b), len(b))
	var viewBox string
	if match := viewBoxRx.FindStringSubmatch(str); len(match) > 1 {
		viewBox = match[1]
	} else {
		viewBox = "0 0 32 32"
	}
	if match := contentRx.FindStringSubmatch(str); len(match) > 1 {
		ext := filepath.Ext(name)
		key := name[:len(name)-len(ext)]
		target[key] = []string{viewBox, match[1]}
	}
	return nil
}

// Init initializes a package-level instance to allow direct usage of package functions Icon and Image.
func Init(setup Setup) error {
	c, err := New(setup)
	if err != nil {
		return err
	}
	comp = c
	return nil
}

var comp *Comp

func Icon(name string, mods ...htm.Mod) *htm.Node {
	if comp == nil {
		panic(fmt.Errorf("svg.Icon: package was not initialized"))
	}
	return comp.Icon(name, mods...)
}

func Image(name string, mods ...htm.Mod) *htm.Node {
	if comp == nil {
		panic(fmt.Errorf("svg.Image: package was not initialized"))
	}
	return comp.Image(name, mods...)
}
