package htm

import (
	"strconv"
	"sync/atomic"
	"time"
)

// globals

func AccessKey(v string) Mod             { return Attr("accesskey", v) }
func (n *Node) AccessKey(v string) *Node { return n.Attr("accesskey", v) }

func Autocapitalize(v string) Mod             { return Attr("autocapitalize", v) }
func (n *Node) Autocapitalize(v string) *Node { return n.Attr("autocapitalize", v) }

func Autocomplete(v string) Mod             { return Attr("autocomplete", v) }
func (n *Node) Autocomplete(v string) *Node { return n.Attr("autocomplete", v) }

func ContentAttr(v string) Mod             { return Attr("content", v) }
func (n *Node) ContentAttr(v string) *Node { return n.Attr("content", v) }

func Dir(v string) Mod             { return Attr("dir", v) }
func (n *Node) Dir(v string) *Node { return n.Attr("dir", v) }

// Draggable sets the "draggable".
// If value is omitted, it sets a string value "true".
// If value is provided and false, it sets a string value "false".
func Draggable(v ...bool) Mod {
	if len(v) > 0 && v[0] == false {
		return Attr("draggable", "false")
	}
	return Attr("draggable", "true")
}

// Draggable sets the "draggable".
// If value is omitted, it sets a string value "true".
// If value is provided and false, it sets a string value "false".
func (n *Node) Draggable(v ...bool) *Node {
	if len(v) > 0 && v[0] == false {
		return n.Attr("draggable", "false")
	}
	return n.Attr("draggable", "true")
}

func EnterKeyHint(v string) Mod             { return Attr("enterkeyhint", v) }
func (n *Node) EnterKeyHint(v string) *Node { return n.Attr("enterkeyhint", v) }

func InputMode(v string) Mod             { return Attr("inputmode", v) }
func (n *Node) InputMode(v string) *Node { return n.Attr("inputmode", v) }

func Lang(v string) Mod             { return Attr("lang", v) }
func (n *Node) Lang(v string) *Node { return n.Attr("lang", v) }

// Spellcheck sets the "spellcheck" attribute.
// If value is omitted, it sets a boolean attribute.
// If value is provided and false, it sets a string value "false".
func Spellcheck(v ...bool) Mod {
	if len(v) > 0 && v[0] == false {
		return Attr("spellcheck", "false")
	}
	return AttrValue("spellcheck")
}

// Spellcheck sets the "spellcheck" attribute.
// If value is omitted, it sets a boolean attribute.
// If value is provided and false, it sets a string value "false".
func (n *Node) Spellcheck(v ...bool) *Node {
	if len(v) > 0 && v[0] == false {
		return n.Attr("spellcheck", "false")
	}
	return n.Attr("spellcheck")
}

func Style(v string) Mod             { return Attr("style", v) }
func (n *Node) Style(v string) *Node { return n.Attr("style", v) }

func TabIndex(v int) Mod             { return AttrValue("tabindex", Int(v)) }
func (n *Node) TabIndex(v int) *Node { return n.AttrValue("tabindex", Int(v)) }

func Hint(v string) Mod             { return Attr("title", v) }
func (n *Node) Hint(v string) *Node { return n.Attr("title", v) }

func Translate(v string) Mod             { return Attr("translate", v) }
func (n *Node) Translate(v string) *Node { return n.Attr("translate", v) }

// forms

func Accept(v string) Mod             { return Attr("accept", v) }
func (n *Node) Accept(v string) *Node { return n.Attr("accept", v) }

func AcceptCharset(v string) Mod             { return Attr("accept-charset", v) }
func (n *Node) AcceptCharset(v string) *Node { return n.Attr("accept-charset", v) }

func Action(v string) Mod             { return Attr("action", v) }
func (n *Node) Action(v string) *Node { return n.Attr("action", v) }

func Allow(v string) Mod             { return Attr("allow", v) }
func (n *Node) Allow(v string) *Node { return n.Attr("allow", v) }

func Alt(v string) Mod             { return Attr("alt", v) }
func (n *Node) Alt(v string) *Node { return n.Attr("alt", v) }

func As(v string) Mod             { return Attr("as", v) }
func (n *Node) As(v string) *Node { return n.Attr("as", v) }

func Capture(v string) Mod             { return Attr("capture", v) }
func (n *Node) Capture(v string) *Node { return n.Attr("capture", v) }

func Charset(v string) Mod             { return Attr("charset", v) }
func (n *Node) Charset(v string) *Node { return n.Attr("charset", v) }

func CiteAttr(v string) Mod         { return Attr("cite", v) }
func (n *Node) Cite(v string) *Node { return n.Attr("cite", v) }

func CrossOrigin(v string) Mod             { return Attr("crossorigin", v) }
func (n *Node) CrossOrigin(v string) *Node { return n.Attr("crossorigin", v) }

func For(v string) Mod             { return Attr("for", v) }
func (n *Node) For(v string) *Node { return n.Attr("for", v) }

func Href(v string) Mod             { return Attr("href", v) }
func (n *Node) Href(v string) *Node { return n.Attr("href", v) }

func Hreflang(v string) Mod             { return Attr("hreflang", v) }
func (n *Node) Hreflang(v string) *Node { return n.Attr("hreflang", v) }

func HttpEquiv(v string) Mod             { return Attr("http-equiv", v) }
func (n *Node) HttpEquiv(v string) *Node { return n.Attr("http-equiv", v) }

func Kind(v string) Mod             { return Attr("kind", v) }
func (n *Node) Kind(v string) *Node { return n.Attr("kind", v) }

func GroupLabel(v string) Mod             { return Attr("label", v) }
func (n *Node) GroupLabel(v string) *Node { return n.Attr("label", v) }

func List(v string) Mod             { return Attr("list", v) }
func (n *Node) List(v string) *Node { return n.Attr("list", v) }

// max and min can be a number or a date

func Max(v int) Mod             { return AttrValue("max", Int(v)) }
func (n *Node) Max(v int) *Node { return n.AttrValue("max", Int(v)) }

func Min(v int) Mod             { return AttrValue("min", Int(v)) }
func (n *Node) Min(v int) *Node { return n.AttrValue("min", Int(v)) }

func MaxDate(v time.Time) Mod             { return Attr("max", v.Format("2006-01-02")) }
func (n *Node) MaxDate(v time.Time) *Node { return n.Attr("max", v.Format("2006-01-02")) }

func MinDate(v time.Time) Mod             { return Attr("min", v.Format("2006-01-02")) }
func (n *Node) MinDate(v time.Time) *Node { return n.Attr("min", v.Format("2006-01-02")) }

func Name(v string) Mod             { return Attr("name", v) }
func (n *Node) Name(v string) *Node { return n.Attr("name", v) }

func Optimum(v float64) Mod             { return AttrValue("optimum", Float(v)) }
func (n *Node) Optimum(v float64) *Node { return n.AttrValue("optimum", Float(v)) }

func Pattern(v string) Mod             { return Attr("pattern", v) }
func (n *Node) Pattern(v string) *Node { return n.Attr("pattern", v) }

func Placeholder(v string) Mod             { return Attr("placeholder", v) }
func (n *Node) Placeholder(v string) *Node { return n.Attr("placeholder", v) }

func Poster(v string) Mod             { return Attr("poster", v) }
func (n *Node) Poster(v string) *Node { return n.Attr("poster", v) }

func Preload(v string) Mod             { return Attr("preload", v) }
func (n *Node) Preload(v string) *Node { return n.Attr("preload", v) }

func Rel(v string) Mod             { return Attr("rel", v) }
func (n *Node) Rel(v string) *Node { return n.Attr("rel", v) }

func ReferrerPolicy(v string) Mod             { return Attr("referrerpolicy", v) }
func (n *Node) ReferrerPolicy(v string) *Node { return n.Attr("referrerpolicy", v) }

func Role(v string) Mod             { return Attr("role", v) }
func (n *Node) Role(v string) *Node { return n.Attr("role", v) }

func Sizes(v string) Mod             { return Attr("sizes", v) }
func (n *Node) Sizes(v string) *Node { return n.Attr("sizes", v) }

func SlotAttr(v string) Mod             { return Attr("slot", v) }
func (n *Node) SlotAttr(v string) *Node { return n.Attr("slot", v) }

func Src(v string) Mod             { return Attr("src", v) }
func (n *Node) Src(v string) *Node { return n.Attr("src", v) }

func Srcdoc(v string) Mod             { return Attr("srcdoc", v) }
func (n *Node) Srcdoc(v string) *Node { return n.Attr("srcdoc", v) }

func Srclang(v string) Mod             { return Attr("srclang", v) }
func (n *Node) Srclang(v string) *Node { return n.Attr("srclang", v) }

func Srcset(v string) Mod             { return Attr("srcset", v) }
func (n *Node) Srcset(v string) *Node { return n.Attr("srcset", v) }

// Step adds a "step" attribute. If value is omitted, "any" is used.
func Step(value ...int) Mod {
	if len(value) > 0 {
		v := value[0]
		return AttrValue("step", Int(v))
	}
	return Attr("step", "any")
}

// Step adds a "step" attribute. If value is omitted, "any" is used.
func (n *Node) Step(value ...int) *Node {
	if len(value) > 0 {
		v := value[0]
		return n.AttrValue("step", Int(v))
	}
	return n.Attr("step", "any")
}

func Type(v string) Mod             { return Attr("type", v) }
func (n *Node) Type(v string) *Node { return n.Attr("type", v) }

// value is ambiguous (string, int, bool, date), so TypedValue

func Value(v TypedValue) Mod             { return AttrValue("value", v) }
func (n *Node) Value(v TypedValue) *Node { return n.AttrValue("value", v) }

func Width(v int) Mod             { return AttrValue("width", Int(v)) }
func (n *Node) Width(v int) *Node { return n.AttrValue("width", Int(v)) }

func Height(v int) Mod             { return AttrValue("height", Int(v)) }
func (n *Node) Height(v int) *Node { return n.AttrValue("height", Int(v)) }

/**/

func AllowFullscreen(v ...TypedValue) Mod             { return AttrValue("allowfullscreen", v...) }
func (n *Node) AllowFullscreen(v ...TypedValue) *Node { return n.AttrValue("allowfullscreen", v...) }

func Async(v ...TypedValue) Mod             { return AttrValue("async", v...) }
func (n *Node) Async(v ...TypedValue) *Node { return n.AttrValue("async", v...) }

func Autofocus(v ...TypedValue) Mod             { return AttrValue("autofocus", v...) }
func (n *Node) Autofocus(v ...TypedValue) *Node { return n.AttrValue("autofocus", v...) }

func Autoplay(v ...TypedValue) Mod             { return AttrValue("autoplay", v...) }
func (n *Node) Autoplay(v ...TypedValue) *Node { return n.AttrValue("autoplay", v...) }

func Checked(v ...TypedValue) Mod             { return AttrValue("checked", v...) }
func (n *Node) Checked(v ...TypedValue) *Node { return n.AttrValue("checked", v...) }

func Controls(v ...TypedValue) Mod             { return AttrValue("controls", v...) }
func (n *Node) Controls(v ...TypedValue) *Node { return n.AttrValue("controls", v...) }

func Default(v ...TypedValue) Mod             { return AttrValue("default", v...) }
func (n *Node) Default(v ...TypedValue) *Node { return n.AttrValue("default", v...) }

func Defer(v ...TypedValue) Mod             { return AttrValue("defer", v...) }
func (n *Node) Defer(v ...TypedValue) *Node { return n.AttrValue("defer", v...) }

func Disabled(v ...TypedValue) Mod             { return AttrValue("disabled", v...) }
func (n *Node) Disabled(v ...TypedValue) *Node { return n.AttrValue("disabled", v...) }

func FormNoValidate(v ...TypedValue) Mod             { return AttrValue("formnovalidate", v...) }
func (n *Node) FormNoValidate(v ...TypedValue) *Node { return n.AttrValue("formnovalidate", v...) }

func Inert(v ...TypedValue) Mod             { return AttrValue("inert", v...) }
func (n *Node) Inert(v ...TypedValue) *Node { return n.AttrValue("inert", v...) }

func IsMap(v ...TypedValue) Mod             { return AttrValue("ismap", v...) }
func (n *Node) IsMap(v ...TypedValue) *Node { return n.AttrValue("ismap", v...) }

func ItemScope(v ...TypedValue) Mod             { return AttrValue("itemscope", v...) }
func (n *Node) ItemScope(v ...TypedValue) *Node { return n.AttrValue("itemscope", v...) }

func Loop(v ...TypedValue) Mod             { return AttrValue("loop", v...) }
func (n *Node) Loop(v ...TypedValue) *Node { return n.AttrValue("loop", v...) }

func Muted(v ...TypedValue) Mod             { return AttrValue("muted", v...) }
func (n *Node) Muted(v ...TypedValue) *Node { return n.AttrValue("muted", v...) }

func NoModule(v ...TypedValue) Mod             { return AttrValue("nomodule", v...) }
func (n *Node) NoModule(v ...TypedValue) *Node { return n.AttrValue("nomodule", v...) }

func Novalidate(v ...TypedValue) Mod             { return AttrValue("novalidate", v...) }
func (n *Node) Novalidate(v ...TypedValue) *Node { return n.AttrValue("novalidate", v...) }

func Open(v ...TypedValue) Mod             { return AttrValue("open", v...) }
func (n *Node) Open(v ...TypedValue) *Node { return n.AttrValue("open", v...) }

func PlaysInline(v ...TypedValue) Mod             { return AttrValue("playsinline", v...) }
func (n *Node) PlaysInline(v ...TypedValue) *Node { return n.AttrValue("playsinline", v...) }

func Readonly(v ...TypedValue) Mod             { return AttrValue("readonly", v...) }
func (n *Node) Readonly(v ...TypedValue) *Node { return n.AttrValue("readonly", v...) }

func Required(v ...TypedValue) Mod             { return AttrValue("required", v...) }
func (n *Node) Required(v ...TypedValue) *Node { return n.AttrValue("required", v...) }

func Selected(v ...TypedValue) Mod             { return AttrValue("selected", v...) }
func (n *Node) Selected(v ...TypedValue) *Node { return n.AttrValue("selected", v...) }

func Multiple(v ...TypedValue) Mod             { return AttrValue("multiple", v...) }
func (n *Node) Multiple(v ...TypedValue) *Node { return n.AttrValue("multiple", v...) }

func Reversed(v ...TypedValue) Mod             { return AttrValue("reversed", v...) }
func (n *Node) Reversed(v ...TypedValue) *Node { return n.AttrValue("reversed", v...) }

/**/

func ColSpan(v int) Mod             { return AttrValue("colspan", Int(v)) }
func (n *Node) ColSpan(v int) *Node { return n.AttrValue("colspan", Int(v)) }

func RowSpan(v int) Mod             { return AttrValue("rowspan", Int(v)) }
func (n *Node) RowSpan(v int) *Node { return n.AttrValue("rowspan", Int(v)) }

func Start(v int) Mod             { return AttrValue("start", Int(v)) }
func (n *Node) Start(v int) *Node { return n.AttrValue("start", Int(v)) }

func Headers(v string) Mod             { return Attr("headers", v) }
func (n *Node) Headers(v string) *Node { return n.Attr("headers", v) }

func DateTime(v string) Mod             { return Attr("datetime", v) }
func (n *Node) DateTime(v string) *Node { return n.Attr("datetime", v) }

func Loading(v string) Mod             { return Attr("loading", v) }
func (n *Node) Loading(v string) *Node { return n.Attr("loading", v) }

func Decoding(v string) Mod             { return Attr("decoding", v) }
func (n *Node) Decoding(v string) *Node { return n.Attr("decoding", v) }

// Download sets the "download" attribute,
// If value is omitted, it sets a boolean attribute.
func Download(value ...string) Mod {
	if len(value) > 0 {
		v := value[0]
		return Attr("download", v)
	}
	return Attr("download")
}

// Download sets the "download" attribute,
// If value is omitted, it sets a boolean attribute.
func (n *Node) Download(value ...string) *Node {
	if len(value) > 0 {
		v := value[0]
		return n.Attr("download", v)
	}
	return n.Attr("download")
}

func Ping(v string) Mod             { return Attr("ping", v) }
func (n *Node) Ping(v string) *Node { return n.Attr("ping", v) }

func Wrap(v string) Mod             { return Attr("wrap", v) }
func (n *Node) Wrap(v string) *Node { return n.Attr("wrap", v) }

func Viewport(v string) Mod {
	return func(n *Node) { n.Attr("name", "viewport").Attr("content", v) }
}
func (n *Node) Viewport(v string) *Node {
	return n.Attr("name", "viewport").Attr("content", v)
}

func ID(v string) Mod             { return Attr("id", v) }
func (n *Node) ID(v string) *Node { return n.Attr("id", v) }

func Data(name string, v TypedValue) Mod             { return AttrValue("data-"+name, v) }
func (n *Node) Data(name string, v TypedValue) *Node { return n.AttrValue("data-"+name, v) }

var uids atomic.Uint64

func UniqueID() string {
	v := uids.Add(1)
	return "id-" + strconv.FormatUint(v, 10)
}

func (n *Node) UniqueID() *Node {
	id := UniqueID()
	return n.ID(id).Var("htm_unique_id", id)
}

// event handlers

func On(event string, js string) Mod             { return Attr("on"+event, js) }
func (n *Node) On(event string, js string) *Node { return n.Attr("on"+event, js) }

// mouse

func OnClick(js string) Mod             { return Attr("onclick", js) }
func (n *Node) OnClick(js string) *Node { return n.Attr("onclick", js) }

func OnDblClick(js string) Mod             { return Attr("ondblclick", js) }
func (n *Node) OnDblClick(js string) *Node { return n.Attr("ondblclick", js) }

func OnMouseDown(js string) Mod             { return Attr("onmousedown", js) }
func (n *Node) OnMouseDown(js string) *Node { return n.Attr("onmousedown", js) }

func OnMouseUp(js string) Mod             { return Attr("onmouseup", js) }
func (n *Node) OnMouseUp(js string) *Node { return n.Attr("onmouseup", js) }

func OnMouseEnter(js string) Mod             { return Attr("onmouseenter", js) }
func (n *Node) OnMouseEnter(js string) *Node { return n.Attr("onmouseenter", js) }

func OnMouseLeave(js string) Mod             { return Attr("onmouseleave", js) }
func (n *Node) OnMouseLeave(js string) *Node { return n.Attr("onmouseleave", js) }

func OnMouseMove(js string) Mod             { return Attr("onmousemove", js) }
func (n *Node) OnMouseMove(js string) *Node { return n.Attr("onmousemove", js) }

func OnMouseOver(js string) Mod             { return Attr("onmouseover", js) }
func (n *Node) OnMouseOver(js string) *Node { return n.Attr("onmouseover", js) }

func OnMouseOut(js string) Mod             { return Attr("onmouseout", js) }
func (n *Node) OnMouseOut(js string) *Node { return n.Attr("onmouseout", js) }

func OnWheel(js string) Mod             { return Attr("onwheel", js) }
func (n *Node) OnWheel(js string) *Node { return n.Attr("onwheel", js) }

// keyboard

func OnKeyDown(js string) Mod             { return Attr("onkeydown", js) }
func (n *Node) OnKeyDown(js string) *Node { return n.Attr("onkeydown", js) }

func OnKeyUp(js string) Mod             { return Attr("onkeyup", js) }
func (n *Node) OnKeyUp(js string) *Node { return n.Attr("onkeyup", js) }

func OnKeyPress(js string) Mod             { return Attr("onkeypress", js) }
func (n *Node) OnKeyPress(js string) *Node { return n.Attr("onkeypress", js) }

// controls

func OnChange(js string) Mod             { return Attr("onchange", js) }
func (n *Node) OnChange(js string) *Node { return n.Attr("onchange", js) }

func OnInput(js string) Mod             { return Attr("oninput", js) }
func (n *Node) OnInput(js string) *Node { return n.Attr("oninput", js) }

func OnSubmit(js string) Mod             { return Attr("onsubmit", js) }
func (n *Node) OnSubmit(js string) *Node { return n.Attr("onsubmit", js) }

func OnReset(js string) Mod             { return Attr("onreset", js) }
func (n *Node) OnReset(js string) *Node { return n.Attr("onreset", js) }

func OnFocus(js string) Mod             { return Attr("onfocus", js) }
func (n *Node) OnFocus(js string) *Node { return n.Attr("onfocus", js) }

func OnBlur(js string) Mod             { return Attr("onblur", js) }
func (n *Node) OnBlur(js string) *Node { return n.Attr("onblur", js) }

func OnSelect(js string) Mod             { return Attr("onselect", js) }
func (n *Node) OnSelect(js string) *Node { return n.Attr("onselect", js) }

// drag & drop

func OnDrag(js string) Mod             { return Attr("ondrag", js) }
func (n *Node) OnDrag(js string) *Node { return n.Attr("ondrag", js) }

func OnDragStart(js string) Mod             { return Attr("ondragstart", js) }
func (n *Node) OnDragStart(js string) *Node { return n.Attr("ondragstart", js) }

func OnDragEnd(js string) Mod             { return Attr("ondragend", js) }
func (n *Node) OnDragEnd(js string) *Node { return n.Attr("ondragend", js) }

func OnDragEnter(js string) Mod             { return Attr("ondragenter", js) }
func (n *Node) OnDragEnter(js string) *Node { return n.Attr("ondragenter", js) }

func OnDragLeave(js string) Mod             { return Attr("ondragleave", js) }
func (n *Node) OnDragLeave(js string) *Node { return n.Attr("ondragleave", js) }

func OnDragOver(js string) Mod             { return Attr("ondragover", js) }
func (n *Node) OnDragOver(js string) *Node { return n.Attr("ondragover", js) }

func OnDrop(js string) Mod             { return Attr("ondrop", js) }
func (n *Node) OnDrop(js string) *Node { return n.Attr("ondrop", js) }

// clipboard

func OnCopy(js string) Mod             { return Attr("oncopy", js) }
func (n *Node) OnCopy(js string) *Node { return n.Attr("oncopy", js) }

func OnCut(js string) Mod             { return Attr("oncut", js) }
func (n *Node) OnCut(js string) *Node { return n.Attr("oncut", js) }

func OnPaste(js string) Mod             { return Attr("onpaste", js) }
func (n *Node) OnPaste(js string) *Node { return n.Attr("onpaste", js) }

// other

func OnLoad(js string) Mod             { return Attr("onload", js) }
func (n *Node) OnLoad(js string) *Node { return n.Attr("onload", js) }

func OnError(js string) Mod             { return Attr("onerror", js) }
func (n *Node) OnError(js string) *Node { return n.Attr("onerror", js) }

func OnScroll(js string) Mod             { return Attr("onscroll", js) }
func (n *Node) OnScroll(js string) *Node { return n.Attr("onscroll", js) }

// aria

func Aria(name string, v TypedValue) Mod             { return AttrValue("aria-"+name, v) }
func (n *Node) Aria(name string, v TypedValue) *Node { return n.AttrValue("aria-"+name, v) }
