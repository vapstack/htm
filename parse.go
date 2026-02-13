package htm

import (
	"bytes"
	"fmt"
	"html"
	"slices"
	"sync"
	"unsafe"
)

type ParseFlag uint8

const (
	// ParseReuseBuffer allows parsed node values to reference the input buffer directly.
	// When enabled, mutating the original input bytes may corrupt the output.
	ParseReuseBuffer ParseFlag = 1 << iota
	// ParseKeepWhitespace keeps all whitespace text as-is.
	ParseKeepWhitespace
	// ParseKeepEdgeWhitespace keeps whitespace in non-empty text nodes and drops whitespace-only text nodes.
	ParseKeepEdgeWhitespace
	// ParseTopLevelRawContent parses exactly one top-level element and stores its entire inner HTML as one raw child node.
	// If input contains more than one top-level node, Parse returns an error.
	ParseTopLevelRawContent
	// ParseAllowScriptContent allows non-empty <script> content.
	ParseAllowScriptContent
	// ParseKeepComments keeps HTML comments as raw nodes.
	ParseKeepComments
)

// Parse builds Node tree from HTML bytes.
//
// For multiple top-level nodes, Parse returns a Group node.
func Parse(data []byte, flags ...ParseFlag) (*Node, error) {
	var mode ParseFlag
	for _, f := range flags {
		mode |= f
	}

	p := getParser(data, mode)

	var (
		root *Node
		err  error
	)
	if mode&ParseTopLevelRawContent != 0 {
		root, err = p.parseTopLevelRawContent()
	} else {
		root, err = p.parseTree()
	}
	if err != nil {
		p.releaseOnError()
		putParser(p)
		return nil, err
	}
	putParser(p)
	return root, nil
}

type htmlParser struct {
	data []byte
	pos  int
	mode ParseFlag

	roots []*Node
	stack []*Node
}

var htmlParserPool = sync.Pool{
	New: func() any {
		p := &htmlParser{}
		p.roots = make([]*Node, 0, 16)
		p.stack = make([]*Node, 0, 16)
		return p
	},
}

func getParser(data []byte, mode ParseFlag) *htmlParser {
	p := htmlParserPool.Get().(*htmlParser)
	p.data = data
	p.mode = mode
	return p
}

func putParser(p *htmlParser) {
	p.data = nil
	p.pos = 0
	p.mode = 0

	clear(p.roots)
	p.roots = p.roots[:0]

	clear(p.stack)
	p.stack = p.stack[:0]

	htmlParserPool.Put(p)
}

func (p *htmlParser) parseTree() (*Node, error) {
	for {
		if len(p.stack) == 0 {
			p.skipSpaces()
		}
		if p.pos >= len(p.data) {
			break
		}
		if p.data[p.pos] != '<' {
			if err := p.parseText(); err != nil {
				return nil, err
			}
			continue
		}
		if err := p.parseTagOrMarkup(); err != nil {
			return nil, err
		}
	}

	if len(p.stack) > 0 {
		n := p.stack[len(p.stack)-1]
		return nil, p.errAtf("unclosed tag %q", n.tag)
	}

	switch len(p.roots) {
	case 0:
		return Group(), nil
	case 1:
		return p.roots[0], nil
	default:
		return Group(p.roots...), nil
	}
}

func (p *htmlParser) parseTopLevelRawContent() (*Node, error) {
	p.skipSpaces()
	if p.pos >= len(p.data) {
		return Group(), nil
	}
	if p.data[p.pos] != '<' {
		return nil, p.errAt("expected a top-level element")
	}

	node, selfClose, err := p.parseStartTag()
	if err != nil {
		return nil, err
	}
	p.roots = append(p.roots, node)

	if selfClose || node.flag&flagVoid != 0 {
		p.skipSpaces()
		if p.pos != len(p.data) {
			return nil, p.errAt("multiple top-level nodes are not allowed in TopLevelRawContent mode")
		}
		return node, nil
	}

	innerStart := p.pos
	innerEnd, afterClose, err := p.scanUntilMatchingClose(node.tag)
	if err != nil {
		return nil, err
	}
	if isScriptTag(node.tag) && innerEnd > innerStart {
		if p.mode&ParseAllowScriptContent == 0 {
			return nil, p.errAt("script content is not allowed by parse flags")
		}
		node.flag |= flagScript
	}
	if innerEnd > innerStart {
		node.content = append(node.content, p.newRawNode(p.data[innerStart:innerEnd]))
	}
	p.pos = afterClose

	p.skipSpaces()
	if p.pos != len(p.data) {
		return nil, p.errAt("multiple top-level nodes are not allowed in TopLevelRawContent mode")
	}
	return node, nil
}

func (p *htmlParser) scanUntilMatchingClose(rootTag string) (innerEnd int, afterClose int, err error) {
	if isRawTextTag(rootTag) {
		start, after, ok := findRawTextClose(p.data, rootTag, p.pos)
		if !ok {
			return 0, 0, p.errAtf("unclosed tag %q", rootTag)
		}
		return start, after, nil
	}

	tags := []string{rootTag}
	i := p.pos
	n := len(p.data)

	// keep a local open-tag stack to find the matching close for the root element
	for i < n {
		if p.data[i] != '<' {
			i++
			continue
		}

		if i+3 < n && p.data[i+1] == '!' && p.data[i+2] == '-' && p.data[i+3] == '-' {
			end := bytes.Index(p.data[i+4:], []byte("-->"))
			if end < 0 {
				return 0, 0, p.errAtPos(i, "unterminated comment")
			}
			i += 4 + end + 3
			continue
		}

		if i+1 >= n {
			return 0, 0, p.errAtPos(i, "unexpected EOF after '<'")
		}

		switch p.data[i+1] {

		case '!':
			k := bytes.IndexByte(p.data[i+2:], '>')
			if k < 0 {
				return 0, 0, p.errAtPos(i, "unterminated declaration")
			}
			i += 2 + k + 1
			continue

		case '/':
			name, next, e := scanEndTag(p.data, i)
			if e != nil {
				return 0, 0, p.errAtPos(i, e.Error())
			}
			if len(tags) == 0 {
				return 0, 0, p.errAtPos(i, "unexpected closing tag")
			}

			top := tags[len(tags)-1]
			if !equalASCIIFold(top, name) {
				return 0, 0, p.errAtPos(i, fmt.Sprintf("mismatched closing tag: expected </%v>, got </%v>", top, name))
			}

			tags = tags[:len(tags)-1]
			if len(tags) == 0 {
				return i, next, nil
			}
			i = next
			continue

		case '?':
			return 0, 0, p.errAtPos(i, "processing instructions are not supported")

		default:
			name, selfClose, next, e := scanStartTagMeta(p.data, i)
			if e != nil {
				return 0, 0, p.errAtPos(i, e.Error())
			}
			if isRawTextTag(name) && !selfClose {
				_, after, ok := findRawTextClose(p.data, name, next)
				if !ok {
					return 0, 0, p.errAtPos(i, fmt.Sprintf("unclosed tag %q", name))
				}
				i = after
				continue
			}
			if !selfClose && !isVoidTag(name) {
				tags = append(tags, name)
			}
			i = next
		}
	}

	return 0, 0, p.errAtf("unclosed tag %q", rootTag)
}

func (p *htmlParser) parseTagOrMarkup() error {
	if p.pos+1 >= len(p.data) {
		return p.errAt("unexpected EOF after '<'")
	}
	switch p.data[p.pos+1] {

	case '!':
		if p.pos+3 < len(p.data) && p.data[p.pos+2] == '-' && p.data[p.pos+3] == '-' {
			return p.parseComment()
		}
		return p.parseDeclaration()

	case '/':
		return p.parseEndTag()

	case '?':
		return p.errAt("processing instructions are not supported")
	}

	node, selfClose, err := p.parseStartTag()
	if err != nil {
		return err
	}
	p.attach(node)

	handled, err := p.parseRawTextElement(node)
	if err != nil {
		return err
	}
	if handled {
		return nil
	}

	if selfClose || node.flag&flagVoid != 0 {
		return nil
	}
	p.stack = append(p.stack, node)
	return nil
}

func (p *htmlParser) parseRawTextElement(node *Node) (bool, error) {
	if !isRawTextTag(node.tag) {
		return false, nil
	}

	innerStart := p.pos
	closeStart, closeAfter, ok := findRawTextClose(p.data, node.tag, p.pos)
	if !ok {
		return true, p.errAtf("unclosed tag %q", node.tag)
	}
	if closeStart > innerStart {
		if isScriptTag(node.tag) {
			if p.mode&ParseAllowScriptContent == 0 {
				return true, p.errAt("script content is not allowed by parse flags")
			}
			node.flag |= flagScript
		}
		node.content = append(node.content, p.newRawNode(p.data[innerStart:closeStart]))
	}
	p.pos = closeAfter
	return true, nil
}

func (p *htmlParser) parseComment() error {
	start := p.pos
	end := bytes.Index(p.data[start+4:], []byte("-->"))
	if end < 0 {
		return p.errAt("unterminated comment")
	}
	p.pos = start + 4 + end + 3
	if p.mode&ParseKeepComments != 0 {
		p.attach(p.newRawNode(p.data[start:p.pos]))
	}
	return nil
}

func (p *htmlParser) parseDeclaration() error {
	start := p.pos
	end := bytes.IndexByte(p.data[start+2:], '>')
	if end < 0 {
		return p.errAt("unterminated declaration")
	}
	p.pos = start + 2 + end + 1
	p.attach(p.newRawNode(p.data[start:p.pos]))
	return nil
}

func (p *htmlParser) parseText() error {
	start := p.pos
	if start >= len(p.data) {
		return nil
	}
	// jump to next tag boundary in one call instead of byte-by-byte scanning
	if i := bytes.IndexByte(p.data[start:], '<'); i >= 0 {
		p.pos = start + i
	} else {
		p.pos = len(p.data)
	}
	if p.pos <= start {
		return nil
	}
	segment := p.data[start:p.pos]

	if !p.inWhitespaceSensitiveContext() {
		if p.mode&ParseKeepWhitespace == 0 {
			if isAllSpace(segment) {
				return nil
			}
			if p.mode&ParseKeepEdgeWhitespace == 0 {
				node := p.newTextNodeTrimmed(segment)
				if node != nil {
					p.attach(node)
				}
				return nil
			}
		}
	}

	node := p.newTextNode(segment)
	if node != nil {
		p.attach(node)
	}
	return nil
}

func (p *htmlParser) parseStartTag() (*Node, bool, error) {

	data := p.data
	nData := len(data)

	if p.pos >= nData || data[p.pos] != '<' {
		return nil, false, p.errAt("expected '<'")
	}
	p.pos++

	nameStart := p.pos
	if p.pos >= nData || tagCharTable[data[p.pos]] != 1 {
		return nil, false, p.errAt("invalid tag name")
	}
	p.pos++
	for p.pos < nData && tagCharTable[data[p.pos]] != 0 {
		p.pos++
	}

	tagStr := p.bytesToString(data[nameStart:p.pos])

	n := Get()
	n.tag = tagStr

	// inline the full start-tag state machine to keep hot-path call overhead minimal
	selfClose := false
	for {
		for p.pos < nData && isSpace(data[p.pos]) {
			p.pos++
		}
		if p.pos >= nData {
			n.Release()
			return nil, false, p.errAt("unexpected EOF in start tag")
		}

		ch := data[p.pos]
		if ch == '>' {
			p.pos++
			break
		}
		if ch == '/' {
			if p.pos+1 < nData && data[p.pos+1] == '>' {
				selfClose = true
				p.pos += 2
				break
			}
			n.Release()
			return nil, false, p.errAt("unexpected '/' in start tag")
		}

		attrNameStart := p.pos
		first := data[p.pos]

		if isSpace(first) || first == '=' || first == '>' || first == '/' {
			n.Release()
			return nil, false, p.errAt("invalid attribute")
		}

		if (first >= '0' && first <= '9') || first == '-' || invalidAttrTable[first] != 0 {
			n.Release()
			return nil, false, p.errAt("invalid attribute")
		}

		p.pos++

		for p.pos < nData {
			c := data[p.pos]
			if isSpace(c) || c == '=' || c == '>' || c == '/' {
				break
			}
			if invalidAttrTable[c] != 0 {
				n.Release()
				return nil, false, p.errAt("invalid attribute")
			}
			p.pos++
		}

		nameBytes := data[attrNameStart:p.pos]
		name := p.bytesToString(nameBytes)

		for p.pos < nData && isSpace(data[p.pos]) {
			p.pos++
		}
		hasValue := false
		var value string
		if p.pos < nData && data[p.pos] == '=' {
			hasValue = true
			p.pos++
			for p.pos < nData && isSpace(data[p.pos]) {
				p.pos++
			}
			if p.pos >= nData {
				n.Release()
				return nil, false, p.errAt("unexpected EOF after '='")
			}

			if data[p.pos] == '"' || data[p.pos] == '\'' {
				quoteChar := data[p.pos]

				p.pos++
				vStart := p.pos

				for p.pos < nData && data[p.pos] != quoteChar {
					p.pos++
				}

				if p.pos >= nData {
					n.Release()
					return nil, false, p.errAt("unterminated quoted attribute value")
				}

				value = p.valueString(data[vStart:p.pos])
				p.pos++

			} else {

				vStart := p.pos

				for p.pos < nData {
					c := data[p.pos]
					if isSpace(c) || c == '>' || c == '/' {
						break
					}
					if c == '"' || c == '\'' || c == '<' || c == '=' || c == '`' {
						n.Release()
						return nil, false, p.errAt("invalid unquoted attribute value")
					}
					p.pos++
				}

				if p.pos == vStart {
					n.Release()
					return nil, false, p.errAt("empty attribute value")
				}

				value = p.valueString(data[vStart:p.pos])
			}
		}

		if isClassAttr(nameBytes) {

			// class is parsed into classMap directly to avoid generic attr writes

			if !hasValue {
				n.Release()
				return nil, false, p.errAt("class attribute must have a value")
			}

			hasClass := false
			for i := 0; i < len(value); i++ {
				if !isSpace(value[i]) {
					hasClass = true
					break
				}
			}
			if hasClass {
				n.class.setMulti(value, true)
			} else {
				n.class.setOne("", true)
			}
			continue
		}

		if hasValue {
			n.attrs.set(name, String(value))
		} else {
			n.attrs.set(name, Bool(true))
		}
	}

	if selfClose || isVoidTag(tagStr) {
		n.flag |= flagVoid
	}
	return n, selfClose, nil
}

func (p *htmlParser) parseEndTag() error {
	if p.pos+1 >= len(p.data) || p.data[p.pos] != '<' || p.data[p.pos+1] != '/' {
		return p.errAt("expected closing tag")
	}

	name, next, err := scanEndTag(p.data, p.pos)
	if err != nil {
		return p.errAt(err.Error())
	}

	p.pos = next

	if len(p.stack) == 0 {
		return p.errAtf("unexpected closing tag </%v>", name)
	}

	top := p.stack[len(p.stack)-1]
	if !equalASCIIFold(top.tag, name) {
		return p.errAtf("mismatched closing tag: expected </%v>, got </%v>", top.tag, name)
	}

	p.stack = p.stack[:len(p.stack)-1]
	return nil
}

func findRawTextClose(data []byte, tag string, from int) (closeStart int, closeAfter int, ok bool) {

	n := len(data)
	tagLen := len(tag)

	if tagLen == 0 {
		return 0, 0, false
	}

	if from < 0 {
		from = 0
	}

	// scan only "</" candidates and validate tag name lazily

	for i := from; i+tagLen+3 <= n; {
		j := bytes.Index(data[i:], rawClosePrefix)
		if j < 0 {
			return 0, 0, false
		}
		i += j
		if !equalASCIIFoldBytesString(data[i+2:i+2+tagLen], tag) {
			i += 2
			continue
		}
		k := i + 2 + tagLen
		for k < n && isSpace(data[k]) {
			k++
		}
		if k < n && data[k] == '>' {
			return i, k + 1, true
		}
		i += 2
	}
	return 0, 0, false
}

var rawClosePrefix = []byte("</")

func (p *htmlParser) attach(n *Node) {
	if n == nil {
		return
	}
	if len(p.stack) == 0 {
		p.roots = append(p.roots, n)
		return
	}
	parent := p.stack[len(p.stack)-1]
	parent.content = append(parent.content, n)
}

func (p *htmlParser) releaseOnError() {
	for _, root := range p.roots {
		if root != nil {
			root.Release()
		}
	}
	p.roots = nil
	p.stack = nil
}

func (p *htmlParser) newTextNode(b []byte) *Node {
	return p.newTextNodeFromString(p.valueString(b))
}

func (p *htmlParser) newTextNodeFromString(s string) *Node {
	if s == "" {
		return nil
	}
	n := Get()
	n.tag = "@text"
	n.value = String(s)
	n.writeFn = renderText
	return n
}

func (p *htmlParser) newTextNodeTrimmed(b []byte) *Node {
	if len(b) == 0 {
		return nil
	}

	if bytes.IndexByte(b, '&') < 0 {
		start, end := trimASCIISpaceBoundsBytes(b)
		if start >= end {
			return nil
		}
		return p.newTextNodeFromString(p.valueStringNoUnescape(b[start:end]))
	}

	s := p.valueString(b)
	if s == "" {
		return nil
	}

	start, end := trimASCIISpaceBoundsString(s)
	if start >= end {
		return nil
	}

	return p.newTextNodeFromString(s[start:end])
}

func (p *htmlParser) newRawNode(b []byte) *Node {
	if len(b) == 0 {
		return nil
	}
	n := Get()
	n.tag = "@raw"
	n.value = Bytes(p.byteSlice(b))
	n.writeFn = renderRaw
	return n
}

func (p *htmlParser) bytesToString(b []byte) string {
	if p.mode&ParseReuseBuffer != 0 {
		return unsafe.String(unsafe.SliceData(b), len(b))
	}
	return string(b)
}

func (p *htmlParser) valueString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if bytes.IndexByte(b, '&') < 0 {
		return p.valueStringNoUnescape(b)
	}
	return html.UnescapeString(p.valueStringNoUnescape(b))
}

func (p *htmlParser) valueStringNoUnescape(b []byte) string {
	return p.bytesToString(b)
}

func (p *htmlParser) byteSlice(b []byte) []byte {
	if p.mode&ParseReuseBuffer != 0 {
		return b
	}
	return slices.Clone(b)
}

func (p *htmlParser) skipSpaces() {
	for p.pos < len(p.data) && isSpace(p.data[p.pos]) {
		p.pos++
	}
}

func (p *htmlParser) inWhitespaceSensitiveContext() bool {
	if len(p.stack) == 0 {
		return false
	}
	tag := p.stack[len(p.stack)-1].tag
	return isWhitespaceSensitiveTag(tag)
}

func isWhitespaceSensitiveTag(tag string) bool {
	return eqASCIIFoldLit(tag, "pre") || eqASCIIFoldLit(tag, "textarea") || eqASCIIFoldLit(tag, "script") ||
		eqASCIIFoldLit(tag, "style") || eqASCIIFoldLit(tag, "title")
}

func (p *htmlParser) errAt(msg string) error {
	return fmt.Errorf("parse error at byte %v: %v", p.pos, msg)
}

func (p *htmlParser) errAtf(format string, args ...any) error {
	return p.errAt(fmt.Sprintf(format, args...))
}

func (p *htmlParser) errAtPos(pos int, msg string) error {
	return fmt.Errorf("parse error at byte %v: %v", pos, msg)
}

func scanEndTag(data []byte, pos int) (name string, next int, err error) {
	n := len(data)
	if pos+2 >= n || data[pos] != '<' || data[pos+1] != '/' {
		return "", pos, fmt.Errorf("invalid closing tag")
	}
	i := pos + 2
	start := i
	if i >= n || tagCharTable[data[i]] != 1 {
		return "", pos, fmt.Errorf("invalid closing tag name")
	}
	i++
	for i < n && tagCharTable[data[i]] != 0 {
		i++
	}
	name = unsafe.String(unsafe.SliceData(data[start:i]), i-start)
	for i < n && isSpace(data[i]) {
		i++
	}
	if i >= n || data[i] != '>' {
		return "", pos, fmt.Errorf("malformed closing tag")
	}
	return name, i + 1, nil
}

func scanStartTagMeta(data []byte, pos int) (name string, selfClose bool, next int, err error) {
	n := len(data)
	if pos >= n || data[pos] != '<' {
		return "", false, pos, fmt.Errorf("expected '<'")
	}
	i := pos + 1
	start := i
	if i >= n || tagCharTable[data[i]] != 1 {
		return "", false, pos, fmt.Errorf("invalid tag name")
	}
	i++
	for i < n && tagCharTable[data[i]] != 0 {
		i++
	}
	name = unsafe.String(unsafe.SliceData(data[start:i]), i-start)

	for i < n {
		if isSpace(data[i]) {
			i++
			continue
		}
		if data[i] == '>' {
			return name, false, i + 1, nil
		}
		if data[i] == '/' {
			if i+1 < n && data[i+1] == '>' {
				return name, true, i + 2, nil
			}
			return "", false, pos, fmt.Errorf("unexpected '/' in start tag")
		}

		// skip attribute name
		for i < n {
			c := data[i]
			if isSpace(c) || c == '=' || c == '>' || c == '/' {
				break
			}
			i++
		}
		if i >= n {
			return "", false, pos, fmt.Errorf("unexpected EOF in start tag")
		}
		for i < n && isSpace(data[i]) {
			i++
		}
		if i < n && data[i] == '=' {
			i++
			for i < n && isSpace(data[i]) {
				i++
			}
			_, _, nextV, scanErr := scanAttrValue(data, i)
			if scanErr != nil {
				return "", false, pos, scanErr
			}
			i = nextV
		}
	}
	return "", false, pos, fmt.Errorf("unexpected EOF in start tag")
}

func scanAttrValue(data []byte, pos int) (start int, end int, next int, err error) {
	n := len(data)
	if pos >= n {
		return 0, 0, pos, fmt.Errorf("unexpected EOF after '='")
	}
	if data[pos] == '"' || data[pos] == '\'' {
		q := data[pos]
		start = pos + 1
		i := start
		for i < n && data[i] != q {
			i++
		}
		if i >= n {
			return 0, 0, pos, fmt.Errorf("unterminated quoted attribute value")
		}
		return start, i, i + 1, nil
	}

	start = pos
	i := pos
	for i < n {
		c := data[i]
		if isSpace(c) || c == '>' || c == '/' {
			break
		}
		if c == '"' || c == '\'' || c == '<' || c == '=' || c == '`' {
			return 0, 0, pos, fmt.Errorf("invalid unquoted attribute value")
		}
		i++
	}
	if i == start {
		return 0, 0, pos, fmt.Errorf("empty attribute value")
	}
	return start, i, i, nil
}

func isRawTextTag(tag string) bool {
	return isScriptTag(tag) || isStyleTag(tag)
}

func isStyleTag(tag string) bool {
	if len(tag) != 5 {
		return false
	}
	return (tag[0]|0x20) == 's' && (tag[1]|0x20) == 't' && (tag[2]|0x20) == 'y' &&
		(tag[3]|0x20) == 'l' && (tag[4]|0x20) == 'e'
}

func isVoidTag(tag string) bool {
	switch len(tag) {
	case 2:
		return eqASCIIFoldLit(tag, "br") || eqASCIIFoldLit(tag, "hr")
	case 3:
		return eqASCIIFoldLit(tag, "col") || eqASCIIFoldLit(tag, "img") || eqASCIIFoldLit(tag, "wbr")
	case 4:
		return eqASCIIFoldLit(tag, "area") || eqASCIIFoldLit(tag, "base") || eqASCIIFoldLit(tag, "link") || eqASCIIFoldLit(tag, "meta")
	case 5:
		return eqASCIIFoldLit(tag, "embed") || eqASCIIFoldLit(tag, "input") || eqASCIIFoldLit(tag, "track")
	case 6:
		return eqASCIIFoldLit(tag, "source")
	default:
		return false
	}
}

func isClassAttr(name []byte) bool {
	return len(name) == 5 &&
		(name[0]|0x20) == 'c' &&
		(name[1]|0x20) == 'l' &&
		(name[2]|0x20) == 'a' &&
		(name[3]|0x20) == 's' &&
		(name[4]|0x20) == 's'
}

var spaceTable = [256]byte{
	' ':  1,
	'\n': 1,
	'\r': 1,
	'\t': 1,
	'\f': 1,
}

func isSpace(b byte) bool { return spaceTable[b] != 0 }

func isAllSpace(b []byte) bool {
	for _, c := range b {
		if !isSpace(c) {
			return false
		}
	}
	return true
}

func trimASCIISpaceBoundsBytes(b []byte) (start int, end int) {
	start = 0
	end = len(b)
	for start < end && isSpace(b[start]) {
		start++
	}
	for end > start && isSpace(b[end-1]) {
		end--
	}
	return
}

func trimASCIISpaceBoundsString(s string) (start int, end int) {
	start = 0
	end = len(s)
	for start < end && isSpace(s[start]) {
		start++
	}
	for end > start && isSpace(s[end-1]) {
		end--
	}
	return
}

func equalASCIIFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if (a[i] | 0x20) != (b[i] | 0x20) {
			return false
		}
	}
	return true
}

func equalASCIIFoldBytesString(a []byte, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if (a[i] | 0x20) != (b[i] | 0x20) {
			return false
		}
	}
	return true
}

func eqASCIIFoldLit(s, lit string) bool {
	if len(s) != len(lit) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if (s[i] | 0x20) != lit[i] {
			return false
		}
	}
	return true
}
