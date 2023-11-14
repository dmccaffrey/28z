# About
28z is a simple VM loosely inspired by the RPL environment on the HP-28 calculator.
![](ui.png)

# Building
go build .

# Running
./28z

## Data types

### Floating point
All values are assumed to be floating point by default.

### String
String values are identified by a single preceeding quotation mark (').

### Sequence
Sequence values contain a sequence of instructions created dynamically through the define and reduce sequence instructions.

### Reference
Reference values are identified by a preceeding dollar sign ($) and are replaced with the corresponding register or variable value either when input, or in storing mode when evaluated.

References can also be specified with a preceeding percent sign (%) in interactive input to have the reference be resolved and immediately evaluated.

## Supported instructions

### produce
- Description: Pop from this stack and push to previous
- Arg count: 1
- Result count: 0
- Usage: produce ⤶

### purge
- Description: Purge x
- Arg count: 1
- Result count: 0
- Usage: 'a ⤶ purge ⤶ undefined⥗a

### inspect
- Description: Write a raw object to file
- Arg count: 1
- Result count: 0
- Usage: 

### sin
- Description: sin of x
- Arg count: 1
- Result count: 1
- Usage: 

### collect
- Description: Collect stack into x
- Arg count: 1
- Result count: 1
- Usage: 1 ⤶ 2 ⤶ collect ⤶ ⤒[2]:1,2

### put
- Description: Put y into x
- Arg count: 2
- Result count: 0
- Usage: 2 ⤶ 'a ⤶ store ⤶ y⥗a

### clear
- Description: Clear stack
- Arg count: 0
- Result count: 0
- Usage: clear ⤶

### print
- Description: Print x
- Arg count: 1
- Result count: 0
- Usage: 'Hello world ⤶ print ⤶ Hello world⥱Console

### clearbuf
- Description: Clear the output buffer
- Arg count: 0
- Result count: 0
- Usage: 

### <
- Description: Define sequence
- Arg count: 0
- Result count: 0
- Usage: < ⤶

### recall
- Description: Recall x
- Arg count: 1
- Result count: 1
- Usage: 'a ⤶ recall ⤶ ⤒a

### prompt
- Description: Prompt the user for a value
- Arg count: 1
- Result count: 1
- Usage: 'Enter x ⤶ prompt ⤶

### mmap
- Description: Map a file to RAM
- Arg count: 1
- Result count: 0
- Usage: 'rom/file.raw ⤶ mmap ⤶ file.byes⥱RAM

### *
- Description: Multiply y by x
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ * ⤶ ⤒12

### drop
- Description: Drop x
- Arg count: 1
- Result count: 0
- Usage: drop ⤶

### ==
- Description: Set the result flag to 1 if x = y
- Arg count: 2
- Result count: 0
- Usage: 

### /
- Description: Divide y by x
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ / ⤶ ⤒3

### >
- Description: Define sequence
- Arg count: 0
- Result count: 0
- Usage: > ⤶

### end
- Description: Return from function
- Arg count: 0
- Result count: 0
- Usage: end ⤶

### store
- Description: Store y into x
- Arg count: 2
- Result count: 1
- Usage: 2 ⤶ 'a ⤶ store ⤶ ⤒2; y⥗a

### loop
- Description: Execute x if the loop counter is not zero
- Arg count: 0
- Result count: 0
- Usage: 5 ⤶ setloop ⤶ ⤒<sequence> | loop ⤶

### ceval
- Description: Conditionally evaluate x if result flag is 1
- Arg count: 1
- Result count: 0
- Usage: ⤒<sequence> | ceval ⤶

### +
- Description: Add x and y
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ + ⤶ ⤒8

### <=
- Description: Set the result flag to 1 if y <= x
- Arg count: 2
- Result count: 0
- Usage: 

### >=
- Description: Set the result flag to 1 if y >= x
- Arg count: 2
- Result count: 0
- Usage: 

### inverse
- Description: Inverts x
- Arg count: 1
- Result count: 1
- Usage: 

### consume
- Description: Pop from previous stack and push to current
- Arg count: 0
- Result count: 1
- Usage: consume ⤶

### swap
- Description: Swap x and y
- Arg count: 2
- Result count: 2
- Usage: swap ⤶ ⤒x,y

### -
- Description: Subtract x from y
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ - ⤶ ⤒4

### enter
- Description: Enter function
- Arg count: 0
- Result count: 0
- Usage: enter ⤶

### show
- Description: Render and pause
- Arg count: 0
- Result count: 0
- Usage: 

### zero
- Description: Zero RAM
- Arg count: 0
- Result count: 0
- Usage: 

### graph
- Description: Graph a sequence
- Arg count: 3
- Result count: 0
- Usage: graph ⤶

### repeat
- Description: Execute x repeatedly
- Arg count: 4
- Result count: 0
- Usage: 0 ⤶ < ⤶'f ⤶ repeat ⤶

### dec
- Description: Decrement the loop register
- Arg count: 0
- Result count: 0
- Usage: dec

### cos
- Description: cos of x
- Arg count: 1
- Result count: 1
- Usage: 

### this
- Description: Refer to the current sequence
- Arg count: 0
- Result count: 1
- Usage: this ⤶

### exchange
- Description: Exchange y and the value in var x
- Arg count: 2
- Result count: 1
- Usage: 3 ⤶ 'a ⤶ exchange ⤶ ⤒a 3⥗a

### expand
- Description: Expand x into the stack
- Arg count: 1
- Result count: -1
- Usage: ⤒[2]:1,2 | expand ⤶ ⤒1 ⤒2

### render
- Description: Render RAM as buffer
- Arg count: 0
- Result count: 0
- Usage: render ⤶

### stream
- Description: Apply x to renderable RAM
- Arg count: 1
- Result count: 0
- Usage: 

### unset
- Description: Sets the result flat to 0
- Arg count: 0
- Result count: 0
- Usage: 

### halt
- Description: Halt execution
- Arg count: 0
- Result count: 0
- Usage: halt ⤶

### eval
- Description: Evaluate x
- Arg count: 1
- Result count: 0
- Usage: 

### rand
- Description: Generate random between 0 and 1
- Arg count: 0
- Result count: 1
- Usage: 

### setloop
- Description: Set loop counter to x
- Arg count: 1
- Result count: 0
- Usage: 5 ⤶ setloop ⤶

### sleep
- Description: Sleep for x ms
- Arg count: 1
- Result count: 0
- Usage: 

### status
- Description: Display status
- Arg count: 0
- Result count: 0
- Usage: 

### mod
- Description: y modulus by x
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ / ⤶ ⤒0

### apply
- Description: Evalue x against all entries in y
- Arg count: 2
- Result count: 1
- Usage: apply ⤶

### reduce
- Description: Use x to reduce y to a single value
- Arg count: 2
- Result count: 1
- Usage: reduce ⤶

### files
- Description: List availabel files in ROM
- Arg count: 0
- Result count: 0
- Usage: files ⤶ [files]⥱Console

### !=
- Description: Set the result flag to 1 if x != y
- Arg count: 2
- Result count: 0
- Usage: 

