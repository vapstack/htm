package htm

func DOCTYPE() *Node            { return RawString("<!DOCTYPE html>") }
func A(m ...Mod) *Node          { return Build("a", m...) }
func Link(m ...Mod) *Node       { return Build("link", m...) }
func Abbr(m ...Mod) *Node       { return Build("abbr", m...) }
func Address(m ...Mod) *Node    { return Build("address", m...) }
func Area(m ...Mod) *Node       { return Build("area", m...) }
func Article(m ...Mod) *Node    { return Build("article", m...) }
func Aside(m ...Mod) *Node      { return Build("aside", m...) }
func Audio(m ...Mod) *Node      { return Build("audio", m...) }
func B(m ...Mod) *Node          { return Build("b", m...) }
func Base(m ...Mod) *Node       { return Build("base", m...) }
func Bdi(m ...Mod) *Node        { return Build("bdi", m...) }
func Bdo(m ...Mod) *Node        { return Build("bdo", m...) }
func Blockquote(m ...Mod) *Node { return Build("blockquote", m...) }
func Body(m ...Mod) *Node       { return Build("body", m...) }
func Br(m ...Mod) *Node         { return Build("br", m...) }
func Button(m ...Mod) *Node     { return Build("button", m...) }
func Canvas(m ...Mod) *Node     { return Build("canvas", m...) }
func Caption(m ...Mod) *Node    { return Build("caption", m...) }
func Cite(m ...Mod) *Node       { return Build("cite", m...) }
func Code(m ...Mod) *Node       { return Build("code", m...) }
func Col(m ...Mod) *Node        { return Build("col", m...) }
func Colgroup(m ...Mod) *Node   { return Build("colgroup", m...) }
func Datalist(m ...Mod) *Node   { return Build("datalist", m...) }
func Dd(m ...Mod) *Node         { return Build("dd", m...) }
func Del(m ...Mod) *Node        { return Build("del", m...) }
func Details(m ...Mod) *Node    { return Build("details", m...) }
func Dfn(m ...Mod) *Node        { return Build("dfn", m...) }
func Dialog(m ...Mod) *Node     { return Build("dialog", m...) }
func Div(m ...Mod) *Node        { return Build("div", m...) }
func Dl(m ...Mod) *Node         { return Build("dl", m...) }
func Dt(m ...Mod) *Node         { return Build("dt", m...) }
func Em(m ...Mod) *Node         { return Build("em", m...) }
func Embed(m ...Mod) *Node      { return Build("embed", m...) }
func Fieldset(m ...Mod) *Node   { return Build("fieldset", m...) }
func Figcaption(m ...Mod) *Node { return Build("figcaption", m...) }
func Figure(m ...Mod) *Node     { return Build("figure", m...) }
func Footer(m ...Mod) *Node     { return Build("footer", m...) }
func Form(m ...Mod) *Node       { return Build("form", m...) }
func H1(m ...Mod) *Node         { return Build("h1", m...) }
func H2(m ...Mod) *Node         { return Build("h2", m...) }
func H3(m ...Mod) *Node         { return Build("h3", m...) }
func H4(m ...Mod) *Node         { return Build("h4", m...) }
func H5(m ...Mod) *Node         { return Build("h5", m...) }
func H6(m ...Mod) *Node         { return Build("h6", m...) }
func Head(m ...Mod) *Node       { return Build("head", m...) }
func Header(m ...Mod) *Node     { return Build("header", m...) }
func Hr(m ...Mod) *Node         { return Build("hr", m...) }
func Html(m ...Mod) *Node       { return Build("html", m...) }
func I(m ...Mod) *Node          { return Build("i", m...) }
func Iframe(m ...Mod) *Node     { return Build("iframe", m...) }
func Img(m ...Mod) *Node        { return Build("img", m...) }
func Input(m ...Mod) *Node      { return Build("input", m...) }
func Ins(m ...Mod) *Node        { return Build("ins", m...) }
func Kbd(m ...Mod) *Node        { return Build("kbd", m...) }
func Label(m ...Mod) *Node      { return Build("label", m...) }
func Legend(m ...Mod) *Node     { return Build("legend", m...) }
func Li(m ...Mod) *Node         { return Build("li", m...) }
func Main(m ...Mod) *Node       { return Build("main", m...) }
func Map(m ...Mod) *Node        { return Build("map", m...) }
func Mark(m ...Mod) *Node       { return Build("mark", m...) }
func Meta(m ...Mod) *Node       { return Build("meta", m...) }
func Meter(m ...Mod) *Node      { return Build("meter", m...) }
func Nav(m ...Mod) *Node        { return Build("nav", m...) }
func Noscript(m ...Mod) *Node   { return Build("noscript", m...) }
func Object(m ...Mod) *Node     { return Build("object", m...) }
func Ol(m ...Mod) *Node         { return Build("ol", m...) }
func Optgroup(m ...Mod) *Node   { return Build("optgroup", m...) }
func Option(m ...Mod) *Node     { return Build("option", m...) }
func Output(m ...Mod) *Node     { return Build("output", m...) }
func P(m ...Mod) *Node          { return Build("p", m...) }
func Param(m ...Mod) *Node      { return Build("param", m...) }
func Picture(m ...Mod) *Node    { return Build("picture", m...) }
func Pre(m ...Mod) *Node        { return Build("pre", m...) }
func Progress(m ...Mod) *Node   { return Build("progress", m...) }
func Q(m ...Mod) *Node          { return Build("q", m...) }
func Rp(m ...Mod) *Node         { return Build("rp", m...) }
func Rt(m ...Mod) *Node         { return Build("rt", m...) }
func Ruby(m ...Mod) *Node       { return Build("ruby", m...) }
func S(m ...Mod) *Node          { return Build("s", m...) }
func Samp(m ...Mod) *Node       { return Build("samp", m...) }
func Script(m ...Mod) *Node     { return Build("script", m...) }
func Section(m ...Mod) *Node    { return Build("section", m...) }
func Select(m ...Mod) *Node     { return Build("select", m...) }
func Small(m ...Mod) *Node      { return Build("small", m...) }
func Source(m ...Mod) *Node     { return Build("source", m...) }
func Span(m ...Mod) *Node       { return Build("span", m...) }
func Strong(m ...Mod) *Node     { return Build("strong", m...) }
func Sub(m ...Mod) *Node        { return Build("sub", m...) }
func Summary(m ...Mod) *Node    { return Build("summary", m...) }
func Sup(m ...Mod) *Node        { return Build("sup", m...) }
func Table(m ...Mod) *Node      { return Build("table", m...) }
func Tbody(m ...Mod) *Node      { return Build("tbody", m...) }
func Td(m ...Mod) *Node         { return Build("td", m...) }
func Template(m ...Mod) *Node   { return Build("template", m...) }
func Textarea(m ...Mod) *Node   { return Build("textarea", m...) }
func Tfoot(m ...Mod) *Node      { return Build("tfoot", m...) }
func Th(m ...Mod) *Node         { return Build("th", m...) }
func Thead(m ...Mod) *Node      { return Build("thead", m...) }
func Time(m ...Mod) *Node       { return Build("time", m...) }
func Tr(m ...Mod) *Node         { return Build("tr", m...) }
func Track(m ...Mod) *Node      { return Build("track", m...) }
func U(m ...Mod) *Node          { return Build("u", m...) }
func Ul(m ...Mod) *Node         { return Build("ul", m...) }
func Video(m ...Mod) *Node      { return Build("video", m...) }
func Wbr(m ...Mod) *Node        { return Build("wbr", m...) }

func VarTag(value string, mods ...Mod) *Node {
	n := Build("var").Apply(mods)
	if !n.HasContent() {
		n.Text(value)
	}
	return n
}

func Title(value string, mods ...Mod) *Node {
	n := Build("title").Apply(mods)
	if !n.HasContent() {
		n.Text(value)
	}
	return n
}
func TitleValue(value TypedValue, mods ...Mod) *Node {
	n := Build("title").Apply(mods)
	if !n.HasContent() {
		n.TextValue(value)
	}
	return n
}

func Stylesheet(href string, mods ...Mod) *Node {
	return Link().Attr("rel", "stylesheet").Attr("href", href).Apply(mods)
}

func Icon(href string, mods ...Mod) *Node {
	return Link().Attr("rel", "icon").Attr("href", href).Apply(mods)
}

func SlotTag(name string, m ...Mod) *Node { return Build("slot").Name(name).Apply(m) }
func StyleTag(m ...Mod) *Node             { return Build("style", m...) }
func DataTag(v string, mods ...Mod) *Node { return Build("data").Attr("value", v).Apply(mods) }

/**/

/*
package htm

func DOCTYPE() *Node { return RawString("<!DOCTYPE html>") }

func Link(m ...Mod) *Node { return Build("a", m...) }

func Rel(rel string, mods ...Mod) *Node {
	return Build("link").SetAttr("rel", rel).Apply(mods)
}

func Title(value any, mods ...Mod) *Node {
	n := Build("title").Apply(mods)
	if !n.HasContent() {
		n.SetText(value)
	}
	return n
}
func Variable(value any, mods ...Mod) *Node {
	n := Build("var").Apply(mods)
	if !n.HasContent() {
		n.SetText(value)
	}
	return n
}

func Abbr(m ...Mod) *Node       { return Build("abbr", m...) }
func Address(m ...Mod) *Node    { return Build("address", m...) }
func Area(m ...Mod) *Node       { return Build("area", m...) }
func Article(m ...Mod) *Node    { return Build("article", m...) }
func Aside(m ...Mod) *Node      { return Build("aside", m...) }
func Audio(m ...Mod) *Node      { return Build("audio", m...) }
func B(m ...Mod) *Node          { return Build("b", m...) }
func Base(m ...Mod) *Node       { return Build("base", m...) }
func Bdi(m ...Mod) *Node        { return Build("bdi", m...) }
func Bdo(m ...Mod) *Node        { return Build("bdo", m...) }
func Blockquote(m ...Mod) *Node { return Build("blockquote", m...) }
func Body(m ...Mod) *Node       { return Build("body", m...) }
func Br(m ...Mod) *Node         { return Build("br", m...) }
func Button(m ...Mod) *Node     { return Build("button", m...) }
func Canvas(m ...Mod) *Node     { return Build("canvas", m...) }
func Caption(m ...Mod) *Node    { return Build("caption", m...) }
func Cite(m ...Mod) *Node       { return Build("cite", m...) }
func Code(m ...Mod) *Node       { return Build("code", m...) }
func Col(m ...Mod) *Node        { return Build("col", m...) }
func Colgroup(m ...Mod) *Node   { return Build("colgroup", m...) }

func DataValue(v any, mods ...Mod) *Node {
	return Build("data").SetAttr("value", v).Apply(mods)
}

func Datalist(m ...Mod) *Node   { return Build("datalist", m...) }
func Dd(m ...Mod) *Node         { return Build("dd", m...) }
func Del(m ...Mod) *Node        { return Build("del", m...) }
func Details(m ...Mod) *Node    { return Build("details", m...) }
func Dfn(m ...Mod) *Node        { return Build("dfn", m...) }
func Dialog(m ...Mod) *Node     { return Build("dialog", m...) }
func Div(m ...Mod) *Node        { return Build("div", m...) }
func Dl(m ...Mod) *Node         { return Build("dl", m...) }
func Dt(m ...Mod) *Node         { return Build("dt", m...) }
func Em(m ...Mod) *Node         { return Build("em", m...) }
func Embed(m ...Mod) *Node      { return Build("embed", m...) }
func Fieldset(m ...Mod) *Node   { return Build("fieldset", m...) }
func Figcaption(m ...Mod) *Node { return Build("figcaption", m...) }
func Figure(m ...Mod) *Node     { return Build("figure", m...) }
func Footer(m ...Mod) *Node     { return Build("footer", m...) }
func Form(m ...Mod) *Node       { return Build("form", m...) }
func H1(m ...Mod) *Node         { return Build("h1", m...) }
func H2(m ...Mod) *Node         { return Build("h2", m...) }
func H3(m ...Mod) *Node         { return Build("h3", m...) }
func H4(m ...Mod) *Node         { return Build("h4", m...) }
func H5(m ...Mod) *Node         { return Build("h5", m...) }
func H6(m ...Mod) *Node         { return Build("h6", m...) }
func Head(m ...Mod) *Node       { return Build("head", m...) }
func Header(m ...Mod) *Node     { return Build("header", m...) }
func Hr(m ...Mod) *Node         { return Build("hr", m...) }
func Html(m ...Mod) *Node       { return Build("html", m...) }
func I(m ...Mod) *Node          { return Build("i", m...) }
func Iframe(m ...Mod) *Node     { return Build("iframe", m...) }
func Img(m ...Mod) *Node        { return Build("img", m...) }
func Input(m ...Mod) *Node      { return Build("input", m...) }
func Ins(m ...Mod) *Node        { return Build("ins", m...) }
func Kbd(m ...Mod) *Node        { return Build("kbd", m...) }
func Label(m ...Mod) *Node      { return Build("label", m...) }
func Legend(m ...Mod) *Node     { return Build("legend", m...) }
func Li(m ...Mod) *Node         { return Build("li", m...) }
func Main(m ...Mod) *Node       { return Build("main", m...) }
func Map(m ...Mod) *Node        { return Build("map", m...) }
func Mark(m ...Mod) *Node       { return Build("mark", m...) }
func Meta(m ...Mod) *Node       { return Build("meta", m...) }
func Meter(m ...Mod) *Node      { return Build("meter", m...) }
func Nav(m ...Mod) *Node        { return Build("nav", m...) }
func Noscript(m ...Mod) *Node   { return Build("noscript", m...) }
func Object(m ...Mod) *Node     { return Build("object", m...) }
func Ol(m ...Mod) *Node         { return Build("ol", m...) }
func Optgroup(m ...Mod) *Node   { return Build("optgroup", m...) }
func Option(m ...Mod) *Node     { return Build("option", m...) }
func Output(m ...Mod) *Node     { return Build("output", m...) }
func P(m ...Mod) *Node          { return Build("p", m...) }
func Param(m ...Mod) *Node      { return Build("param", m...) }
func Picture(m ...Mod) *Node    { return Build("picture", m...) }
func Pre(m ...Mod) *Node        { return Build("pre", m...) }
func Progress(m ...Mod) *Node   { return Build("progress", m...) }
func Q(m ...Mod) *Node          { return Build("q", m...) }
func Rp(m ...Mod) *Node         { return Build("rp", m...) }
func Rt(m ...Mod) *Node         { return Build("rt", m...) }
func Ruby(m ...Mod) *Node       { return Build("ruby", m...) }
func S(m ...Mod) *Node          { return Build("s", m...) }
func Samp(m ...Mod) *Node       { return Build("samp", m...) }
func Script(m ...Mod) *Node     { return Build("script", m...) }
func Section(m ...Mod) *Node    { return Build("section", m...) }
func Select(m ...Mod) *Node     { return Build("select", m...) }
func Small(m ...Mod) *Node      { return Build("small", m...) }
func Source(m ...Mod) *Node     { return Build("source", m...) }
func Span(m ...Mod) *Node       { return Build("span", m...) }
func Strong(m ...Mod) *Node     { return Build("strong", m...) }
func StyleBlock(m ...Mod) *Node { return Build("style", m...) }
func Sub(m ...Mod) *Node        { return Build("sub", m...) }
func Summary(m ...Mod) *Node    { return Build("summary", m...) }
func Sup(m ...Mod) *Node        { return Build("sup", m...) }
func Table(m ...Mod) *Node      { return Build("table", m...) }
func Tbody(m ...Mod) *Node      { return Build("tbody", m...) }
func Td(m ...Mod) *Node         { return Build("td", m...) }
func Template(m ...Mod) *Node   { return Build("template", m...) }
func Textarea(m ...Mod) *Node   { return Build("textarea", m...) }
func Tfoot(m ...Mod) *Node      { return Build("tfoot", m...) }
func Th(m ...Mod) *Node         { return Build("th", m...) }
func Thead(m ...Mod) *Node      { return Build("thead", m...) }
func Time(m ...Mod) *Node       { return Build("time", m...) }
func Tr(m ...Mod) *Node         { return Build("tr", m...) }
func Track(m ...Mod) *Node      { return Build("track", m...) }
func U(m ...Mod) *Node          { return Build("u", m...) }
func Ul(m ...Mod) *Node         { return Build("ul", m...) }
func Video(m ...Mod) *Node      { return Build("video", m...) }
func Wbr(m ...Mod) *Node        { return Build("wbr", m...) }

func TemplateSlot(name string, m ...Mod) *Node { return Build("slot").Name(name).Apply(m) }
*/
