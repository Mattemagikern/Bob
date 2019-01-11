# Bob - The builder

TODO:
!! refactor variable substitution code so the variable substitution and variable
declaration code can coexist and use the same codebase.
!! refactor variable declarations

Define the Build command
- Syntax (@? == filepath)
- mkdir -p objects
- Compile command over all files matching inc.Sf.src regex
- append to build state and write build state file

Define clean command
- turncate .state file
- continue

utils: 
- pre-commit hook formater (gofmt)

