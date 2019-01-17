# Bob - The builder


Bob were created to keep track of the current status of my C projects. 
The problem with conventional build systems like Make is that it is
**Stateless**, it dosn't care about previous build states, it will only check
wheather or not the \*.c timestamp is less than the corresponding object file.
This creates the problem of include dependencies. If I were to change a \*.h in
the project make would not update the object due to it hasn't checked the \*.c
files dependencies.

Bob on the other hand is **statefull** it will remember how it build your \*.o
files and mapp the corresponding \*.h files to \*.c files and see if the object
should update. This will eliminate the use of `make clean` for example. If you
are unsure or just want to clean you may do so, bob will turncate the .state
file and rebuild all files accordingly.

## Builder file
```make
<variable_name> = <expression>

<recepie>: <ingredients>
	commands
```
Build command:
```make
obj/%.o:%.c
	$CC $CFLAGS -c -o $@ $< 
```
```make
inc = .*\/inc\/(.*)\.h$
src = .*\/src\/(.*)\.c$
inc_pattern = (?m)(?:^#include[\s]*)[<|"](.*)[>|"]
```
