package libxml2

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const stdXmlDecl = `<?xml version="1.0"?>` + "\n"

var (
	goodWFNSStrings = []string{
		stdXmlDecl + `<foobar xmlns:bar="xml://foo" bar:foo="bar"/>` + "\n",
		stdXmlDecl + `<foobar xmlns="xml://foo" foo="bar"><foo/></foobar>` + "\n",
		stdXmlDecl + `<bar:foobar xmlns:bar="xml://foo" foo="bar"><foo/></bar:foobar>` + "\n",
		stdXmlDecl + `<bar:foobar xmlns:bar="xml://foo" foo="bar"><bar:foo/></bar:foobar>` + "\n",
		stdXmlDecl + `<bar:foobar xmlns:bar="xml://foo" bar:foo="bar"><bar:foo/></bar:foobar>` + "\n",
	}
	goodWFStrings = []string{
		`<foobar/>`,
		`<foobar></foobar>`,
		`<foobar></foobar>`,
		`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<foobar></foobar>`,
		`<?xml version="1.0" encoding="ISO-8859-1"?>` + "\n" + `<foobar></foobar>`,
		stdXmlDecl + `<foobar> </foobar>` + "\n",
		stdXmlDecl + `<foobar><foo/></foobar> `,
		stdXmlDecl + `<foobar> <foo/> </foobar> `,
		stdXmlDecl + `<foobar><![CDATA[<>&"\` + "`" + `]]></foobar>`,
		stdXmlDecl + `<foobar>&lt;&gt;&amp;&quot;&apos;</foobar>`,
		stdXmlDecl + `<foobar>&#x20;&#160;</foobar>`,
		stdXmlDecl + `<!--comment--><foobar>foo</foobar>`,
		stdXmlDecl + `<foobar>foo</foobar><!--comment-->`,
		stdXmlDecl + `<foobar>foo<!----></foobar>`,
		stdXmlDecl + `<foobar foo="bar"/>`,
		stdXmlDecl + `<foobar foo="\` + "`" + `bar>"/>`,
	}
	goodWFDTDStrings = []string{
		stdXmlDecl + `<!DOCTYPE foobar [` + "\n" + `<!ENTITY foo " test ">` + "\n" + `]>` + "\n" + `<foobar>&foo;</foobar>`,
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar>&foo;</foobar>`,
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar>&foo;&gt;</foobar>`,
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar=&quot;foo&quot;">]><foobar>&foo;&gt;</foobar>`,
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar>&foo;&gt;</foobar>`,
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar foo="&foo;"/>`,
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar foo="&gt;&foo;"/>`,
	}
	badWFStrings = []string{
		"",                                      // totally empty document
		stdXmlDecl,                              // only XML Declaration
		"<!--ouch-->",                           // comment only is like an empty document
		`<!DOCTYPE ouch [<!ENTITY foo "bar">]>`, // no good either ...
		"<ouch>",                // single tag (tag mismatch)
		"<ouch/>foo",            // trailing junk
		"foo<ouch/>",            // leading junk
		"<ouch foo=bar/>",       // bad attribute
		`<ouch foo="bar/>`,      // bad attribute
		"<ouch>&</ouch>",        // bad char
		`<ouch>&//0x20;</ouch>`, // bad chart
		"<foob<e4>r/>",          // bad encoding
		"<ouch>&foo;</ouch>",    // undefind entity
		"<ouch>&gt</ouch>",      // unterminated entity
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar &foo;="ouch"/>`,          // bad placed entity
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar=&quot;foo&quot;">]><foobar &foo;/>`, // even worse
		"<ouch><!---></ouch>",   // bad comment
		"<ouch><!-----></ouch>", // bad either... (is this conform with the spec????)
	}
)

func parseShouldSucceed(t *testing.T, opts ParseOption, inputs []string) {
	t.Logf("Test parsing with parser %v", opts)
	for _, s := range inputs {
		d, err := ParseString(s, opts)
		if !assert.NoError(t, err, "Parse should succeed") {
			return
		}
		d.Free()
	}
}

func parseShouldFail(t *testing.T, opts ParseOption, inputs []string) {
	for _, s := range inputs {
		d, err := ParseString(s, opts)
		if err == nil {
			d.Free()
			t.Errorf("Expected failure to parse '%s'", s)
		}
	}
}

type ParseOptionToString struct {
	v ParseOption
	e string
}

func TestParseOptionStringer(t *testing.T) {
	values := []ParseOptionToString{
		ParseOptionToString{
			v: XmlParseRecover,
			e: "Recover",
		},
		ParseOptionToString{
			v: XmlParseNoEnt,
			e: "NoEnt",
		},
		ParseOptionToString{
			v: XmlParseDTDLoad,
			e: "DTDLoad",
		},
		ParseOptionToString{
			v: XmlParseDTDAttr,
			e: "DTDAttr",
		},
		ParseOptionToString{
			v: XmlParseDTDValid,
			e: "DTDValid",
		},
		ParseOptionToString{
			v: XmlParseNoError,
			e: "NoError",
		},
		ParseOptionToString{
			v: XmlParseNoWarning,
			e: "NoWarning",
		},
		ParseOptionToString{
			v: XmlParsePedantic,
			e: "Pedantic",
		},
		ParseOptionToString{
			v: XmlParseNoBlanks,
			e: "NoBlanks",
		},
		ParseOptionToString{
			v: XmlParseSAX1,
			e: "SAX1",
		},
		ParseOptionToString{
			v: XmlParseXInclude,
			e: "XInclude",
		},
		ParseOptionToString{
			v: XmlParseNoNet,
			e: "NoNet",
		},
		ParseOptionToString{
			v: XmlParseNoDict,
			e: "NoDict",
		},
		ParseOptionToString{
			v: XmlParseNsclean,
			e: "Nsclean",
		},
		ParseOptionToString{
			v: XmlParseNoCDATA,
			e: "NoCDATA",
		},
		ParseOptionToString{
			v: XmlParseNoXIncNode,
			e: "NoXIncNode",
		},
		ParseOptionToString{
			v: XmlParseCompact,
			e: "Compact",
		},
		ParseOptionToString{
			v: XmlParseOld10,
			e: "Old10",
		},
		ParseOptionToString{
			v: XmlParseNoBaseFix,
			e: "NoBaseFix",
		},
		ParseOptionToString{
			v: XmlParseHuge,
			e: "Huge",
		},
		ParseOptionToString{
			v: XmlParseOldSAX,
			e: "OldSAX",
		},
		ParseOptionToString{
			v: XmlParseIgnoreEnc,
			e: "IgnoreEnc",
		},
		ParseOptionToString{
			v: XmlParseBigLines,
			e: "BigLines",
		},
	}

	for _, d := range values {
		if d.v.String() != "["+d.e+"]" {
			t.Errorf("e '%s', got '%s'", d.e, d.v.String())
		}
	}
}

func TestParseEmpty(t *testing.T) {
	doc, err := ParseString(``)
	if err == nil {
		t.Errorf("Parse of empty string should fail")
		defer doc.Free()
	}
}

func TestParse(t *testing.T) {
	inputs := [][]string{
		goodWFStrings,
		goodWFNSStrings,
		goodWFDTDStrings,
	}

	for _, input := range inputs {
		parseShouldSucceed(t, 0, input)
	}
}

func TestParseBad(t *testing.T) {
	inputs := [][]string{
		badWFStrings,
	}

	for _, input := range inputs {
		parseShouldFail(t, 0, input)
	}
}

func TestParseNoBlanks(t *testing.T) {
	inputs := [][]string{
		goodWFStrings,
		goodWFNSStrings,
		goodWFDTDStrings,
	}
	for _, input := range inputs {
		parseShouldSucceed(t, XmlParseNoBlanks, input)
	}
}

func TestRoundtripNoBlanks(t *testing.T) {
	doc, err := ParseString(`<a>    <b/> </a>`, XmlParseNoBlanks)
	if err != nil {
		t.Errorf("failed to parse string: %s", err)
		return
	}

	if !assert.Regexp(t, regexp.MustCompile(`<a><b/></a>`), doc.Dump(false), "stringified xml should have no blanks") {
		return
	}
}