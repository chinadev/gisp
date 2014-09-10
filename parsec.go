package gisp

import (
	"strconv"

	p "github.com/Dwarfartisan/goparsec"
)

// IntParser 解析整数
func IntParser(st p.ParseState) (interface{}, error) {
	i, err := p.Int(st)
	if err == nil {
		val, err := strconv.Atoi(i.(string))
		if err == nil {
			return Int(val), nil
		}
		return nil, err
	}
	return nil, err

}

var EscapeChar = p.Bind_(p.Rune('\\'), func(st p.ParseState) (interface{}, error) {
	r, err := p.OneOf("nrt\"\\")(st)
	if err == nil {
		ru := r.(rune)
		switch ru {
		case 'r':
			return '\r', nil
		case 'n':
			return '\n', nil
		// FIXME:引号的解析偷懒了，单双引号的应该分开。
		case '\'':
			return '\'', nil
		case '"':
			return '"', nil
		case '\\':
			return '\\', nil
		case 't':
			return '\t', nil
		default:
			return nil, st.Trap("Unknown escape sequence \\%c", r)
		}
	} else {
		return nil, err
	}
})

var RuneParser = p.Bind(
	p.Between(p.Rune('\''), p.Rune('\''),
		p.Either(p.Try(EscapeChar), p.NoneOf("'"))),
	p.ReturnString)

var StringParser = p.Bind(
	p.Between(p.Rune('"'), p.Rune('"'),
		p.Many(p.Either(p.Try(EscapeChar), p.NoneOf("\"")))),
	p.ReturnString)

func bodyParser(st p.ParseState) (interface{}, error) {
	value, err := p.SepBy(ValueParser, p.Many1(p.Space))(st)
	return value, err
}

func ListParser(st p.ParseState) (interface{}, error) {
	one := p.Bind(AtomParser, func(atom interface{}) p.Parser {
		return p.Bind_(p.Rune(')'), p.Return(List{atom}))
	})
	list, err := p.Either(p.Try(p.Bind_(p.Rune('('), one)),
		p.Between(p.Rune('('), p.Rune(')'), bodyParser))(st)
	if err == nil {
		return List(list.([]interface{})), nil
	} else {
		return nil, err
	}
}

func QuoteParser(st p.ParseState) (interface{}, error) {
	lisp, err := p.Bind_(p.Rune('\''), ValueParser)(st)
	if err == nil {
		return Quote{lisp}, nil
	} else {
		return nil, err
	}
}

func ValueParser(st p.ParseState) (interface{}, error) {
	value, err := p.Choice(p.Try(StringParser),
		p.Try(FloatParser),
		p.Try(IntParser),
		p.Try(QuoteParser),
		p.Try(RuneParser),
		p.Try(StringParser),
		p.Try(BoolParser),
		p.Try(NilParser),
		// p.Try(AtomParser),
		p.Try(p.Bind(AtomParser, SuffixParser)),
		// ListParser,
		p.Bind(ListParser, SuffixParser),
	)(st)
	return value, err
}
