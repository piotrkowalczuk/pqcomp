package pqcomp

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/piotrkowalczuk/nilt"
)

const (
	// E represents equal operator.
	E = "="
	// NE represents not equal operator.
	NE = "<>"
	// GT represents greater than operator.
	GT = ">"
	// LT represents lower than operator.
	LT = "<"
	// GTE represents greater than or equal operator.
	GTE = ">="
	// LTE represents lower than or equal operator.
	LTE = "<="
	// DESC represents descendant way of sorting.
	DESC = "DESC"
	// ASC represents ascendant way of sorting.
	ASC = "ASC"
)

// Appearer wraps Appear function.
type Appearer interface {
	// Appear returns true if object should be used by AddExpr method.
	// Otherwise it is ignored.
	Appear() bool
}

// Composer is some sort of stateful iterator that helps to build complex SQL queries.
type Composer struct {
	composed        int
	keys, operators []string
	arguments       []interface{}
	idx, diff       int
	parent          *Composer
	childs          []*Composer
}

// New allocates new Composer and pre-allocates space for given amount of arguments and expressions.
// Each child expression passed to the constructor creates new child Composer
// with pre-allocated space for expressions.
func New(arguments, nbOfExpressions int, nbOfChildExpressions ...int) *Composer {
	return neww(nil, arguments, nbOfExpressions, nbOfChildExpressions...)
}

func neww(parent *Composer, args, pexpr int, cexpr ...int) *Composer {
	comp := Composer{
		keys:      make([]string, 0, pexpr),
		operators: make([]string, 0, pexpr),
		arguments: make([]interface{}, 0, args),
		diff:      args,
		childs:    make([]*Composer, len(cexpr)),
		parent:    parent,
	}

	for i := 0; i < len(cexpr); i++ {
		comp.childs[i] = neww(&comp, cexpr[i], cexpr[i])
	}
	return &comp
}

// AddArg add static argument.
func (c *Composer) AddArg(arg interface{}) {
	c.arguments = append(c.arguments, arg)
}

// AddExpr adds expression if value meet certain requirements.
// To know more please read the source code.
func (c *Composer) AddExpr(key, operator string, value interface{}) {
	if value == nil {
		return
	}

	switch v := value.(type) {
	case Appearer:
		if v.Appear() {
			c.addExpr(key, operator, value)
		}
	case *sql.NullBool:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case *nilt.Bool:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case nilt.Bool:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case *sql.NullString:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case *nilt.String:
		fmt.Println("*nilt.String WTF", v.Valid)
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case nilt.String:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case *sql.NullInt64:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case *nilt.Int64:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case nilt.Int64:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case *nilt.Int32:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case nilt.Int32:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case *sql.NullFloat64:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case *nilt.Float64:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	case nilt.Float64:
		if v.Valid {
			c.addExpr(key, operator, value)
		}
	default:
		vo := reflect.ValueOf(v)
		switch vo.Kind() {
		case reflect.Slice:
			if !vo.IsNil() {
				c.addExpr(key, operator, value)
			}
		default:
			c.addExpr(key, operator, value)
		}
	}
}

func (c *Composer) addExpr(key, expr string, value interface{}) {
	c.keys = append(c.keys, key)
	c.operators = append(c.operators, expr)
	c.arguments = append(c.arguments, value)
}

// Compose returns next available composer
// or if pool of pre-allocated Composer's is empty allocates new one.
func (c *Composer) Compose(nbOfChildExpressions ...int) (comp *Composer) {
	if len(c.childs) > c.composed {
		comp = c.childs[c.composed]

		if len(nbOfChildExpressions) != 0 {
			comp.childs = make([]*Composer, 0, len(nbOfChildExpressions))
			for i := range comp.childs {
				comp.childs[i] = neww(c, nbOfChildExpressions[i], nbOfChildExpressions[i])
			}
		}
		c.composed++
		return
	}

	comp = neww(c, 0, 0, nbOfChildExpressions...)
	c.childs = append(c.childs, comp)
	c.composed++
	return
}

// Args returns slice of arguments that was passed to the composer
// or to any child.
func (c *Composer) Args() []interface{} {
	if len(c.childs) == 0 {
		return c.arguments
	}

	args := make([]interface{}, 0, c.lenWithChilds())
	args = append(args, c.arguments...)
	for _, ch := range c.childs {
		args = append(args, ch.arguments...)
	}

	return args
}

func (c *Composer) lenWithChilds() (count int) {
	for _, ch := range c.childs {
		count += len(ch.arguments)
	}

	return
}

// Len returns number of expressions that was passed.
func (b *Composer) Len() int {
	return len(b.keys)
}

// Next move cursor to next position. Returns false if it's not possible.
func (b *Composer) Next() bool {
	if b.idx < b.Len() {
		b.idx++
		if b.parent != nil {
			b.parent.idx++
		}
		return true
	}

	return false
}

// Reset set cursor back to 0
func (b *Composer) Reset() {
	b.idx = 0
}

// Key returns key for current cursor position.
func (b *Composer) Key() string {
	return b.keys[b.idx-1]
}

// Oper returns operator for current cursor position.
func (b *Composer) Oper() string {
	return b.operators[b.idx-1]
}

// PlaceHolder returns placeholder for current cursor position.
func (b *Composer) PlaceHolder() string {
	if b.parent != nil {
		return b.parent.PlaceHolder()
	}
	return "$" + strconv.FormatInt(int64(b.diff+b.idx), 10)
}

// First returns true cursor is on first position.
func (b *Composer) First() bool {
	return b.idx == 1
}
