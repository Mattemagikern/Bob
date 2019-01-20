# Bob - The builder

Bob were created to keep track of the current status of my C projects. 
The problem with conventional build systems like Make is that it is
**Stateless**, it doesn't care about previous build states, it will only check
weather or not the \*.c timestamp is less than the corresponding object file.
This creates the problem of include dependencies. If I were to change a \*.h in
the project, make would not update the object due to it hasn't checked the \*.c
files dependencies.

Bob on the other hand is **statefull** it will remember how it build your \*.o
files and map the corresponding \*.h files to \*.c files and see if the object
should update. This will eliminate the use of `make clean` for example. If you
are unsure or just want to clean you may do so, bob will truncate the .state
file and rebuild all files accordingly.

When Bob is executed it will try to find the Blueprint (./Blueprint) file. It
contains the **recepies**, **variabels** and the **build command**. 
## The Blueprint
#### Variables 
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


The ``+=`` appends the expression on the right hand side to the expression of the
variable on the left hand side.


The ``-=`` removes the expression on the right hand side from the expression of the
variable on the left hand side.

The ``$(<shell_command>)`` assigns the output from the shell command to the
variable name. It can be used in conjecture with the other operations (``+=``,
``-=``, ``=``).

To reference a variable prepend a ``$`` to the variable name. A variable name
may not contain space or tabs.

##### Special variables
In order for Bob to find and to be able to find and map all dependencies you need
to provide 3 regex patterns:

* src - Provides a regex pattern that finds the source files (\*.c) you'd like the
  builder command to execute.
* inc - Provides a regex pattern that finds the include files (\*.h) your source
  files depend upon.
* inc\_pattern - Provides a regex pattern for how a include is declared in the
  source or include files.

You only need to set ``src`` but to use all Bobs features I would advise you to
apply the other two regex patterns as well. 
##### Recipes
```make
<recepie>: <ingredient> <ingredient> ...
	commands
```
A recipe consists of a list of ingredients and a list of commands. 

All ingredients is a reference to another recipe. If a recipe only consists of
commands it is called a root recipe. 


The root recipes are very use full since Bob remembers the previous builds and
doesn't reset the variables to its initial state between input recipes( command
line arguments). This enables Bob to work along the line of functional
programing and chaining in particular. This feature was the secondary motivator
for the creation of Bob. In large projects, makefiles may be long and
hard to decipher while if we apply this way of thinking we can vastly reduce the
size of the make files and chain commands together to create unique builds for
testing or debugging purposes. Since we now can supply the debug flag to a small
sub project we can filter out the interesting logs which makes expensive logging
tools unnecessary. 


The commands that belong to a specific recipe starts with a \t directly under
the recipe. The commands may be an update of a variable or a shell command.
#### The build recepie
```make
obj/%.o:%.c
	$CC $CFLAGS -c -o $@ $< 
```
The build recipe is slightly different from a regular recipe. It consists of a
placement of the object file: ``obj/%.o:`` and the suffix of the source files:
``%.c``. The examples provided is for a C project but change the suffix to the
appropriate language and it will work as intended. 

# A full example of a Blueprint
An example of an Blueprint for a C project of mine:
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
