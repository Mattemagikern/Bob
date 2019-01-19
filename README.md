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

When Bob is executed it will try to find the BUILDER (./BUILDER) file. It
contains the **recepies**, **variabels** and the **build command**. 
## Builder file
#### Variabels 
```make
<variable_name> = <expression>
$<variable_name> += <expression>
$<variable_name> -= <expression>
<variable_name> = $(<shell_command>)
```
Just like in Make variables consists of a name and an expression separated by a
``=,+=, -=``. 

The ``=`` assigns the expression on the right hand side to a variable with the
variable name on the left side.


The ``+=`` appends the expression on the right hand side to the expession of the
variable on the left hand side.


The ``-=`` removes the expression on the right hand side from the expession of the
variable on the left hand side.

The ``$(<shell_command>)`` assigns the output from the shell command to the
variable name. It can be used in congecture with the other operations (``+=``,
``-=``, ``=``).

To reference a variable prepend a ``$`` to the variable name. A variable name
may not contain space or tabs.

##### Recepies
```make
<recepie>: <ingredient> <ingredient> ...
	commands
```
A recepie consists of a list of ingredients and a list of commands. 

All ingredients is a reference to another recepie. If a recepie only consists of
commands it is called a root recepie. 


The root recipes are very use full since Bob remembers the previous builds and
dosn't reset the variables to its initial state between input recepies( command
line arguments). This enables Bob to work along the line of functional
programing and chaining in particular. This feature was the secondary motivator
for the creation of Bob. In large projects, makefiles may be enourmosly long and
hard to desipher while if we apply this way of thinking we can vastly reduce the
size of the make files and chain commands together to create unique builds for
testing or debugging purpuses. Since we now can supply the debug flag to a small
sub project we can filterout the intresting logs which makes expensive logging
tools uneccesarry. 


The commands that belonge to a specific recepie starts with a \t directly under
the recepie. The commands may be an update of a variable or a shell command.
#### The build recepie
```make
obj/%.o:%.c
	$CC $CFLAGS -c -o $@ $< 
```
The build recepie is slightly different from a regular recepie. It consists of a
placement of the object file: ``obj/%.o:`` and the suffix of the source files:
``%.c``. The examples provided is for a C project but change the suffix to the
appropriate language and it will work as inteded. 

# A full example of a BUILDER file
An example of an BUILDER file for a C project of mine:
```
inc = .*\/inc\/(.*)\.h$
src = .*\/src\/(.*)\.c$
inc_pattern = (?m)(?:^#include[\s]*)[<|"](.*)[>|"]

CC      = gcc
dep     =  ./inc
CFLAGS   = -Wall -Wextra -std=c99 -Wno-format -Wno-parentheses -Wno-empty-body
$CFLAGS  += $dep

all: build debug build server client
	echo $LA

server:
	$CC $CFLAGS -o server $Objects src/server.c

client:
	$CC $CFLAGS -o client $Objects src/client.c

debug:
	echo DEBUG
	src = .*tests/.*\.c$
	$LA = $(ls -la)
	$CFLAGS += -DDEBUG

obj/%.o:%.c
	$CC $CFLAGS -c -o $@ $< 
```
