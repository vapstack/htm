package htm

import (
	"bytes"
	"html/template"
	"io"
	"strconv"
	"testing"
)

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
