/*
   Environment

   An environment in this interpreter is what is used to keep track of values by associating them with a name.
   Under the hood, the environment is basically an hash map that associates strings with objects.
*/

package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

/*
   Enclosing Environments

   Here is a problem case, lets say in monna I would want to type this:

   ```
   let i = 5;
   let print_num = fn(i) {
      puts(i);
   }

   print_num(10);
   puts(i);
   ```

  The ideal result of the above code in the monna programming language is for 10 and 5 to be the outputs respectively.
  In a situation where enclosed environment does not exists, both outputs will be 10 because the current value of i
  would be overwritten. The ideal situation would be to preserve the previous binding to 'i' while also making a a new
  one.

  This works be creating a new instance of object.Environment with a pointer to the environment it should extend, doing this
  encloses a fresh and empty environment with an existing one. When the Get method is called and it itself doesn't have the value
  associated with the given name, it calls the Get of the enclosing environment. That's the environment it's extending. If that
  enclosing environment can't find the value, it calls its own enclosing environment and so on until there is no enclosing environment
  anymore and it will error out to an unknown identifier.
*/
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (l_environment *Environment) Get(name string) (Object, bool) {
	obj, ok := l_environment.store[name]

	if !ok && l_environment.outer != nil {
		obj, ok = l_environment.outer.Get(name)
	}

	return obj, ok
}

func (l_environment *Environment) Set(name string, value Object) Object {
	l_environment.store[name] = value
	return value
}
