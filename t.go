package snoc

type Any interface{}

type Env struct {
	Chain *Pair
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
	F    func(args []Any, env Env) Any // args are evaluated.
}

type Special struct {
	Name string
	F    func(args []Any, env Env) Any // args are unevaluated.
}
