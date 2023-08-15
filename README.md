## Monk Programming Language
This project is a result of learning the innards of compilers. This is an interpreter but there are really not so different other than not so minor differences :) The book in questinon [Writing an Interpreter In Go](https://interpreterbook.com/).

![Demo of monkey working](/doc/demo.png)

### Using
To use the monkey interpreter, simply run `go run main.go` in the project directory.To get an executable, run `go build` in the project directory.

### Features
#### Variables:
	let x = 5;

#### Return Statements:
	return 5; 
	return false;

#### Expressions: 
	if (x < y) { x }

	if (x > y) { y } else { x }

#### Functions: 
	fn(x, y) { x + y; }

#### String Literals:
	"Hello World"

#### Built-in Functions:
- **len()**: Returns the number of characters in a string.
```
len("Dexter's Laboratory");
19
```
- **puts()**: Prints the given arguments to STDOUT. It return a null. It only cares about printing not returning a value.
```
puts("Sugar, spice, and everything nice!");
Sugar, spice, and everything nice!

```
