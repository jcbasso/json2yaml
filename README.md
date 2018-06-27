Json2Yaml is a json parser using yacc.

To build it:

	$ go generate
	$ go build

or

	$ go generate
	$ go run expr.go

Then run, for example 
	$ echo "{\"a\":true}" | go run expr.go


The file main.go contains the "go generate" command to run yacc to
create expr.go from expr.y. It also has the package doc comment,
as godoc will not scan the .y file.

The actual implementation is in expr.y.