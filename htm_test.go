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

func Test_Clear_PreservesTagAndFlags(t *testing.T) {
	customWrite := func(*Node, io.Writer) error { return nil }
	n := Build("img").
		Attr("id", "x").
		Class("a").
		Content(Text("x")).
		Var("k", "v").
		Slot("s", Text("slot")).
		Postpone(func(*Node) {})

	n.UnsafeScript()
	n.SetWriteFn(customWrite)
	n.value = String("payload")
	tagBefore, voidBefore := n.GetTag()
	flagBefore := n.flag

	n.Clear()

	tagAfter, voidAfter := n.GetTag()
	if tagAfter != tagBefore || voidAfter != voidBefore {
		t.Fatalf("clear must preserve tag/void flag: before=(%q,%v) after=(%q,%v)", tagBefore, voidBefore, tagAfter, voidAfter)
	}
	if n.flag != flagBefore {
		t.Fatalf("clear must preserve flags: before=%08b after=%08b", flagBefore, n.flag)
	}
	if n.HasAttr("id") || n.HasClass("a") || n.HasContent() || n.HasVar("k") || n.HasSlot("s") {
		t.Fatalf("clear must remove attrs/classes/content/vars/slots")
	}
	if len(n.postponed) != 0 {
		t.Fatalf("clear must remove postponed mods")
	}
	if n.writeFn == nil {
		t.Fatalf("clear must preserve write fn")
	}
	if !n.value.Valid() {
		t.Fatalf("clear must not remove value")
	}
}

func Test_Reset_RestoresPoolDefaults(t *testing.T) {
	customWrite := func(*Node, io.Writer) error { return nil }
	n := Build("img").
		Attr("id", "x").
		Class("a").
		Content(Text("x")).
		Var("k", "v").
		Slot("s", Text("slot")).
		Postpone(func(*Node) {})

	n.UnsafeScript()
	n.SetWriteFn(customWrite)
	n.value = String("payload")

	n.Reset()

	tag, isVoid := n.GetTag()
	if tag != "div" || isVoid {
		t.Fatalf("reset must restore default tag, got (%q,%v)", tag, isVoid)
	}
	if n.flag != 0 {
		t.Fatalf("reset must clear all non-owned flags, got %08b", n.flag)
	}
	if n.HasAttr("id") || n.HasClass("a") || n.HasContent() || n.HasVar("k") || n.HasSlot("s") {
		t.Fatalf("reset must remove attrs/classes/content/vars/slots")
	}
	if len(n.postponed) != 0 {
		t.Fatalf("reset must remove postponed mods")
	}
	if n.writeFn != nil {
		t.Fatalf("reset must remove write fn")
	}
	if n.value.Valid() {
		t.Fatalf("reset must remove value")
	}
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
	if len(n.ExtractSlotNodes("x")) != 2 {
		t.Fatalf("expected 2 nodes after append")
	}

	// restore slot and test prepend
	n.Slot("x", a, b)
	c := Text("C")
	n.PrependSlot("x", c)
	ex := n.ExtractSlotNodes("x")
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

	ex := dst.ExtractSlotNodes("a")
	if len(ex) != 2 {
		t.Fatalf("expected 2 group nodes, got %v", len(ex))
	}
	if ex[0].String() != "1" || ex[1].String() != "2" {
		t.Fatalf("unexpected moved slot content: %q, %q", ex[0].String(), ex[1].String())
	}
	for _, x := range ex {
		x.Release()
	}
}

/**/

func Benchmark_Build(b *testing.B) {
	cnt := buildBasic().Count(CountTree)

	buildBasic().Release() // warmup

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		buildBasic().Release()
	}
	b.StopTimer()
	b.ReportMetric(float64(cnt), "nodes/op")
}

func buildBasic() *Node {
	return Div().Class("flex flex-col items-center p-7 rounded-2xl").Attr("id", "root").Content(
		Span().Class("a b c").Text("hello"),
		Span().Attr("data-x", "1").Text("world"),
	)
}

func Benchmark_Build_Mods(b *testing.B) {
	build := func() *Node {
		return Div(Class("flex flex-col items-center p-7 rounded-2xl"), Attr("id", "root"), Content(
			Span(Class("a b c"), TextContent("hello")),
			Span(Attr("data-x", "1"), TextContent("world")),
		))
	}

	cnt := build().Count(CountTree)

	build().Release() // warmup

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		build().Release()
	}
	b.StopTimer()
	b.ReportMetric(float64(cnt), "nodes/op")
}

func Benchmark_Render(b *testing.B) {
	n := buildBasic()
	cnt := n.Count(CountTree)

	b.Run("Default", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			if err := n.Render(io.Discard); err != nil {
				b.Fatal(err)
			}
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})

	b.Run("AssumeNoReplace", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			if err := n.Render(io.Discard, RenderAssumeNoReplace); err != nil {
				b.Fatal(err)
			}
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})
}

func Benchmark_Render_Mods(b *testing.B) {
	n := Div(
		Class("flex flex-col items-center p-7 rounded-2xl"),
		Attr("id", "root"),
		Content(
			Span(Class("a b c"), Content(Text("hello"))),
			Span(Attr("data-x", "1"), Content(Text("world"))),
		),
	)

	cnt := n.Count(CountTree)

	b.Run("Default", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			if err := n.Render(io.Discard); err != nil {
				b.Fatal(err)
			}
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})

	b.Run("AssumeNoReplace", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			if err := n.Render(io.Discard); err != nil {
				b.Fatal(err)
			}
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})
}

func Benchmark_BuildRender(b *testing.B) {

	build := func() *Node {
		return Div().Class("flex flex-col items-center p-7 rounded-2xl").Attr("id", "root").Content(
			Span().Class("a b c").Text("hello"),
			Span().Attr("data-x", "1").Text("world"),
		)
	}

	cnt := build().Count(CountTree)

	build().Release() // warmup

	b.Run("Default", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			n := build()
			if err := n.Render(io.Discard); err != nil {
				b.Fatal(err)
			}
			n.Release()
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})

	build().Release() // warmup

	b.Run("AssumeNoReplace", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			n := build()
			if err := n.Render(io.Discard, RenderAssumeNoReplace); err != nil {
				b.Fatal(err)
			}
			n.Release()
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})

}

func Benchmark_BuildRender_Mods(b *testing.B) {
	build := func() *Node {
		return Div(
			Class("flex flex-col items-center p-7 rounded-2xl"),
			Attr("id", "root"),
			Content(
				Span(Class("a b c"), Content(Text("hello"))),
				Span(Attr("data-x", "1"), Content(Text("world"))),
			),
		)
	}

	cnt := build().Count(CountTree)

	build().Release() // warmup

	b.Run("Default", func(b *testing.B) {

		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			n := build()
			if err := n.Render(io.Discard); err != nil {
				b.Fatal(err)
			}
			n.Release()
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})

	build().Release() // warmup

	b.Run("AssumeNoReplace", func(b *testing.B) {

		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			n := build()
			if err := n.Render(io.Discard, RenderAssumeNoReplace); err != nil {
				b.Fatal(err)
			}
			n.Release()
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})
}

func Benchmark_Class_SetMulti(b *testing.B) {
	n := Build("div")
	defer n.Release()

	s1 := "flex flex-col items-center p-7 rounded-2xl"
	s2 := "size-48 shadow-xl rounded-md"
	s3 := "flex gap-2 font-medium text-gray-600 dark:text-gray-400"

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		n.class.reset()
		n.Class(s1)
		n.Class(s2)
		n.Class(s3)
	}
}

func Benchmark_Class_SetAndMoveClassPrefix(b *testing.B) {
	src := Build("div")
	dst := Build("div")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		src.class.reset()
		dst.class.reset()

		src.Class("abc1 xyz1 abc2 xyz2 xyz3 abc3")
		src.MoveClassPrefixTo(dst, "abc")
	}
}

func Benchmark_Attrs_SetAndMoveAttrPrefix(b *testing.B) {
	src := Build("div")
	dst := Build("div")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		src.attrs.reset()
		dst.attrs.reset()

		src.Attr("data-a", "1").Attr("data-b", "2").Attr("x-a", "3").Attr("data-c", "4")
		src.MoveAttrPrefixTo(dst, "data-")
	}
}

/**/

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

func buildCompareChaining() *Node {
	list := Ul().Class("user-list")
	for _, item := range benchData {
		list.Append(
			Li().Class("user-item").AttrValue("id", Int(item.ID)).Content(
				Span().Class("name").Text(item.Name),
				A().Href("mailto:"+item.Email).Text(item.Email), // strings concatenation will allocate
			),
		)
	}
	return list
}

func Benchmark_Compare_Htm(b *testing.B) {
	cnt := buildCompareChaining().Count(CountTree)

	buildCompareChaining().Release() // warmup

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		list := buildCompareChaining()
		_ = list.Render(io.Discard)
		list.Release()
	}
	b.StopTimer()
	b.ReportMetric(float64(cnt), "nodes/op")
}

func Benchmark_Compare_Htm_Mods(b *testing.B) {
	build := func() *Node {
		list := Ul(Class("user-list"))
		for _, item := range benchData {
			list.Append(
				Li(Class("user-item"), AttrValue("id", Int(item.ID)), Content(
					Span(Class("name"), TextContent(item.Name)),
					A(Href("mailto:"+item.Email), TextContent(item.Email)), // strings concatenation will allocate
				)),
			)
		}
		return list
	}
	cnt := build().Count(CountTree)

	build().Release() // warmup

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		list := build()
		_ = list.Render(io.Discard)
		list.Release()
	}
	b.StopTimer()
	b.ReportMetric(float64(cnt), "nodes/op")
}

// standard html/template

func Benchmark_Compare_StdTemplate(b *testing.B) {
	const tplString = `<ul class="user-list">{{range .}}<li class="user-item" id="{{.ID}}"><span class="name">{{.Name}}</span><a href="mailto:{{.Email}}">{{.Email}}</a></li>{{end}}</ul>`
	tpl := template.Must(template.New("list").Parse(tplString))

	cnt := buildCompareChaining().Count(CountTree)

	_ = tpl.Execute(io.Discard, benchData) // warmup?

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = tpl.Execute(io.Discard, benchData)
	}
	b.StopTimer()
	b.ReportMetric(float64(cnt), "nodes/op")
}

func Benchmark_Compare_Htm_Static(b *testing.B) {
	static := Static(func() *Node {
		x := Ul().Class("user-list")
		for _, item := range benchData {
			x.Append(
				Li().Class("user-item").AttrValue("id", Int(item.ID)).Content(
					Span().Class("name").Text(item.Name),
					A().Href("mailto:"+item.Email).Text(item.Email), // strings concatenation will allocate
				),
			)
		}
		return x
	})

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = static.Render(io.Discard)
	}
	b.StopTimer()
	b.ReportMetric(1, "nodes/op")
}

// the "baseline" speed, unsafe but fast

func Benchmark_Compare_RawString(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, 64*1024))
	run := func() {
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
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		run()
	}
	b.StopTimer()
	b.ReportMetric(float64(5*len(benchData)+1), "nodes/op")
}

/**/

func buildBenchPage() *Node {
	const rows = 100
	return Div().
		Class("min-h-screen bg-gray-50").
		Attr("id", "root").
		Content(
			Header().Class("sticky top-0 z-10 bg-white border-b").Content(
				Div().Class("mx-auto max-w-6xl px-6 py-4 flex items-center justify-between").Content(
					Div().Class("flex items-center gap-3").Content(
						Span().Class("text-xl font-semibold").Text("Acme"),
						Span().Class("text-sm text-gray-500").Text("Dashboard"),
					),
					Div().Class("flex items-center gap-3").Content(
						A().Class("text-sm text-blue-600 hover:underline").Href("#").Text("Docs"),
						A().Class("text-sm text-blue-600 hover:underline").Href("#").Text("Support"),
						Span().Class("inline-flex items-center rounded-full bg-gray-100 px-3 py-1 text-sm").Text("alice@acme.test"),
					),
				),
			),

			Div().Class("mx-auto max-w-6xl px-6 py-8 grid grid-cols-12 gap-6").Content(
				Aside().Class("col-span-3").Content(
					Nav().Class("rounded-2xl bg-white border p-4").Content(
						Div().Class("text-xs font-semibold text-gray-500 mb-3").Text("NAVIGATION"),
						Ul().Class("space-y-2").Content(
							Li().Content(A().Class("block rounded-lg px-3 py-2 bg-blue-50 text-blue-700").Href("#").Text("Overview")),
							Li().Content(A().Class("block rounded-lg px-3 py-2 hover:bg-gray-50").Href("#").Text("Users")),
							Li().Content(A().Class("block rounded-lg px-3 py-2 hover:bg-gray-50").Href("#").Text("Billing")),
							Li().Content(A().Class("block rounded-lg px-3 py-2 hover:bg-gray-50").Href("#").Text("Settings")),
						),
					),
				),

				Main().Class("col-span-9 space-y-6").Content(
					Div().Class("grid grid-cols-3 gap-4").Content(
						Div().Class("rounded-2xl bg-white border p-5").Content(
							Div().Class("text-sm text-gray-500").Text("Active users"),
							Div().Class("mt-2 text-2xl font-semibold").Text("1,284"),
							Div().Class("mt-1 text-xs text-green-600").Text("+8.2% this week"),
						),
						Div().Class("rounded-2xl bg-white border p-5").Content(
							Div().Class("text-sm text-gray-500").Text("Revenue"),
							Div().Class("mt-2 text-2xl font-semibold").Text("$32,140"),
							Div().Class("mt-1 text-xs text-green-600").Text("+3.1% this week"),
						),
						Div().Class("rounded-2xl bg-white border p-5").Content(
							Div().Class("text-sm text-gray-500").Text("Churn"),
							Div().Class("mt-2 text-2xl font-semibold").Text("1.7%"),
							Div().Class("mt-1 text-xs text-red-600").Text("+0.3pp this week"),
						),
					),

					Section().Class("rounded-2xl bg-white border").Content(
						Div().Class("px-6 py-4 border-b flex items-center justify-between").Content(
							Div().Content(
								Div().Class("text-base font-semibold").Text("Users"),
								Div().Class("text-sm text-gray-500").Text("A small sample of recent signups"),
							),
							Div().Class("flex items-center gap-2").Content(
								A().Class("text-sm rounded-lg border px-3 py-2 hover:bg-gray-50").Href("#").Text("Export"),
								A().Class("text-sm rounded-lg bg-blue-600 text-white px-3 py-2").Href("#").Text("Invite"),
							),
						),

						Div().Class("overflow-hidden").Content(
							Table().Class("w-full text-sm").Content(
								Thead().Class("bg-gray-50 text-gray-600").Content(
									Tr().Content(
										Th().Class("text-left font-medium px-6 py-3").Text("Name"),
										Th().Class("text-left font-medium px-6 py-3").Text("Email"),
										Th().Class("text-left font-medium px-6 py-3").Text("Role"),
										Th().Class("text-right font-medium px-6 py-3").Text("Status"),
									),
								),
								Tbody().Class("divide-y").Content(
									buildRows(rows),
								),
							),
						),

						Div().Class("px-6 py-4 border-t flex items-center justify-between").Content(
							Div().Class("text-sm text-gray-500").Text("Showing 1 to 10 of 1,284 results"),
							Div().Class("flex items-center gap-2").Content(
								A().Class("rounded-lg border px-3 py-2 text-sm hover:bg-gray-50").Href("#").Text("Previous"),
								A().Class("rounded-lg bg-gray-900 text-white px-3 py-2 text-sm").Href("#").Text("1"),
								A().Class("rounded-lg border px-3 py-2 text-sm hover:bg-gray-50").Href("#").Text("2"),
								A().Class("rounded-lg border px-3 py-2 text-sm hover:bg-gray-50").Href("#").Text("3"),
								A().Class("rounded-lg border px-3 py-2 text-sm hover:bg-gray-50").Href("#").Text("Next"),
							),
						),
					),
				),
			),

			Footer().Class("mt-10 pb-10").Content(
				Div().Class("mx-auto max-w-6xl px-6 text-sm text-gray-500").Text("(c) 2026 Acme Inc. All rights reserved."),
			),
		)
}

func buildRows(rows int) *Node {
	out := Group()
	for i := 0; i < rows; i++ {
		out.Append(
			Tr().Class("hover:bg-gray-50").Content(
				Td().Class("px-6 py-3").Content(
					Div().Class("font-medium").Text("Alice Johnson"),
					Div().Class("text-xs text-gray-500").Text("Amsterdam"),
				),
				Td().Class("px-6 py-3").Text("alice@example.com"),
				Td().Class("px-6 py-3").Text("Member"),
				Td().Class("px-6 py-3 text-right").Content(
					Span().Class("inline-flex items-center rounded-full bg-green-50 text-green-700 px-2 py-1 text-xs").Text("Active"),
				),
			),
		)
	}
	return out
}

func Benchmark_Page_Build(b *testing.B) {
	cnt := buildBenchPage().Count(CountTree)

	buildBenchPage().Release() // warmup

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		n := buildBenchPage()
		n.Release()
	}
	b.StopTimer()
	b.ReportMetric(float64(cnt), "nodes/op")
}

func Benchmark_Page_Render(b *testing.B) {
	buildBenchPage().Release() // warmup

	n := buildBenchPage()

	cnt := n.Count(CountTree)

	b.Run("Default", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			if err := n.Render(io.Discard); err != nil {
				b.Fatal(err)
			}
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})

	buildBenchPage().Release() // warmup

	b.Run("AssumeNoReplace", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			if err := n.Render(io.Discard, RenderAssumeNoReplace); err != nil {
				b.Fatal(err)
			}
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})

}

func Benchmark_Page_BuildRender(b *testing.B) {

	cnt := buildBenchPage().Count(CountTree)

	buildBenchPage().Release() // warmup

	b.Run("Default", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			n := buildBenchPage()
			if err := n.Render(io.Discard); err != nil {
				b.Fatal(err)
			}
			n.Release()
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})

	buildBenchPage().Release() // warmup

	b.Run("AssumeNoReplace", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			n := buildBenchPage()
			if err := n.Render(io.Discard, RenderAssumeNoReplace); err != nil {
				b.Fatal(err)
			}
			n.Release()
		}
		b.StopTimer()
		b.ReportMetric(float64(cnt), "nodes/op")
	})
}

func Benchmark_Count(b *testing.B) {

	n := buildBenchPage()

	b.Run("CountContent", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			n.Count()
		}
		b.StopTimer()
	})

	b.Run("CountRecursive", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			n.Count(CountRecursive)
		}
		b.StopTimer()
	})

	b.Run("CountTree", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			n.Count(CountTree)
		}
		b.StopTimer()
	})

	b.Run("CountTreeNonNil", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			n.Count(CountTree | CountSkipNil)
		}
		b.StopTimer()
	})

}
