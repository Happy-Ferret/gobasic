// Package builtin contains the implementation for the core functions
// implemented for BASIC, such as SIN, COS, RND, PRINT, etc.
//
// The builtin package also provides an interface which with you can
// extent the interpreter to supply new primitives, for example you
// might implement RANDOM, PEEK, POKE, or something entirely different.
//
// Examples of extending the interpreter exist within the `goserver/`
// and `embed/` packages within the distribution.
//
package builtin

import (
	"bufio"
	"sync"

	"github.com/skx/gobasic/object"
)

// Signature is the signature of a builtin-function.
//
// Each built-in will receive an array of objects, and will return a
// single object back to the caller.
//
// In the case of an error then the object will be an error-object.
type Signature func(env Environment, args []object.Object) object.Object

// Builtins holds our state.
type Builtins struct {
	// lock holds a mutex to prevent corruption.
	lock sync.Mutex

	// argRegistry holds the number of arguments the given name requires.
	argRegistry map[string]int

	// fnRegistry holds a reference to the golang function which
	// implements the builtin.
	fnRegistry map[string]Signature
}

// New returns a new helper/holder for builtin functions.
func New() *Builtins {
	t := &Builtins{}
	t.argRegistry = make(map[string]int)
	t.fnRegistry = make(map[string]Signature)

	return t
}

// Register records a built-in function.
// The three arguments are:
//  NAME  - The thing that the BASIC program will call
//  nARGS - The number of arguments the built-in requires.
//          NOTE: Arguments are comma-separated in the BASIC program,
//          but commas are stripped out.
//  FT    - The function which provides the implementation.
func (b *Builtins) Register(name string, nArgs int, ft Signature) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// Record the details.
	b.argRegistry[name] = nArgs
	b.fnRegistry[name] = ft
}

// Get the values associated with the given built-in.
func (b *Builtins) Get(name string) (int, Signature) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.argRegistry[name], b.fnRegistry[name]
}

// Environment is an interface which is passed to all built-in functions.
type Environment interface {
	// StdInput is a handle to a reader-object, allowing input to
	// be processed by the built-ins.
	StdInput() *bufio.Reader

	// StdOutput is a handle to a writer-object, allowing output to
	// be generated by the built-ins.
	StdOutput() *bufio.Writer

	// StdError is a handle to a writer-object, allowing user errors to
	// be communicated by the built-ins.
	StdError() *bufio.Writer

	// LineEnding specifies any additional characters that should be
	// appended to PRINT commands - for example '\n'
	LineEnding() string

	// Data allows the builtins to get a reference to the intepreter.
	Data() interface{}
}
