package snoc

type Any interface{}

type Terp struct {
	Globals map[*Sym]Any
}

type Env struct {
	Proto *ProtoFunc
	Up    *Env
	Slots []Any
	Terp  *Terp
}

type ProtoFunc struct {
	Outer  *ProtoFunc
	Params []*Sym
	Values []Any // Only for Let
	Body   Any
	Name   string
	IsLet  bool
}

type Func struct {
	Outer  *Env
	Params []*Sym
	Values []Any // Only for Let
	Body   Any
	Name   string
	IsLet  bool
	Proto  *ProtoFunc
}

type Var struct {
	Proto *ProtoFunc
	Slot  int
	Sym   *Sym
}

type Sym struct {
	S string
}

type Pair struct {
	H Any
	T *Pair
}

type Prim struct {
	Name string
	F    func(args []Any, env *Env) Any // args are evaluated.
}

type Special struct {
	Name string
	F    func(args []Any, env *Env) Any // args are unevaluated.
}
