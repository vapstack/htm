package htm

import (
	"strings"
	"testing"
)

func Test_Render_BasicTagVoidAndNormal(t *testing.T) {

	div := Build("div").ID("x").Class("a b")
	got1 := div.String()

	if got1 != `<div class="a b" id="x"></div>` {
		t.Fatalf("unexpected render: %q", got1)
	}

	img := Img().Src("x")
	got2 := img.String()

	if got2 != `<img src="x"/>` {
		t.Fatalf("unexpected void render: %q", got2)
	}
}

func Test_Render_TextEscaping(t *testing.T) {
	n := Div().Text(`<>"'&`)
	got := n.String()

	if got != `<div>&lt;&gt;&#34;&#39;&amp;</div>` {
		t.Fatalf("unexpected escaping: %q", got)
	}
}

func Test_Render_Class_NoLeadingSpacesWhenSomeInactive(t *testing.T) {
	n := Div()
	n.Class("a")
	n.Class("b")
	n.RemoveClass("a")
	n.Class("c")

	got := n.String()

	if got != `<div class="b c"></div>` {
		t.Fatalf("unexpected class render: %q", got)
	}
}

func Test_Render_Attr_FalseMeansAbsent(t *testing.T) {
	n := Div()
	n.Attr("data-x", "1")
	n.AttrValue("data-y", Unset)
	n.AttrValue("disabled", Bool(false))

	got := n.String()

	if strings.Contains(got, "data-y") {
		t.Fatalf("data-y must not be rendered: %q", got)
	}
	if strings.Contains(got, "disabled") {
		t.Fatalf("disabled=false must not be rendered: %q", got)
	}
	if got != `<div data-x="1"></div>` {
		t.Fatalf("unexpected attrs render: %q", got)
	}
}

func Test_Render_BoolAttr_TrueRendersNameOnly(t *testing.T) {
	n := Input().AttrValue("disabled", Bool(true))
	got := n.String()
	if got != `<input disabled/>` {
		t.Fatalf("unexpected bool attr render: %q", got)
	}
}

func Test_Attr_MovePrefixSuffix(t *testing.T) {
	src := Div().Attr("data-a", "1").Attr("data-b", "2").Attr("x-a", "3")
	dst := Div()
	defer src.Release()
	defer dst.Release()

	src.MoveAttrPrefixTo(dst, "data-")

	if src.GetAttr("data-a").Valid() || src.GetAttr("data-b").Valid() {
		t.Fatalf("expected data-* moved out of src")
	}
	if dst.GetAttr("data-a").StringOrZero() != "1" || dst.GetAttr("data-b").StringOrZero() != "2" {
		t.Fatalf("expected data-* moved into dst")
	}
	if src.GetAttr("x-a").StringOrZero() != "3" {
		t.Fatalf("expected x-a to remain in src")
	}
}

func Test_Class_MovePrefixSuffix(t *testing.T) {
	src := Build("div").Class("a x:1 y:2 z")
	dst := Build("div")
	defer src.Release()
	defer dst.Release()

	src.MoveClassPrefixTo(dst, "x:", "y:")

	if src.HasClass("x:1") || src.HasClass("y:2") {
		t.Fatalf("expected x:/y: classes moved out of src")
	}
	if !dst.HasClass("x:1") || !dst.HasClass("y:2") {
		t.Fatalf("expected x:/y: classes present in dst")
	}
	if !src.HasClass("a") || !src.HasClass("z") {
		t.Fatalf("expected a and z to remain in src")
	}
}

func Test_Pool_GetReleaseCanReuse(t *testing.T) {
	n := Get()
	n.SetTag("div").Attr("id", "x").Class("a")
	n.Release()

	m := Get()

	if got := m.String(); got != `<div></div>` {
		m.Release()
		t.Fatalf("expected reset node, got: %q", got)
	}
	m.Release()
}

func Test_Static_CachesAndAvoidsReRender(t *testing.T) {
	var calls int
	fn := func() *Node {
		calls++
		return Build("div").Attr("id", "x").Text("hi")
	}

	a := Static(fn)
	s1 := a.String()
	a.Release()

	b := Static(fn)
	s2 := b.String()
	b.Release()

	if calls > 1 {
		t.Fatalf("expected fn to be called once, got %d", calls)
	}
	if s1 != `<div id="x">hi</div>` || s2 != s1 {
		t.Fatalf("unexpected output: %q / %q", s1, s2)
	}
}

func Test_StaticContent_UsesCachedRaw(t *testing.T) {
	fn := func() *Node {
		return Group(
			Build("span").Content(Text("a")),
			Build("span").Content(Text("b")),
		)
	}

	n1 := Build("div").StaticContent(fn)
	n2 := Build("div").StaticContent(fn)

	s1 := n1.String()
	s2 := n2.String()

	n1.Release()
	n2.Release()

	if s1 != `<div><span>a</span><span>b</span></div>` {
		t.Fatalf("unexpected: %q", s1)
	}
	if s2 != s1 {
		t.Fatalf("expected same output, got %q vs %q", s2, s1)
	}

	if !strings.Contains(s1, "<span>") {
		t.Fatalf("expected spans, got: %q", s1)
	}
}

func Test_Slots_SetAppendPrependDeleteExtract(t *testing.T) {
	n := Build("div")
	defer n.Release()

	a := Text("A")
	n.Slot("x", a)
	if !n.HasSlot("x") {
		t.Fatalf("expected slot x to exist")
	}

	b := Text("B")
	n.AppendSlot("x", b)
	if len(n.ExtractSlot("x")) != 2 {
		t.Fatalf("expected 2 nodes after append")
	}

	// restore slot and test prepend
	n.Slot("x", a, b)
	c := Text("C")
	n.PrependSlot("x", c)
	ex := n.ExtractSlot("x")
	if len(ex) != 3 {
		t.Fatalf("expected 3 nodes after prepend, got %d", len(ex))
	}
	if ex[0].String() != `C` { // Text node renders as plain text
		t.Fatalf("expected first to be C, got %q", ex[0].String())
	}
	// release extracted manually (ExtractSlot detaches without put)
	for _, x := range ex {
		x.Release()
	}

	// DeleteSlot releases children
	n.Slot("x", Text("X"))
	n.DeleteSlot("x")
	if n.HasSlot("x") {
		t.Fatalf("expected slot x to be empty after delete")
	}
}

func Test_Slots_MoveSlotTo(t *testing.T) {
	src := Build("div")
	dst := Build("div")
	defer src.Release()
	defer dst.Release()

	src.Slot("a", Text("1"), Text("2"))
	dst.Slot("a", Text("X"))

	src.MoveSlotTo(dst, "a")

	if src.HasSlot("a") {
		t.Fatalf("expected src slot to be empty after move")
	}
	if !dst.HasSlot("a") {
		t.Fatalf("expected dst slot to have content after move")
	}

	ex := dst.ExtractSlot("a")
	if len(ex) != 2 {
		t.Fatalf("expected 2 moved nodes, got %d", len(ex))
	}
	if ex[0].String() != "1" || ex[1].String() != "2" {
		t.Fatalf("unexpected moved slot content: %q, %q", ex[0].String(), ex[1].String())
	}
	for _, x := range ex {
		x.Release()
	}
}
