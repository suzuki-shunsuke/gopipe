package gopipe

type Arg struct {
	Name  string
	Value interface{}
}

type Args struct {
	Args map[string]Arg
}

type GetArgs func(mArgs map[string]*Args) (*Args, error)

func (args *Args) Set(name string, value interface{}) {
	if args.Args == nil {
		args.Args = map[string]Arg{}
	}
	args.Args[name] = Arg{
		Name:  name,
		Value: value,
	}
}

func (args *Args) Get(name string) interface{} {
	if args.Args == nil {
		return nil
	}
	return args.Args[name].Value
}

func (args *Args) GetOk(name string) (interface{}, bool) {
	if args.Args == nil {
		return nil, false
	}
	arg, ok := args.Args[name]
	if ok {
		return arg.Value, true
	}
	return nil, false
}

func (args *Args) GetString(name string) string {
	if args.Args == nil {
		return ""
	}
	return args.Args[name].Value.(string)
}

func (args *Args) GetStringOk(name string) (string, bool) {
	if args.Args == nil {
		return "", false
	}
	arg, ok := args.Args[name]
	if ok {
		return arg.Value.(string), true
	}
	return "", false
}
