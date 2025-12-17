package htm

import (
	"bytes"
	"html/template"
	"io"
	"strconv"
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

/**/

func Benchmark_Render_SimpleTree(b *testing.B) {
	root := Div().
		Class("flex flex-col items-center p-7 rounded-2xl").
		Attr("id", "root").
		Content(
			Build("span").Class("a b c").Text("hello"),
			Build("span").Attr("data-x", "1").Text("world"),
		)
	defer root.Release()

	var buf bytes.Buffer
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := root.Render(&buf); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Render_SimpleTree_Mods(b *testing.B) {
	root := Div(
		Class("flex flex-col items-center p-7 rounded-2xl"),
		Attr("id", "root"),
		Content(
			Span(
				Class("a b c"),
				Content(Text("hello")),
			),
			Span(
				Attr("data-x", "1"),
				Content(Text("world")),
			),
		),
	)
	defer root.Release()

	var buf bytes.Buffer
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := root.Render(&buf); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Render_TreeFullCycle(b *testing.B) {

	buf := new(bytes.Buffer)
	b.ReportAllocs()
	b.ResetTimer()

	fn := func(buf *bytes.Buffer) {
		n := Div().Class("flex flex-col items-center p-7 rounded-2xl").Attr("id", "root").
			Content(
				Span().Class("a b c").Text("hello"),
				Span().Attr("data-x", "1").Text("world"),
				Div().Class("flex flex-col items-center p-7 rounded-2xl").Attr("id", "root").Content(
					Span().Class("a b c").Text("hello"),
					Span().Attr("data-x", "1").Text("world"),
					Div().Class("flex flex-col items-center p-7 rounded-2xl").Attr("id", "root").Content(
						Span().Class("a b c").Text("hello"),
						Span().Attr("data-x", "1").Text("world"),
					),
				),
			)
		defer n.Release()

		if err := n.Render(buf); err != nil {
			b.Fatal(err)
		}
	}

	for i := 0; i < b.N; i++ {
		buf.Reset()
		fn(buf)
	}
}

func Benchmark_Build(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	fn := func() {
		n := Div().Class("flex flex-col items-center p-7 rounded-2xl").Attr("id", "root").
			Content(
				Span().Class("a b c").Text("hello"),
				Span().Attr("data-x", "1").Text("world"),
				Div().Class("flex flex-col items-center p-7 rounded-2xl").Attr("id", "root").Content(
					Span().Class("a b c").Text("hello"),
					Span().Attr("data-x", "1").Text("world"),
					Div().Class("flex flex-col items-center p-7 rounded-2xl").Attr("id", "root").Content(
						Span().Class("a b c").Text("hello"),
						Span().Attr("data-x", "1").Text("world"),
					),
				),
			)
		defer n.Release()
	}

	for i := 0; i < b.N; i++ {
		fn()
	}
}

func Benchmark_Render_TreeFullCycle_Mods(b *testing.B) {
	buf := new(bytes.Buffer)
	b.ReportAllocs()
	b.ResetTimer()

	fn := func(buf *bytes.Buffer) {
		n := Div(
			Class("flex flex-col items-center p-7 rounded-2xl"),
			Attr("id", "root"),
			Content(
				Span(
					Class("a b c"),
					Content(Text("hello")),
				),
				Span(
					Attr("data-x", "1"),
					Content(Text("world")),
				),
				Div(
					Class("flex flex-col items-center p-7 rounded-2xl"),
					Attr("id", "root"),
					Content(
						Span(
							Class("a b c"),
							Content(Text("hello")),
						),
						Span(
							Attr("data-x", "1"),
							Content(Text("world")),
						),
						Div(
							Class("flex flex-col items-center p-7 rounded-2xl"),
							Attr("id", "root"),
							Content(
								Span(
									Class("a b c"),
									Content(Text("hello")),
								),
								Span(
									Attr("data-x", "1"),
									Content(Text("world")),
								),
							),
						),
					),
				),
			),
		)
		defer n.Release()

		if err := n.Render(buf); err != nil {
			b.Fatal(err)
		}
	}

	for i := 0; i < b.N; i++ {
		buf.Reset()
		fn(buf)
	}
}

func Benchmark_Build_Mods(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	fn := func() {
		n := Div(
			Class("flex flex-col items-center p-7 rounded-2xl"),
			Attr("id", "root"),
			Content(
				Span(
					Class("a b c"),
					Content(Text("hello")),
				),
				Span(
					Attr("data-x", "1"),
					Content(Text("world")),
				),
				Div(
					Class("flex flex-col items-center p-7 rounded-2xl"),
					Attr("id", "root"),
					Content(
						Span(
							Class("a b c"),
							Content(Text("hello")),
						),
						Span(
							Attr("data-x", "1"),
							Content(Text("world")),
						),
						Div(
							Class("flex flex-col items-center p-7 rounded-2xl"),
							Attr("id", "root"),
							Content(
								Span(
									Class("a b c"),
									Content(Text("hello")),
								),
								Span(
									Attr("data-x", "1"),
									Content(Text("world")),
								),
							),
						),
					),
				),
			),
		)
		defer n.Release()
	}

	for i := 0; i < b.N; i++ {
		fn()
	}
}

func Benchmark_Class_SetMulti(b *testing.B) {
	n := Build("div")
	defer n.Release()

	s1 := "flex flex-col items-center p-7 rounded-2xl"
	s2 := "size-48 shadow-xl rounded-md"
	s3 := "flex gap-2 font-medium text-gray-600 dark:text-gray-400"

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// mimic typical build usage
		n.class.reset()
		n.Class(s1)
		n.Class(s2)
		n.Class(s3)
	}
}

func Benchmark_Attrs_SetAndMovePrefix(b *testing.B) {
	src := Build("div")
	dst := Build("div")
	defer src.Release()
	defer dst.Release()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		src.attrs.reset()
		dst.attrs.reset()

		src.Attr("data-a", "1").Attr("data-b", "2").Attr("x-a", "3").Attr("data-c", "4")
		src.MoveAttrPrefixTo(dst, "data-")
	}
}

/**/

// Typical list item model
type Item struct {
	ID    int
	Name  string
	Email string
}

var benchData = []Item{
	{1, "Alice", "alice@example.com"},
	{2, "Bob", "bob@example.com"},
	{3, "Charlie", "charlie@example.com"},
	{4, "Dave", "dave@example.com"},
	{5, "Eve", "eve@example.com"},
}

func Benchmark_Compare_Htm(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		list := Ul().Class("user-list")

		for _, item := range benchData {
			list.Append(
				Li().Class("user-item").AttrValue("id", Int(item.ID)).Content(
					Span().Class("name").Text(item.Name),
					A().Href("mailto:"+item.Email).Text(item.Email), // strings concatenation will allocate
				),
			)
		}

		_ = list.Render(io.Discard)
		list.Release()
	}
}

func Benchmark_Compare_Htm_Mods(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		list := Ul(Class("user-list"))

		for _, item := range benchData {
			list.Append(
				Li(Class("user-item"), AttrValue("id", Int(item.ID)), Content(
					Span(Class("name"), TextContent(item.Name)),
					A(Href("mailto:"+item.Email), TextContent(item.Email)), // strings concatenation will allocate
				)),
			)
		}

		_ = list.Render(io.Discard)
		list.Release()
	}
}

// standard html/template

func Benchmark_Compare_StdTemplate(b *testing.B) {
	const tplString = `<ul class="user-list">{{range .}}<li class="user-item" id="{{.ID}}"><span class="name">{{.Name}}</span><a href="mailto:{{.Email}}">{{.Email}}</a></li>{{end}}</ul>`
	tpl := template.Must(template.New("list").Parse(tplString))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = tpl.Execute(io.Discard, benchData)
	}
}

// the "baseline" speed, unsafe but fast

func Benchmark_Compare_RawString(b *testing.B) {
	b.ReportAllocs()
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.WriteString(`<ul class="user-list">`)

		for _, item := range benchData {
			buf.WriteString(`<li class="user-item" id="`)
			buf.WriteString(strconv.Itoa(item.ID))
			buf.WriteString(`"><span class="name">`)
			buf.WriteString(item.Name)
			buf.WriteString(`</span><a href="`)
			buf.WriteString("mailto:" + item.Email)
			buf.WriteString(`">`)
			buf.WriteString(item.Email)
			buf.WriteString(`</a></li>`)
		}

		buf.WriteString(`</ul>`)
		_, _ = buf.WriteTo(io.Discard)
	}
}
