package htm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Mod represents a function that modifies a Node.
type Mod = func(n *Node)

// Node represents an HTML element or a special rendering node (e.g. text/raw/group).
type Node struct {
	tag  string
	flag byte

	attrs *attrMap
	class *classMap
	vars  []valueEntry

	content []*Node
	slots   []slotNode

	postponed []Mod

	writeFn func(*Node, io.Writer) error // render override

	value TypedValue

	attached []*Node
	acquired atomic.Bool
}

type slotNode struct {
	name  string
	group *Node

	// content []*Node
}

type valueEntry struct {
	name  string
	value TypedValue
}

const (
	flagVoid = 1 << iota
	flagOwned
	flagScript
)

/**/

// Build retrieves a Node from the pool, sets its tag, and applies the provided modifiers.
func Build(tag string, mods ...Mod) *Node {
	n := Get()
	n.SetTag(tag)
	for _, mod := range mods {
		if mod != nil {
			mod(n)
		}
	}
	return n
}

// Group creates a logical group of nodes that renders its children without a wrapping parent tag.
func Group(nodes ...*Node) *Node {
	n := Get()
	n.tag = "$group"
	n.content = append(n.content, nodes...)
	n.writeFn = renderGroup
	return n
}

// RawString creates a raw HTML node from a string.
// The content is written directly to the output without escaping.
func RawString(s string) *Node {
	if len(s) == 0 {
		return nil
	}
	n := Get()
	n.tag = "$raw"
	n.value = String(s)
	n.writeFn = renderRaw
	return n
}

// RawBytes creates a raw HTML node from a byte slice.
// The content is written directly to the output without escaping.
func RawBytes(b []byte) *Node {
	if len(b) == 0 {
		return nil
	}
	n := Get()
	n.tag = "$raw"
	n.value = Bytes(b)
	n.writeFn = renderRaw
	return n
}

// Mods combines multiple modifiers into a single Mod.
func Mods(mods ...Mod) Mod {
	for _, mod := range mods {
		if mod != nil {
			return func(n *Node) {
				for _, m := range mods {
					if m != nil {
						m(n)
					}
				}
			}
		}
	}
	return nil
}

// If conditionally returns a node produced by fn when cond is true.
// If cond is false, If returns nil.
func If(cond bool, fn func() *Node) *Node {
	if cond {
		return fn()
	}
	return nil
}

// ModIf conditionally returns a modifier produced by fn when cond is true.
// If cond is false, ModIf returns nil.
func ModIf(cond bool, fn func() Mod) Mod {
	if cond {
		return fn()
	}
	return nil
}

/**/

// Mod applies the provided modifiers to the node.
func (n *Node) Mod(mods ...Mod) *Node { return n.Apply(mods) }

// Apply applies a slice of modifiers to the node.
func (n *Node) Apply(mods []Mod) *Node {
	for _, m := range mods {
		if m != nil {
			m(n)
		}
	}
	return n
}

// If executes fn on the node if cond is true.
func (n *Node) If(cond bool, fn func(*Node)) *Node {
	if cond {
		fn(n)
	}
	return n
}

// DefaultContent postpones a function that executes fn if the node's content is empty.
func (n *Node) DefaultContent(fn func(*Node)) *Node {
	n.postponed = append(n.postponed, func(n *Node) {
		if n.NoContent() {
			fn(n)
		}
	})
	return n
}

// ModIf executes fn and returns the resulting Mod if cond is true.
func (n *Node) ModIf(cond bool, fn func(*Node) Mod) Mod {
	if cond {
		return fn(n)
	}
	return nil
}

/**/

// GetTag returns the current tag name and whether it is a void (self-closing) element.
func (n *Node) GetTag() (string, bool) { return n.tag, n.flag&flagVoid != 0 }

// Tag returns a Mod that sets the HTML tag name.
func Tag(tag string) Mod { return func(n *Node) { n.SetTag(tag) } }

// TagEx returns a Mod that sets the HTML tag name and void status.
func TagEx(tag string, void bool) Mod { return func(n *Node) { n.SetTagEx(tag, void) } }

// Tag sets the HTML tag name of the node.
func (n *Node) Tag(tag string) *Node { return n.SetTag(tag) }

// SetTag sets the HTML tag name of the node.
func (n *Node) SetTag(tag string) *Node {
	n.SetTagEx(tag, false)
	return n
}

// SetTagEx sets the HTML tag name and explicitly controls the void element status.
func (n *Node) SetTagEx(tag string, void bool) *Node {
	if tag == "" {
		return n
	}
	switch tag {
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "source", "track", "wbr":
		n.tag = tag
		n.flag |= flagVoid
	default:
		n.tag = tag
		if void {
			n.flag |= flagVoid
		} else {
			n.flag &^= flagVoid
		}
	}
	return n
}

/**/

// GetAttr retrieves the attribute value by name. Returns a zero value if not found.
// Use TypedValue.Valid to check for validity.
func (n *Node) GetAttr(name string) TypedValue {
	v, _ := n.attrs.get(name)
	return v
}

// Attr returns a Mod that sets an attribute.
// If value is omitted, it sets a boolean attribute.
// Attr should not be used to set a class attribute; use Class instead.
func Attr(name string, value ...string) Mod {
	if len(value) > 0 {
		v := value[0]
		return func(n *Node) { n.Attr(name, v) }
	}
	return func(n *Node) { n.Attr(name) }
}

// AttrBool returns a Mod that sets a boolean attribute.
// If value is omitted, it is treated as enabled.
func AttrBool(name string, value ...bool) Mod {
	if len(value) > 0 {
		v := value[0]
		return func(n *Node) { n.AttrBool(name, v) }
	}
	return func(n *Node) { n.Attr(name) }
}

// AttrValue returns a Mod that sets an attribute.
// If value is omitted, it sets a boolean attribute.
// AttrValue should not be used to set a class attribute; use Class instead.
func AttrValue(name string, value ...TypedValue) Mod {
	if len(value) > 0 {
		v := value[0]
		return func(n *Node) { n.AttrValue(name, v) }
	}
	return func(n *Node) { n.AttrValue(name) }
}

// Attr sets the value of an attribute.
// If value is omitted, it sets a boolean attribute.
// To unset a boolean attribute, use BoolAttr(name, false) or RemoveAttr(name) or AttrValue(name, Unset).
// Attr should not be used to set a class attribute; use Class instead.
func (n *Node) Attr(name string, value ...string) *Node {
	if len(value) > 0 {
		v := value[0]
		n.attrs.set(name, String(v))
		return n
	}
	n.attrs.set(name, Bool(true))
	return n
}

// A is an alias for Attr
func (n *Node) A(name string, value ...string) *Node { return n.Attr(name, value...) }

// AttrBool sets the value of an attribute.
// If value is omitted, it is treated as enabled.
// To unset a boolean attribute, use BoolAttr(name, false) or RemoveAttr(name) or AttrValue(name, Unset).
func (n *Node) AttrBool(name string, value ...bool) *Node {
	if len(value) > 0 {
		v := value[0]
		n.attrs.set(name, Bool(v))
		return n
	}
	n.attrs.set(name, Bool(true))
	return n
}

// AB is an alias for AttrValue
func (n *Node) AB(name string, value ...bool) *Node { return n.AttrBool(name, value...) }

// AttrValue sets the value of an attribute.
// If value is omitted, it sets a boolean attribute.
// To unset a boolean attribute, use AttrValue(name, Unset) or RemoveAttr(name).
// AttrValue should not be used to set a class attribute; use Class instead.
func (n *Node) AttrValue(name string, value ...TypedValue) *Node {
	if len(value) > 0 {
		v := value[0]
		n.attrs.set(name, v)
		return n
	}
	n.attrs.set(name, Bool(true))
	return n
}

// AV is an alias for AttrValue
func (n *Node) AV(name string, value ...TypedValue) *Node { return n.AttrValue(name, value...) }

// RemoveAttr removes the specified attributes from the node.
func (n *Node) RemoveAttr(names ...string) *Node {
	for _, name := range names {
		n.attrs.set(name, Unset)
	}
	return n
}

// HasAttr checks if the node has at least one of the specified attributes.
func (n *Node) HasAttr(names ...string) bool { return n.attrs.hasAny(names...) }

// HasAttrAll checks if the node has all of the specified attributes.
func (n *Node) HasAttrAll(names ...string) bool { return n.attrs.hasAll(names...) }

// HasAttrPrefix checks if the node has any attribute starting with the given prefix.
func (n *Node) HasAttrPrefix(prefix string) bool { return n.attrs.hasPrefix(prefix) }

// HasAttrSuffix checks if the node has any attribute ending with the given suffix.
func (n *Node) HasAttrSuffix(suffix string) bool { return n.attrs.hasSuffix(suffix) }

// EachAttr iterates over all attributes, calling fn for each.
// Iteration stops if fn returns false.
func (n *Node) EachAttr(fn func(string, TypedValue) bool) *Node {
	n.attrs.each(fn)
	return n
}

// MoveAttrTo moves specific attributes from the current node to the destination node.
func (n *Node) MoveAttrTo(dst *Node, names ...string) *Node {
	if n == dst {
		return n
	}
	for _, name := range names {
		if v, ok := n.attrs.extract(name); ok {
			dst.attrs.set(name, v)
		}
	}
	return n
}

// MoveAttrPrefixTo moves all attributes starting with the given prefix to the destination node.
func (n *Node) MoveAttrPrefixTo(dst *Node, prefix string) *Node {
	if n == dst {
		return n
	}
	n.attrs.movePrefixTo(dst.attrs, prefix)
	return n
}

// MoveAttrSuffixTo moves all attributes ending with the given suffix to the destination node.
func (n *Node) MoveAttrSuffixTo(dst *Node, suffix string) *Node {
	if n == dst {
		return n
	}
	n.attrs.moveSuffixTo(dst.attrs, suffix)
	return n
}

// MoveAttr returns a Mod that moves specific attributes to a destination node.
func (n *Node) MoveAttr(names ...string) Mod {
	return func(dst *Node) { n.MoveAttrTo(dst, names...) }
}

// MoveAttrPrefix returns a Mod that moves attributes with a specific prefix to a destination node.
func (n *Node) MoveAttrPrefix(prefix string) Mod {
	return func(dst *Node) { n.MoveAttrPrefixTo(dst, prefix) }
}

// MoveAttrSuffix returns a Mod that moves attributes with a specific suffix to a destination node.
func (n *Node) MoveAttrSuffix(suffix string) Mod {
	return func(dst *Node) { n.MoveAttrSuffixTo(dst, suffix) }
}

/**/

// Class returns a Mod that adds a class name to the node.
func Class(class string) Mod { return func(n *Node) { n.Class(class) } }

// Class adds a class name to the node. Multiple classes can be separated by spaces.
func (n *Node) Class(name string) *Node {
	n.class.setMulti(name, true)
	return n
}

// C is an alias for Class.
func (n *Node) C(name string) *Node { return n.Class(name) }

// RemoveClass removes the specified class names from the node.
func (n *Node) RemoveClass(names ...string) *Node {
	for _, name := range names {
		n.class.setMulti(name, false)
	}
	return n
}

// HasClass checks if the node has at least one of the specified classes.
func (n *Node) HasClass(names ...string) bool { return n.class.hasAny(names...) }

// HasClassAll checks if the node has all of the specified classes.
func (n *Node) HasClassAll(names ...string) bool { return n.class.hasAll(names...) }

// HasClassPrefix checks if the node has any class starting with the given prefix.
func (n *Node) HasClassPrefix(prefix string) bool { return n.class.hasPrefix(prefix) }

// HasClassSuffix checks if the node has any class ending with the given suffix.
func (n *Node) HasClassSuffix(suffix string) bool { return n.class.hasSuffix(suffix) }

// EachClass iterates over all active classes, calling fn for each.
// Iteration stops if fn returns false.
func (n *Node) EachClass(fn func(string) bool) *Node {
	n.class.each(fn)
	return n
}

// MoveClassTo moves specific classes from the current node to the destination node.
func (n *Node) MoveClassTo(dst *Node, names ...string) *Node {
	for _, name := range names {
		if n.class.extract(name) {
			dst.class.setOne(name, true)
		}
	}
	return n
}

// CopyClassPrefixTo copies classes starting with the given prefixes to the destination node.
func (n *Node) CopyClassPrefixTo(dst *Node, prefixes ...string) *Node {
OUTER:
	for i := len(n.class.o) - 1; i >= 0; i-- {
		e := n.class.o[i]
		if !e.active {
			continue
		}
		for k := i + 1; k < len(n.class.o); k++ {
			if n.class.o[k].name == e.name {
				continue OUTER
			}
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(e.name, prefix) {
				dst.class.setOne(e.name, true)
				break
			}
		}
	}
	return n
}

// MoveClassPrefixTo moves classes starting with the given prefixes to the destination node.
func (n *Node) MoveClassPrefixTo(dst *Node, prefixes ...string) *Node {
OUTER:
	for i := len(n.class.o) - 1; i >= 0; i-- {
		e := n.class.o[i]
		if !e.active {
			continue
		}
		for k := i + 1; k < len(n.class.o); k++ {
			if n.class.o[k].name == e.name {
				continue OUTER
			}
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(e.name, prefix) {
				dst.class.setOne(e.name, true)
				n.class.o[i].active = false
				break
			}
		}
	}
	return n
}

// CopyClassSuffixTo copies classes ending with the given suffixes to the destination node.
func (n *Node) CopyClassSuffixTo(dst *Node, suffixes ...string) *Node {
OUTER:
	for i := len(n.class.o) - 1; i >= 0; i-- {
		e := n.class.o[i]
		if !e.active {
			continue
		}
		for k := i + 1; k < len(n.class.o); k++ {
			if n.class.o[k].name == e.name {
				continue OUTER
			}
		}
		for _, suffix := range suffixes {
			if strings.HasSuffix(e.name, suffix) {
				dst.class.setOne(e.name, true)
				break
			}
		}
	}
	return n
}

// MoveClassSuffixTo moves classes ending with the given suffixes to the destination node.
func (n *Node) MoveClassSuffixTo(dst *Node, suffixes ...string) *Node {
OUTER:
	for i := len(n.class.o) - 1; i >= 0; i-- {
		e := n.class.o[i]
		if !e.active {
			continue
		}
		for k := i + 1; k < len(n.class.o); k++ {
			if n.class.o[k].name == e.name {
				continue OUTER
			}
		}
		for _, suffix := range suffixes {
			if strings.HasSuffix(e.name, suffix) {
				dst.class.setOne(e.name, true)
				n.class.o[i].active = false
				break
			}
		}
	}
	return n
}

// MoveClass returns a Mod that moves specific classes to a destination node.
func (n *Node) MoveClass(names ...string) Mod {
	return func(dst *Node) { n.MoveClassTo(dst, names...) }
}

// CopyClassPrefix returns a Mod that copies classes with specific prefixes to a destination node.
func (n *Node) CopyClassPrefix(prefixes ...string) Mod {
	return func(dst *Node) { n.CopyClassPrefixTo(dst, prefixes...) }
}

// MoveClassPrefix returns a Mod that moves classes with specific prefixes to a destination node.
func (n *Node) MoveClassPrefix(prefixes ...string) Mod {
	return func(dst *Node) { n.MoveClassPrefixTo(dst, prefixes...) }
}

// CopyClassSuffix returns a Mod that copies classes with specific suffixes to a destination node.
func (n *Node) CopyClassSuffix(suffixes ...string) Mod {
	return func(dst *Node) { n.CopyClassSuffixTo(dst, suffixes...) }
}

// MoveClassSuffix returns a Mod that moves classes with specific suffixes to a destination node.
func (n *Node) MoveClassSuffix(suffixes ...string) Mod {
	return func(dst *Node) { n.MoveClassSuffixTo(dst, suffixes...) }
}

/**/

// GetVar retrieves the value of a user variable by name. Returns unset value if not found.
func (n *Node) GetVar(name string) TypedValue {
	for _, v := range n.vars {
		if v.name == name {
			if v.value.Valid() {
				return v.value
			} else {
				return TypedValue{}
			}
		}
	}
	return TypedValue{}
}

// HasVar checks if the node has at least one of the specified variables.
func (n *Node) HasVar(names ...string) bool {
	for _, v := range n.vars {
		for _, name := range names {
			if name == v.name {
				return v.value.Valid()
			}
		}
	}
	return false
}

// HasVarAll checks if the node has all the specified variables.
func (n *Node) HasVarAll(names ...string) bool {
NAMES:
	for _, name := range names {
		for _, v := range n.vars {
			if v.name == name {
				if v.value.Valid() {
					continue NAMES
				} else {
					break
				}
			}
		}
		return false
	}
	return true
}

// Var returns a Mod that attaches arbitrary user data (variable) to the node.
// These variables are not rendered to HTML.
func Var(name string, value string) Mod { return func(n *Node) { n.Var(name, value) } }

// VarValue returns a Mod that attaches arbitrary user data (variable) to the node.
// These variables are not rendered to HTML.
func VarValue(name string, value TypedValue) Mod { return func(n *Node) { n.VarValue(name, value) } }

// Var attaches arbitrary user data (variable) to the node.
// These variables are not rendered to HTML.
func (n *Node) Var(name string, value string) *Node {
	for i, v := range n.vars {
		if v.name == name {
			n.vars[i].value = String(value)
			return n
		}
	}
	n.vars = append(n.vars, valueEntry{name: name, value: String(value)})
	return n
}

// V is an alias for Var.
func (n *Node) V(name string, value string) *Node { return n.Var(name, value) }

// VV is an alias for VarValue.
func (n *Node) VV(name string, value string) *Node { return n.Var(name, value) }

// VarValue attaches arbitrary user data (variable) to the node.
// These variables are not rendered to HTML.
func (n *Node) VarValue(name string, value TypedValue) *Node {
	for i, v := range n.vars {
		if v.name == name {
			n.vars[i].value = value
			return n
		}
	}
	if !value.Valid() {
		return n
	}
	n.vars = append(n.vars, valueEntry{name: name, value: value})
	return n
}

// RemoveVar removes the specified variables from the node.
func (n *Node) RemoveVar(names ...string) *Node {
	for _, name := range names {
		for i, v := range n.vars {
			if v.name == name {
				n.vars[i].value = Unset
				break
			}
		}
	}
	return n
}

// MoveVarTo moves specific variables from the current node to the destination node.
func (n *Node) MoveVarTo(dst *Node, names ...string) *Node {
NAMES:
	for _, name := range names {
		for i, v := range n.vars {
			if v.name == name {
				if v.value.Valid() {
					dst.VarValue(name, v.value)
					n.vars[i].value = Unset
				}
				continue NAMES
			}
		}
	}
	return n
}

// MoveVarPrefixTo moves variables starting with the given prefix to the destination node.
func (n *Node) MoveVarPrefixTo(dst *Node, prefix string) *Node {
	for i, v := range n.vars {
		if v.value.Valid() && strings.HasPrefix(n.vars[i].name, prefix) {
			dst.VarValue(n.vars[i].name, v.value)
			n.vars[i].value = Unset
		}
	}
	return n
}

// MoveVarSuffixTo moves variables ending with the given suffix to the destination node.
func (n *Node) MoveVarSuffixTo(dst *Node, prefix string) *Node {
	for i, v := range n.vars {
		if v.value.Valid() && strings.HasSuffix(n.vars[i].name, prefix) {
			dst.VarValue(n.vars[i].name, v.value)
			n.vars[i].value = Unset
		}
	}
	return n
}

/**/

var staticMap sync.Map

// Static renders the node returned by fn once and caches the result globally.
// Subsequent calls return a cached raw byte node, avoiding re-rendering.
func Static(fn func() *Node) *Node {
	ptr := reflect.ValueOf(fn).Pointer()
	if v, ok := staticMap.Load(ptr); ok {
		return RawBytes(v.([]byte))
	}
	var buf bytes.Buffer
	if err := fn().Render(&buf); err != nil {
		return RawBytes([]byte(err.Error()))
	}
	staticMap.Store(ptr, buf.Bytes())
	return RawBytes(buf.Bytes())
}

// StaticContent sets the content of the node to the cached output of fn.
func (n *Node) StaticContent(fn func() *Node) *Node { return n.Content(Static(fn)) }

/**/

// Content returns a Mod that sets (replaces) the content of the node.
func Content(nodes ...*Node) Mod { return func(n *Node) { n.Content(nodes...) } }

// Content sets (replaces) the content of the node with the provided nodes.
func (n *Node) Content(nodes ...*Node) *Node {
	n.RemoveContent()
	n.Append(nodes...)
	return n
}

// X is an alias for Content.
func (n *Node) X(nodes ...*Node) *Node { return n.Content(nodes...) }

// TextContent returns a Mod that sets (replaces) the content of the node to a single text node.
// The content is HTML-escaped during rendering.
func TextContent(s string) Mod { return func(n *Node) { n.Text(s) } }

// Text creates a text node from a string. The content is HTML-escaped during rendering.
func Text(s string) *Node {
	if s == "" {
		return nil
	}
	n := Get()
	n.tag = "$text"
	n.value = String(s)
	n.writeFn = renderText
	return n
}

// Text sets (replaces) the content of the node to a single text node. The content is HTML-escaped during rendering.
func (n *Node) Text(s string) *Node {
	return n.Content(Text(s))
}

// T is an alias for Text.
func (n *Node) T(s string) *Node { return n.Text(s) }

// TextValue creates a text node from an arbitrary value. The content is HTML-escaped during rendering.
func TextValue(v TypedValue) *Node {
	if !v.Valid() {
		return nil
	}
	n := Get()
	n.tag = "$text"
	n.value = v
	n.writeFn = renderText
	return n
}

// TextValue sets (replaces) the content of the node to a single text node. The content is HTML-escaped during rendering.
func (n *Node) TextValue(v TypedValue) *Node {
	return n.Content(TextValue(v))
}

// TV is an alias for TextValue.
func (n *Node) TV(v TypedValue) *Node { return n.TextValue(v) }

/**/

// Append adds nodes to the end of the content.
func (n *Node) Append(nodes ...*Node) *Node {
	for _, node := range nodes {
		if node != nil { // if at least one is non-nil - append all
			n.content = append(n.content, nodes...)
			return n
		}
	}
	return n
}

// Prepend adds nodes to the beginning of the content.
func (n *Node) Prepend(nodes ...*Node) *Node {
	for _, node := range nodes {
		if node != nil { // if at least one is non-nil - append all
			n.content = slices.Insert(n.content, 0, nodes...)
			return n
		}
	}
	return n
}

// HasContent checks if the node has any content.
func (n *Node) HasContent() bool {
	for _, v := range n.content {
		if v != nil {
			return true
		}
	}
	return false
}

// NoContent checks if the node has no content.
func (n *Node) NoContent() bool {
	for _, v := range n.content {
		if v != nil {
			return false
		}
	}
	return true
}

// RemoveContent clears the content and recursively releases all child nodes.
func (n *Node) RemoveContent() *Node {
	for _, c := range n.content {
		put(c)
	}
	clear(n.content)
	n.content = n.content[:0]
	return n
}

// ExtractContentNodes removes and returns the content of the node.
func (n *Node) ExtractContentNodes() (extracted []*Node) {
	extracted, n.content = n.content, nil
	return
}

// ExtractContent removes and returns the content of the node as Group node.
func (n *Node) ExtractContent() *Node {
	if len(n.content) > 0 {
		g := Group(n.content...)
		clear(n.content)
		n.content = n.content[:0]
		return g
	}
	return nil
}

// MoveContentTo moves all content from the current node to the destination node.
func (n *Node) MoveContentTo(dst *Node) *Node {
	n.content, dst.content = dst.content, n.content
	n.RemoveContent()
	return n
}

// EachContent calls fn for each child node. Iteration stops if fn returns false.
func (n *Node) EachContent(fn func(*Node) bool) *Node {
	for _, node := range n.content {
		if node != nil && !fn(node) {
			return n
		}
	}
	return n
}

// MoveContent is a Mod that moves the content to a destination node.
func (n *Node) MoveContent(dst *Node) {
	n.MoveContentTo(dst)
}

/**/

// HasSlot checks if a named slot exists and has content.
func (n *Node) HasSlot(name string) bool {
	for _, slot := range n.slots {
		if slot.name == name {
			return slot.group != nil
		}
	}
	return false
}

// Slot returns a Mod that sets the content of a named slot.
// If the slot exists, its content is replaced.
func Slot(name string, nodes ...*Node) Mod {
	return func(n *Node) { n.Slot(name, nodes...) }
}

// Slot sets the content of a named slot. If the slot exists, its content is replaced.
func (n *Node) Slot(name string, nodes ...*Node) *Node {
	for _, node := range nodes {
		if node != nil { // if at least one is not nil - add all
			return n.addSlot(name, slotReplace, nodes...)
		}
	}
	return n
}

// AppendSlot adds nodes to the end of a named slot.
func (n *Node) AppendSlot(name string, nodes ...*Node) *Node {
	for _, node := range nodes {
		if node != nil { // if at least one is not nil - add all
			return n.addSlot(name, slotAppend, nodes...)
		}
	}
	return n
}

// PrependSlot adds nodes to the beginning of a named slot.
func (n *Node) PrependSlot(name string, nodes ...*Node) *Node {
	for _, node := range nodes {
		if node != nil { // if at least one is not nil - add all
			return n.addSlot(name, slotPrepend, nodes...)
		}
	}
	return n
}

const (
	slotReplace = iota
	slotAppend
	slotPrepend
)

func (n *Node) addSlot(name string, op byte, nodes ...*Node) *Node {
	for i, slot := range n.slots {
		if slot.name != name {
			continue
		}
		if slot.group != nil {
			switch op {
			case slotReplace:
				slot.group.Content(nodes...)
			case slotAppend:
				slot.group.Append(nodes...)
			case slotPrepend:
				slot.group.Prepend(nodes...)
			}
		} else {
			n.slots[i].group = Group(nodes...)
		}
		return n
	}
	n.slots = append(n.slots, slotNode{
		name:  name,
		group: Group(nodes...),
	})
	return n
}

// DeleteSlot removes specific named slots and releases their content.
func (n *Node) DeleteSlot(names ...string) *Node {
	for _, name := range names {
		for i, slot := range n.slots {
			if slot.name != name {
				continue
			}
			slot.group.Release()
			n.slots[i].group = nil
			break
		}
	}
	return n
}

// ExtractSlotNodes removes and returns the content of a named slot.
func (n *Node) ExtractSlotNodes(name string) []*Node {
	for i, slot := range n.slots {
		if slot.name != name {
			continue
		}
		if slot.group == nil {
			return nil
		}
		extracted := slot.group.ExtractContentNodes()
		slot.group.Release()
		n.slots[i].group = nil
		return extracted
	}
	return nil
}

// Count returns a number of content nodes including nil nodes.
// For exact number of non-nil nodes use CountExact.
func (n *Node) Count() int {
	return len(n.content)
}

// CountExact returns a number of non-nil content nodes.
func (n *Node) CountExact() int {
	count := 0
	for _, x := range n.content {
		if x != nil {
			count++
		}
	}
	return count
}

// CountRecursive returns a total number of content nodes
// in the node itself and its subtree.
// CountRecursive counts all nodes including nil ones.
// For exact number of non-nil nodes use CountRecursiveExact.
func (n *Node) CountRecursive() int {
	count := n.Count()
	for _, x := range n.content {
		count += x.CountRecursive()
	}
	return count
}

// CountExactRecursive returns a total number of non-nil content nodes
// in the node itself and its subtree.
// CountRecursive counts all nodes including nil ones.
// For exact number of non-nil nodes use CountRecursiveExact.
func (n *Node) CountExactRecursive() int {
	count := n.CountExact()
	for _, x := range n.content {
		count += x.CountExactRecursive()
	}
	return count
}

// ExtractSlot removes and returns the contents of a named slot as a Group node.
// If slot does not exist or has no content, ExtractSlot returns nil.
func (n *Node) ExtractSlot(name string) *Node {
	for i, slot := range n.slots {
		if slot.name != name {
			continue
		}
		g := slot.group
		n.slots[i].group = nil
		return g
	}
	return nil
}

// WithSlot executes fn with the slot's content if the specified slot exists and has content.
// Otherwise, fn is not called. Slot contents is removed before passing it to fn.
func (n *Node) WithSlot(name string, fn func(n *Node, slot *Node)) {
	for i, slot := range n.slots {
		if slot.name != name {
			continue
		}
		if slot.group == nil {
			return
		}
		g := slot.group
		n.slots[i].group = nil
		fn(n, g)
		return
	}
}

// MoveSlotTo moves named slots and their content to the destination node.
func (n *Node) MoveSlotTo(dst *Node, names ...string) *Node {
	for _, name := range names {
		for i, slot := range n.slots {
			if slot.name != name {
				continue
			}
			if slot.group == nil {
				continue
			}
			dst.Slot(name, slot.group)
			n.slots[i].group = nil
			break
		}
	}
	return n
}

// MoveSlot returns a Mod that moves named slots to a destination node.
func (n *Node) MoveSlot(names ...string) Mod {
	return func(dst *Node) { n.MoveSlotTo(dst, names...) }
}

/**/

// Postpone adds mods to be applied just before rendering.
func (n *Node) Postpone(mods ...Mod) *Node {
	n.postponed = append(n.postponed, mods...)
	return n
}

/**/

// Own marks the node as owned, preventing it from being returned to the pool by Release.
func (n *Node) Own() *Node {
	n.flag |= flagOwned
	return n
}

// Disown removes the owned mark, allowing the node to be returned to the pool.
// func (n *Node) Disown() *Node {
// 	n.flag &^= flagOwned
// 	return n
// }

// Owned returns true if the node has been marked as owned and will not be returned to the pool.
func (n *Node) Owned() bool { return n.flag&flagOwned != 0 }

// UnsafeScript enables unsafe <script> rendering.
func (n *Node) UnsafeScript() { n.flag |= flagScript }

// Release returns the node and its children to the pool for reuse.
// If the node is marked as Owned, neither it nor its subtree will be returned to the pool.
func (n *Node) Release() { put(n) }

// SetPoolingNeighbor links another node to be released together with n.
func (n *Node) SetPoolingNeighbor(x *Node) { n.attached = append(n.attached, x) }

// String renders the node to a string.
func (n *Node) String() string {
	b := new(strings.Builder)
	if err := n.Render(b); err != nil {
		return err.Error()
	}
	return b.String()
}

/**/

// WriteFn returns a Mod that overrides the default rendering logic with fn.
func WriteFn(fn func(*Node, io.Writer) error) Mod { return func(n *Node) { n.SetWriteFn(fn) } }

// SetWriteFn overrides the default rendering logic of the node with fn.
func (n *Node) SetWriteFn(fn func(*Node, io.Writer) error) *Node {
	n.writeFn = fn
	return n
}

/**/

func renderGroup(n *Node, w io.Writer) error {
	for _, node := range n.content {
		if node != nil {
			if err := node.Render(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func renderRaw(n *Node, w io.Writer) (err error) {
	switch n.value.any.(type) {
	case stringptr:
		s := unsafe.Slice(n.value.any.(stringptr), n.value.num)
		_, err = w.Write(s)
	case byteptr:
		s := unsafe.Slice(n.value.any.(byteptr), n.value.num)
		_, err = w.Write(s)
	}
	return
}

var (
	trueValue  = []byte("true")
	falseValue = []byte("false")
)

func renderText(n *Node, w io.Writer) (err error) {
	switch n.value.Kind() {

	case KindNone:
		return

	case KindAny:
		_, err = fmt.Fprint(EscapeWriter(w.Write), n.value.any)

	case KindBool:
		if n.value.num == 1 {
			_, err = w.Write(trueValue)
		} else {
			_, err = w.Write(falseValue)
		}

	case KindInt64:
		err = WriteInt(w, int64(n.value.num))

	case KindUint64:
		err = WriteUint(w, n.value.num)

	case KindFloat64:
		err = WriteFloat(w, math.Float64frombits(n.value.num))

	case KindString:
		s := unsafe.Slice(n.value.any.(stringptr), n.value.num)
		_, err = EscapeWriter(w.Write).Write(s)

	case KindBytes:
		s := unsafe.Slice(n.value.any.(byteptr), n.value.num)
		_, err = EscapeWriter(w.Write).Write(s)

	case KindJSON:
		buf := jsonBufPool.Get().(*bytes.Buffer)
		buf.Reset()
		if err = json.NewEncoder(buf).Encode(n.value.any); err != nil {
			jsonBufPool.Put(buf)
			return err
		}
		b := buf.Bytes()
		if len(b) > 0 && b[len(b)-1] == '\n' {
			b = b[:len(b)-1]
		}
		_, err = EscapeWriter(w.Write).Write(b)
		jsonBufPool.Put(buf)

	default:
		_, err = fmt.Fprint(EscapeWriter(w.Write), n.value.any)

	}
	return
}

/**/

var jsonBufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

/**/

// Render writes the HTML representation of the node to w.
func (n *Node) Render(w io.Writer) error {
	if n == nil {
		return nil
	}
	if n.writeFn != nil {
		return n.writeFn(n, w)
	}
	for _, fn := range n.postponed {
		fn(n)
	}
	if !ValidTag(n.tag) {
		return fmt.Errorf("invalid tag: %v", n.tag)
	}

	if _, err := w.Write(openingPart); err != nil {
		return err
	}
	if _, err := WriteString(w, n.tag); err != nil {
		return err
	}
	if len(n.class.o) > 0 {
		if err := writeClass(w, n.class.o); err != nil {
			return err
		}
	}
	if len(n.attrs.o) > 0 {
		if err := writeAttributes(w, n.attrs.o); err != nil {
			return err
		}
	}
	if n.flag&flagVoid != 0 {
		_, err := w.Write(closingPartVoid)
		return err
	}
	if _, err := w.Write(closingPartRight); err != nil {
		return err
	}
	if len(n.content) > 0 {
		if (n.flag&flagScript == 0) && isScriptTag(n.tag) {
			return fmt.Errorf("script tags are not allowed to have content, use UnsafeScript to bypass this error")
		}
		for _, content := range n.content {
			if err := content.Render(w); err != nil {
				return err
			}
		}
	}
	if _, err := w.Write(closingPartLeft); err != nil {
		return err
	}
	if _, err := WriteString(w, n.tag); err != nil {
		return err
	}
	if _, err := w.Write(closingPartRight); err != nil {
		return err
	}
	return nil
}

func isScriptTag(tag string) bool {
	if len(tag) != 6 {
		return false
	}
	return (tag[0]|0x20) == 's' && (tag[1]|0x20) == 'c' && (tag[2]|0x20) == 'r' &&
		(tag[3]|0x20) == 'i' && (tag[4]|0x20) == 'p' && (tag[5]|0x20) == 't'
}

var (
	space            = []byte(" ")
	quote            = []byte(`"`)
	equals           = []byte("=")
	classPrefix      = []byte(` class="`)
	openingPart      = []byte("<")
	closingPartVoid  = []byte("/>")
	closingPartLeft  = []byte("</")
	closingPartRight = []byte(">")
)

func writeClass(w io.Writer, classes []classEntry) error {
	if _, err := w.Write(classPrefix); err != nil {
		return err
	}

	first := true
OUTER:
	for i := 0; i < len(classes); i++ {
		c := classes[i]
		if !c.active || !ValidClass(c.name) {
			continue
		}
		for k := i + 1; k < len(classes); k++ {
			if classes[k].name == c.name {
				continue OUTER
			}
		}
		if !first {
			if _, err := w.Write(space); err != nil {
				return err
			}
		}
		if _, err := WriteString(w, c.name); err != nil {
			return err
		}
		first = false
	}

	if _, err := w.Write(quote); err != nil {
		return err
	}
	return nil
}

func writeAttributes(w io.Writer, attrs []valueEntry) error {
OUTER:
	for i := 0; i < len(attrs); i++ {
		a := attrs[i]
		if !a.value.Valid() {
			continue
		}
		if !ValidAttr(a.name) {
			continue
		}
		for k := i + 1; k < len(attrs); k++ {
			if attrs[k].name == a.name {
				continue OUTER
			}
		}
		kind := a.value.Kind()

		if kind == KindBool && a.value.num == 0 {
			continue
		}

		if _, err := w.Write(space); err != nil {
			return err
		}
		if _, err := WriteString(w, a.name); err != nil {
			return err
		}

		if kind == KindBool {
			continue
		}

		if _, err := w.Write(equals); err != nil {
			return err
		}
		if _, err := w.Write(quote); err != nil {
			return err
		}

		switch kind {
		case KindString:
			ptr := (*byte)(a.value.any.(stringptr))
			s := unsafe.Slice(ptr, a.value.num)
			if _, err := EscapeWriter(w.Write).Write(s); err != nil {
				return err
			}
		case KindBytes:
			s := unsafe.Slice(a.value.any.(byteptr), a.value.num)
			if _, err := EscapeWriter(w.Write).Write(s); err != nil {
				return err
			}
		case KindInt64:
			if err := WriteInt(w, int64(a.value.num)); err != nil {
				return err
			}
		case KindUint64:
			if err := WriteUint(w, a.value.num); err != nil {
				return err
			}
		case KindFloat64:
			if err := WriteFloat(w, math.Float64frombits(a.value.num)); err != nil {
				return err
			}
		case KindJSON:
			buf := jsonBufPool.Get().(*bytes.Buffer)
			buf.Reset()
			if err := json.NewEncoder(buf).Encode(a.value.any); err != nil {
				jsonBufPool.Put(buf)
				return err
			}
			b := buf.Bytes()
			if len(b) > 0 && b[len(b)-1] == '\n' {
				b = b[:len(b)-1]
			}
			_, err := EscapeWriter(w.Write).Write(b)
			jsonBufPool.Put(buf)
			if err != nil {
				return err
			}
		default:
			if _, err := fmt.Fprint(EscapeWriter(w.Write), a.value.any); err != nil {
				return err
			}
		}
		if _, err := w.Write(quote); err != nil {
			return err
		}
	}
	return nil
}

// EscapeWriter is an io.Writer that escapes HTML special characters.
type EscapeWriter func(p []byte) (n int, err error)

var (
	escapedDoubleQuote = []byte("&#34;")
	escapedSingleQuote = []byte("&#39;")
	escapedLt          = []byte("&lt;")
	escapedGt          = []byte("&gt;")
	escapedAmp         = []byte("&amp;")
)

func (w EscapeWriter) Write(p []byte) (int, error) {
	start := 0
	n := 0
	for i, c := range p {
		var esc []byte
		switch c {
		case '"':
			esc = escapedDoubleQuote
		case '\'':
			esc = escapedSingleQuote
		case '<':
			esc = escapedLt
		case '>':
			esc = escapedGt
		case '&':
			esc = escapedAmp
		}
		if esc != nil {
			if start < i {
				x, err := w(p[start:i])
				n += x
				if err != nil {
					return n, err
				}
			}
			x, err := w(esc)
			n += x
			if err != nil {
				return n, err
			}
			start = i + 1
		}
	}
	if start < len(p) {
		x, err := w(p[start:])
		n += x
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func WriteString(w io.Writer, s string) (int, error) {
	return w.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

/**/

//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func WriteInt(w io.Writer, n int64) error {
	var buf [20]byte // int64 max 20 chars + sign
	i := len(buf)

	var u uint64
	if n < 0 {
		u = uint64(^n + 1) // two's complement inversion for absolute value
	} else {
		u = uint64(n)
	}
	if u == 0 {
		i--
		buf[i] = '0'
	} else {
		for u > 0 {
			i--
			buf[i] = byte('0' + u%10)
			u /= 10
		}
	}
	if n < 0 {
		i--
		buf[i] = '-'
	}
	slice := buf[i:]
	safeSlice := *(*[]byte)(noescape(unsafe.Pointer(&slice)))
	_, err := w.Write(safeSlice)
	return err
}

func WriteUint(w io.Writer, n uint64) error {
	var buf [20]byte // uint64 max 20 chars
	i := len(buf)
	if n == 0 {
		i--
		buf[i] = '0'
	} else {
		for n > 0 {
			i--
			buf[i] = byte('0' + n%10)
			n /= 10
		}
	}
	slice := buf[i:]
	safeSlice := *(*[]byte)(noescape(unsafe.Pointer(&slice)))
	_, err := w.Write(safeSlice)
	return err
}

func WriteFloat(w io.Writer, f float64) error {
	var buf [64]byte
	b := strconv.AppendFloat(buf[:0], f, 'g', -1, 64)
	safeSlice := *(*[]byte)(noescape(unsafe.Pointer(&b)))
	_, err := w.Write(safeSlice)
	return err
}

/**/

// NoPool disables automatic pooling.
var NoPool bool

var nodePool = sync.Pool{
	New: func() any { return newNode() },
}

func newNode() *Node {
	n := &Node{
		tag:   "div",
		attrs: newAttrMap(),
		class: newClassMap(),
	}

	if NoPool {
		return n
	}

	n.content = make([]*Node, 0, 16)
	n.slots = make([]slotNode, 0, 4)
	n.vars = make([]valueEntry, 0, 4)

	return n
}

// Get retrieves a zeroed Node from the global pool.
// Most users should use Build, Group, Text and others instead.
func Get() *Node {
	if NoPool {
		return newNode()
	}
	n := nodePool.Get().(*Node)

	if !n.acquired.CompareAndSwap(false, true) {
		panic("htm: got already acquired node; pool is corrupted")
	}

	return n
}

// var (
// 	contentSliceThrown     atomic.Uint64
// 	slotContentSliceThrown atomic.Uint64
// 	slotSliceThrown        atomic.Uint64
// 	varSliceThrown         atomic.Uint64
// 	attachedSliceThrown    atomic.Uint64
// )

func put(n *Node) {
	if n == nil {
		return
	}
	if NoPool {
		return
	}
	if n.flag&flagOwned != 0 {
		return
	}
	if !n.acquired.CompareAndSwap(true, false) {
		panic("htm: attempt to release an already released node")
	}

	n.tag = "div"
	n.flag = 0
	n.value = TypedValue{}

	n.attrs.reset()
	n.class.reset()

	n.writeFn = nil

	// if n.content == nil {
	// 	// n.content = make([]*Node, 0, 16)
	// } else
	if len(n.content) > 0 {
		for _, node := range n.content {
			put(node)
		}
		clear(n.content)
		n.content = n.content[:0]
		// if cap(n.content) > 128 {
		// 	contentSliceThrown.Add(1)
		// 	n.content = make([]*Node, 0, 64)
		// }
	}

	if len(n.slots) > 0 {
		for _, slot := range n.slots {
			slot.group.Release()
			// if cap(slot.content) > 128 {
			// 	slotContentSliceThrown.Add(1)
			// 	slot.content = make([]*Node, 0, 64)
			// }
		}
		n.slots = n.slots[:0]
		// if cap(n.slots) > 16 {
		// 	slotSliceThrown.Add(1)
		// 	n.slots = make([]slotNode, 0, 8)
		// }
	}

	if len(n.vars) > 0 {
		clear(n.vars)
		n.vars = n.vars[:0]
		// if cap(n.vars) > 32 {
		// 	varSliceThrown.Add(1)
		// 	n.vars = make([]valueEntry, 0, 16)
		// }
	}

	if len(n.attached) > 0 {
		for _, attached := range n.attached {
			put(attached)
		}
		clear(n.attached)
		n.attached = n.attached[:0]
		// if cap(n.attached) > 64 {
		// 	attachedSliceThrown.Add(1)
		// 	n.attached = make([]*Node, 0, 32)
		// }
	}

	if len(n.postponed) > 0 {
		clear(n.postponed)
		n.postponed = n.postponed[:0]
	}

	nodePool.Put(n)
}

/**/

type attrMap struct {
	o []valueEntry
}

func newAttrMap() *attrMap {
	return &attrMap{
		o: make([]valueEntry, 0, 8),
	}
}

func (am *attrMap) reset() {
	clear(am.o)
	am.o = am.o[:0]
}

func (am *attrMap) get(name string) (TypedValue, bool) {
	for i := len(am.o) - 1; i >= 0; i-- {
		e := am.o[i]
		if e.name == name {
			if e.value.Valid() {
				return e.value, true
			}
			return Unset, false
		}
	}
	return Unset, false
}

func (am *attrMap) findActiveIndex(name string) int {
	for i := len(am.o) - 1; i >= 0; i-- {
		e := am.o[i]
		if e.name == name {
			if e.value.Valid() {
				return i
			}
			return -1
		}
	}
	return -1
}

func (am *attrMap) hasAny(names ...string) bool {
	for _, name := range names {
		for i := len(am.o) - 1; i >= 0; i-- {
			if am.o[i].name == name {
				if am.o[i].value.Valid() {
					return true
				}
				break
			}
		}
	}
	return false
}

func (am *attrMap) hasAll(names ...string) bool {
	for _, name := range names {
		ok := false
		for i := len(am.o) - 1; i >= 0; i-- {
			if am.o[i].name == name {
				ok = am.o[i].value.Valid()
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func (am *attrMap) hasPrefix(prefix string) bool {
OUTER:
	for i := len(am.o) - 1; i >= 0; i-- {
		e := am.o[i]
		if !e.value.Valid() {
			continue
		}
		if !strings.HasPrefix(e.name, prefix) {
			continue
		}
		for k := i + 1; k < len(am.o); k++ {
			if am.o[k].name == e.name {
				continue OUTER
			}
		}
		return true
	}
	return false
}

func (am *attrMap) hasSuffix(suffix string) bool {
OUTER:
	for i := len(am.o) - 1; i >= 0; i-- {
		e := am.o[i]
		if !e.value.Valid() {
			continue
		}
		if !strings.HasSuffix(e.name, suffix) {
			continue
		}
		for k := i + 1; k < len(am.o); k++ {
			if am.o[k].name == e.name {
				continue OUTER
			}
		}
		return true
	}
	return false
}

func (am *attrMap) each(fn func(string, TypedValue) bool) {
OUTER:
	for i := len(am.o) - 1; i >= 0; i-- {
		e := am.o[i]
		if !e.value.Valid() {
			continue
		}
		for k := i + 1; k < len(am.o); k++ {
			if am.o[k].name == e.name {
				continue OUTER
			}
		}
		if !fn(e.name, e.value) {
			return
		}
	}
}

func (am *attrMap) set(name string, v TypedValue) {
	if name == "" {
		return
	}
	am.o = append(am.o, valueEntry{name: name, value: v})
}

func (am *attrMap) extract(name string) (TypedValue, bool) {
	for i := len(am.o) - 1; i >= 0; i-- {
		if am.o[i].name != name {
			continue
		}
		v := am.o[i].value
		if !v.Valid() {
			return Unset, false
		}
		am.o[i].value = Unset
		return v, true
	}
	return Unset, false
}

func (am *attrMap) movePrefixTo(dst *attrMap, prefix string) {
OUTER:
	for i := len(am.o) - 1; i >= 0; i-- {
		e := am.o[i]
		v := e.value
		if !v.Valid() {
			continue
		}
		if !strings.HasPrefix(e.name, prefix) {
			continue
		}
		for k := i + 1; k < len(am.o); k++ {
			if am.o[k].name == e.name {
				continue OUTER
			}
		}
		dst.set(e.name, v)
		am.o[i].value = Unset
	}
}

func (am *attrMap) moveSuffixTo(dst *attrMap, suffix string) {
OUTER:
	for i := len(am.o) - 1; i >= 0; i-- {
		e := am.o[i]
		v := e.value
		if !v.Valid() {
			continue
		}
		if !strings.HasSuffix(e.name, suffix) {
			continue
		}
		for k := i + 1; k < len(am.o); k++ {
			if am.o[k].name == e.name {
				continue OUTER
			}
		}
		dst.set(e.name, v)
		am.o[i].value = Unset
	}
}

/**/

type (
	classMap struct {
		o []classEntry
	}
	classEntry struct {
		name   string
		active bool
	}
)

func newClassMap() *classMap {
	return &classMap{
		o: make([]classEntry, 0, 16),
	}
}

func (cm *classMap) reset() {
	cm.o = cm.o[:0]
}

func (cm *classMap) has(name string) bool {
	for i := len(cm.o) - 1; i >= 0; i-- {
		if cm.o[i].name == name {
			return cm.o[i].active
		}
	}
	return false
}

func (cm *classMap) extract(name string) bool {
	for i := len(cm.o) - 1; i >= 0; i-- {
		if cm.o[i].name == name {
			out := cm.o[i].active
			cm.o[i].active = false
			return out
		}
	}
	return false
}

func (cm *classMap) hasAny(names ...string) bool {
	for _, name := range names {
		for i := len(cm.o) - 1; i >= 0; i-- {
			if cm.o[i].name == name {
				if cm.o[i].active {
					return true
				}
				break
			}
		}
	}
	return false
}

func (cm *classMap) hasAll(names ...string) bool {
	for _, name := range names {
		ok := false
		for i := len(cm.o) - 1; i >= 0; i-- {
			if cm.o[i].name == name {
				ok = cm.o[i].active
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func (cm *classMap) hasPrefix(prefix string) bool {
OUTER:
	for i := len(cm.o) - 1; i >= 0; i-- {
		e := cm.o[i]
		if !e.active {
			continue
		}
		if !strings.HasPrefix(e.name, prefix) {
			continue
		}
		for k := i + 1; k < len(cm.o); k++ {
			if cm.o[k].name == e.name {
				continue OUTER
			}
		}
		if e.active {
			return true
		}
	}
	return false
}

func (cm *classMap) hasSuffix(suffix string) bool {
OUTER:
	for i := len(cm.o) - 1; i >= 0; i-- {
		e := cm.o[i]
		if !e.active {
			continue
		}
		if !strings.HasSuffix(e.name, suffix) {
			continue
		}
		for k := i + 1; k < len(cm.o); k++ {
			if cm.o[k].name == e.name {
				continue OUTER
			}
		}
		return true
	}
	return false
}

func (cm *classMap) each(fn func(string) bool) {
OUTER:
	for i := len(cm.o) - 1; i >= 0; i-- {
		e := cm.o[i]
		if !e.active {
			continue
		}
		for k := i + 1; k < len(cm.o); k++ {
			if cm.o[k].name == e.name {
				continue OUTER
			}
		}
		if !fn(e.name) {
			return
		}
	}
}

func (cm *classMap) setMulti(s string, active bool) {
	start := -1
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\n' && s[i] != '\t' && s[i] != '\r' && s[i] != '\f' {
			if start == -1 {
				start = i
			}
		} else if start != -1 {
			cm.o = append(cm.o, classEntry{
				name:   s[start:i],
				active: active,
			})
			start = -1
		}
	}
	if start != -1 {
		cm.o = append(cm.o, classEntry{
			name:   s[start:],
			active: active,
		})
	}
}

func (cm *classMap) setOne(name string, active bool) {
	cm.o = append(cm.o, classEntry{
		name:   name,
		active: active,
	})
}

/**/

// ValidTag checks if the string is a valid HTML tag name.
func ValidTag(tag string) bool {
	if tag == "" {
		return false
	}
	for _, c := range tag {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-' {
			continue
		}
		return false
	}
	if (tag[0] >= 'a' && tag[0] <= 'z') || (tag[0] >= 'A' && tag[0] <= 'Z') {
		return true
	}
	return false
}

// ValidAttr checks if the string is a valid HTML attribute name.
func ValidAttr(attr string) bool {
	if attr == "" {
		return false
	}
	for _, c := range attr {
		if c == '"' || c == '<' || c == '>' || c == '#' || c == ' ' || c == '&' || c == '\'' || c == '\\' || c == '=' {
			return false
		}
	}
	if (attr[0] >= '0' && attr[0] <= '9') || attr[0] == '-' {
		return false
	}
	return true
}

// ValidClass checks if the string is a valid CSS class name.
func ValidClass(class string) bool {
	return class != "" && !strings.ContainsAny(class, `"'`) // <>&`)
}

/**/

type (
	stringptr *byte
	byteptr   *byte
)

type ValueKind int

const (
	KindNone ValueKind = iota
	KindAny
	KindBool
	KindInt64
	KindUint64
	KindFloat64
	KindString
	KindJSON
	KindBytes
)

// Unset represents an empty, unset or removed value.
var Unset TypedValue

// TypedValue can hold strings, ints, booleans and floats without allocation.
type TypedValue struct {
	_   [0]func() // idea is taken from slog to disallow equality checks (==)
	num uint64    // bool, int*, uint*, float*, string/bytes len, kind for JSON
	any any       // ValueKind, stringptr, byteptr, any
}

func (v TypedValue) Kind() ValueKind {
	if v.any == nil {
		return KindNone
	}
	switch k := v.any.(type) {
	case ValueKind:
		return k
	case stringptr:
		return KindString
	case byteptr:
		return KindBytes
	default:
		return ValueKind(v.num)
	}
}

// Valid returns true if it contains a valid value.
func (v TypedValue) Valid() bool {
	return v.any != nil
}

func String(s string) TypedValue {
	return TypedValue{num: uint64(len(s)), any: stringptr(unsafe.StringData(s))}
}
func Bytes(b []byte) TypedValue {
	return TypedValue{num: uint64(len(b)), any: byteptr(unsafe.SliceData(b))}
}

func Bool(v bool) TypedValue {
	var b uint64
	if v {
		b = 1
	}
	return TypedValue{num: b, any: KindBool}
}

func Int(v int) TypedValue       { return Int64(int64(v)) }
func Uint(v uint) TypedValue     { return Uint64(uint64(v)) }
func Int64(v int64) TypedValue   { return TypedValue{num: uint64(v), any: KindInt64} }
func Uint64(v uint64) TypedValue { return TypedValue{num: v, any: KindUint64} }
func Float(v float64) TypedValue { return TypedValue{num: math.Float64bits(v), any: KindFloat64} }

func Any(v any) TypedValue  { return TypedValue{any: v, num: uint64(KindAny)} }  // inverse
func JSON(v any) TypedValue { return TypedValue{any: v, num: uint64(KindJSON)} } // inverse

func (v TypedValue) String() (string, bool) {
	if sp, ok := v.any.(stringptr); ok {
		return unsafe.String(sp, v.num), true
	}
	return "", false
}

func (v TypedValue) StringOrDefault(d string) string {
	if x, ok := v.String(); ok {
		return x
	}
	return d
}
func (v TypedValue) StringOrZero() string {
	if x, ok := v.String(); ok {
		return x
	}
	return ""
}

func (v TypedValue) Bytes() ([]byte, bool) {
	if bp, ok := v.any.(byteptr); ok {
		return unsafe.Slice(bp, v.num), true
	}
	return nil, false
}

func (v TypedValue) BytesOrZero() []byte {
	if b, ok := v.Bytes(); ok {
		return b
	}
	return nil
}

func (v TypedValue) Int64() (int64, bool) {
	if k, ok := v.any.(ValueKind); ok && k == KindInt64 {
		return int64(v.num), true
	}
	return 0, false
}

func (v TypedValue) IntOrDefault(d int64) int64 {
	if x, ok := v.Int64(); ok {
		return x
	}
	return d
}

func (v TypedValue) IntOrZero() int64 {
	if x, ok := v.Int64(); ok {
		return x
	}
	return 0
}

func (v TypedValue) Uint64() (uint64, bool) {
	if k, ok := v.any.(ValueKind); ok && k == KindUint64 {
		return v.num, true
	}
	return 0, false
}

func (v TypedValue) UintOrDefault(d uint64) uint64 {
	if x, ok := v.Uint64(); ok {
		return x
	}
	return d
}
func (v TypedValue) UintOrZero() uint64 {
	if x, ok := v.Uint64(); ok {
		return x
	}
	return 0
}

func (v TypedValue) Float64() (float64, bool) {
	if k, ok := v.any.(ValueKind); ok && k == KindFloat64 {
		return math.Float64frombits(v.num), true
	}
	return 0, false
}

func (v TypedValue) FloatOrDefault(d float64) float64 {
	if x, ok := v.Float64(); ok {
		return x
	}
	return d
}

func (v TypedValue) FloatOrZero() float64 {
	if x, ok := v.Float64(); ok {
		return x
	}
	return 0
}

func (v TypedValue) Bool() (bool, bool) {
	if k, ok := v.any.(ValueKind); ok && k == KindBool {
		return v.num == 1, true
	}
	return false, false
}

func (v TypedValue) BoolOrDefault(d bool) bool {
	if x, ok := v.Bool(); ok {
		return x
	}
	return d
}

func (v TypedValue) BoolOrZero() bool {
	if x, ok := v.Bool(); ok {
		return x
	}
	return false
}

func (v TypedValue) JSON() (any, bool) {
	if v.Kind() == KindJSON {
		return v.any, true
	}
	return nil, false
}

func (v TypedValue) JSONOrZero() any {
	if v.Kind() == KindJSON {
		return v.any
	}
	return nil
}

func (v TypedValue) Any() any {
	switch k := v.any.(type) {
	case stringptr:
		return unsafe.String(k, v.num)
	case ValueKind:
		switch k {
		case KindInt64:
			return int64(v.num)
		case KindUint64:
			return v.num
		case KindFloat64:
			return math.Float64frombits(v.num)
		case KindBool:
			return v.num == 1
		}
	}
	return v.any
}
