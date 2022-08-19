/*
   Environment

   An environment in this interpreter is what is used to keep track of values by associating them with a name.
   Under the hood, the environment is basically an hash map that associates strings with objects.
*/

package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (l_environment *Environment) Get(name string) (Object, bool) {
	obj, ok := l_environment.store[name]
	return obj, ok
}

func (l_environment *Environment) Set(name string, value Object) Object {
	l_environment.store[name] = value
	return value
}
