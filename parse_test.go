package htm

import (
	"os"
	"strings"
	"testing"
)

func TestParse_BasicTree(t *testing.T) {
	n, err := Parse([]byte(`<div id="root" class="a b"><span data-x="1">hello</span></div>`))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	got := n.String()
	want := `<div class="a b" id="root"><span data-x="1">hello</span></div>`
	if got != want {
		t.Fatalf("unexpected render\nwant: %s\n got: %s", want, got)
	}
}

func TestParse_MultipleRootsBecomeGroup(t *testing.T) {
	n, err := Parse([]byte(`<span>a</span><span>b</span>`))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != `<span>a</span><span>b</span>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_StrictMismatchReturnsError(t *testing.T) {
	n, err := Parse([]byte(`<div><span></div>`))
	if n != nil {
		n.Release()
	}
	if err == nil {
		t.Fatal("expected error for mismatched closing tag")
	}
}

func TestParse_ScriptContentDeniedByDefault(t *testing.T) {
	n, err := Parse([]byte(`<script>console.log("x")</script>`))
	if n != nil {
		n.Release()
	}
	if err == nil {
		t.Fatal("expected error when script content is denied by default")
	}
}

func TestParse_ScriptContentCanBeAllowed(t *testing.T) {
	n, err := Parse([]byte(`<script>console.log("x")</script>`), ParseAllowScriptContent)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != `<script>console.log("x")</script>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_TopLevelRawContentMode(t *testing.T) {
	n, err := Parse(
		[]byte(`<div id="x"><span>a</span><b>b</b></div>`),
		ParseTopLevelRawContent,
	)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if n.Len() != 1 {
		t.Fatalf("expected exactly one raw child, got %d", n.Len())
	}
	if got := n.String(); got != `<div id="x"><span>a</span><b>b</b></div>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_TopLevelRawContentMode_Script(t *testing.T) {
	n, err := Parse(
		[]byte(`<script>console.log("x")</script>`),
		ParseTopLevelRawContent|ParseAllowScriptContent,
	)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != `<script>console.log("x")</script>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_TopLevelRawContentModeRejectsMultipleRoots(t *testing.T) {
	n, err := Parse(
		[]byte(`<div>1</div><div>2</div>`),
		ParseTopLevelRawContent,
	)
	if n != nil {
		n.Release()
	}
	if err == nil {
		t.Fatal("expected error for multiple roots in TopLevelRawContent mode")
	}
}

func TestParse_ReuseInputBufferOption(t *testing.T) {
	buf := []byte(`<div id="a">x</div>`)
	n, err := Parse(buf, ParseReuseBuffer)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	buf[9] = 'b' // id="a" -> id="b"
	if got := n.String(); !strings.Contains(got, `id="b"`) {
		t.Fatalf("expected render to reflect input mutation, got: %s", got)
	}
}

func TestParse_CopyInputByDefault(t *testing.T) {
	buf := []byte(`<div id="a">x</div>`)
	n, err := Parse(buf)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	buf[9] = 'b'
	if got := n.String(); strings.Contains(got, `id="b"`) {
		t.Fatalf("did not expect render to reflect input mutation, got: %s", got)
	}
	if got := n.String(); !strings.Contains(got, `id="a"`) {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_CommentsDroppedByDefault(t *testing.T) {
	n, err := Parse([]byte(`<!DOCTYPE html><!--x--><div>ok</div>`))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != `<!DOCTYPE html><div>ok</div>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_CommentsCanBeKept(t *testing.T) {
	n, err := Parse([]byte(`<!DOCTYPE html><!--x--><div>ok</div>`), ParseKeepComments)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != `<!DOCTYPE html><!--x--><div>ok</div>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_EntitiesDecodedInTextAndAttrs(t *testing.T) {
	n, err := Parse([]byte(`<div title="Tom &amp; Jerry">A &amp; B</div>`))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != `<div title="Tom &amp; Jerry">A &amp; B</div>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_WhitespaceMode_DropFormatting(t *testing.T) {
	n, err := Parse(
		[]byte("<div>\n  <span>a</span>\n  <span>b</span>\n</div>"),
		ParseKeepEdgeWhitespace,
	)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != `<div><span>a</span><span>b</span></div>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_WhitespaceMode_DropFormatting_PreservesPre(t *testing.T) {
	n, err := Parse(
		[]byte("<pre>\n  x y  \n</pre><div>\n  <span>a</span>\n</div>"),
		ParseKeepEdgeWhitespace,
	)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != "<pre>\n  x y  \n</pre><div><span>a</span></div>" {
		t.Fatalf("unexpected render: %q", got)
	}
}

func TestParse_WhitespaceMode_DropAll(t *testing.T) {
	n, err := Parse([]byte("<div>  hello world  <span>\t x y \n</span>  z  </div>"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != `<div>hello world<span>x y</span>z</div>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func TestParse_WhitespaceMode_DropAll_PreservesSensitiveTags(t *testing.T) {
	n, err := Parse([]byte("<title> A B </title><textarea> C D </textarea><div> E F </div>"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	defer n.Release()

	if got := n.String(); got != `<title> A B </title><textarea> C D </textarea><div>E F</div>` {
		t.Fatalf("unexpected render: %s", got)
	}
}

func Benchmark_Parse_Icon(b *testing.B) {

	data := []byte(`<svg
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="m14 11 4-4 4 4" />
  <path d="M18 16V7" />
  <path d="m2 16 4.039-9.69a.5.5 0 0 1 .923 0L11 16" />
  <path d="M3.304 13h6.392" />
</svg>
`)

	bench := func(b *testing.B, name string, flags ParseFlag) {
		x, err := Parse(data, flags)
		if err != nil {
			b.Fatal(err)
		}
		cnt := x.Count()
		x.Release()

		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				n, err := Parse(data, flags)
				if err != nil {
					b.Fatal(err)
				}
				n.Release()
			}
			b.ReportMetric(float64(cnt), "nodes/op")
		})
	}

	bench(b, "Default", 0)
	bench(b, "ReuseBuffer_DropAll", ParseReuseBuffer)
	bench(b, "ReuseBuffer_KeepEdgeWhitespace", ParseReuseBuffer|ParseKeepEdgeWhitespace)
	bench(b, "ReuseBuffer_KeepWhitespace", ParseReuseBuffer|ParseKeepWhitespace)
	bench(b, "ReuseBuffer_TopLevelOnly", ParseReuseBuffer|ParseTopLevelRawContent)
}

func Benchmark_Parse_Page(b *testing.B) {

	data, _ := os.ReadFile("parse_test.html")

	x, _ := Parse(data, ParseReuseBuffer|ParseAllowScriptContent)
	// _ = os.WriteFile("parse_test.out", []byte(x.String()), 0o600)
	x.Release()

	bench := func(b *testing.B, name string, flags ParseFlag) {
		x, err := Parse(data, flags)
		if err != nil {
			b.Fatal(err)
		}
		cnt := x.Count()
		x.Release()

		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				n, err := Parse(data, flags)
				if err != nil {
					b.Fatal(err)
				}
				n.Release()
			}
			b.ReportMetric(float64(cnt), "nodes/op")
		})
	}

	bench(b, "DropAll", ParseAllowScriptContent)
	bench(b, "ReuseBuffer_DropAll", ParseReuseBuffer|ParseAllowScriptContent)
	bench(b, "ReuseBuffer_KeepEdgeWhitespace", ParseReuseBuffer|ParseAllowScriptContent|ParseKeepEdgeWhitespace)
	bench(b, "ReuseBuffer_KeepWhitespace", ParseReuseBuffer|ParseAllowScriptContent|ParseKeepWhitespace)
	bench(b, "ReuseBuffer_KeepWhitespace_KeepComments", ParseReuseBuffer|ParseAllowScriptContent|ParseKeepWhitespace|ParseKeepComments)
}
