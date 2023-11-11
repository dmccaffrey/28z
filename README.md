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

## Supported instructions

### !=
- Description: Set the result flag to 1 if x != y
- Arg count: 2
- Result count: 0
- Usage: 

### loop
- Description: Execute x if the loop counter is not zero
- Arg count: 0
- Result count: 0
- Usage: 5 ⤶ setloop ⤶ ⤒<sequence> | loop ⤶

### <
- Description: Define sequence
- Arg count: 0
- Result count: 0
- Usage: < ⤶

### produce
- Description: Pop from this stack and push to previous
- Arg count: 1
- Result count: 0
- Usage: produce ⤶

### recall
- Description: Recall x
- Arg count: 1
- Result count: 1
- Usage: 'a ⤶ recall ⤶ ⤒a

### render
- Description: Render RAM as buffer
- Arg count: 0
- Result count: 0
- Usage: render ⤶

### mmap
- Description: Map a file to RAM
- Arg count: 1
- Result count: 0
- Usage: 'rom/file.raw ⤶ mmap ⤶ file.byes⥱RAM

### ==
- Description: Set the result flag to 1 if x = y
- Arg count: 2
- Result count: 0
- Usage: 

### inverse
- Description: Inverts x
- Arg count: 1
- Result count: 1
- Usage: 

### exchange
- Description: Exchange y and the value in var x
- Arg count: 2
- Result count: 1
- Usage: 3 ⤶ 'a ⤶ exchange ⤶ ⤒a 3⥗a

### swap
- Description: Swap x and y
- Arg count: 2
- Result count: 2
- Usage: swap ⤶ ⤒x,y

### consume
- Description: Pop from previous stack and push to current
- Arg count: 0
- Result count: 1
- Usage: consume ⤶

### reduce
- Description: Use x to reduce y to a single value
- Arg count: 2
- Result count: 1
- Usage: reduce ⤶

### clear
- Description: Clear stack
- Arg count: 0
- Result count: 0
- Usage: clear ⤶

### graph
- Description: Graph a sequence
- Arg count: 3
- Result count: 0
- Usage: graph ⤶

### files
- Description: List availabel files in ROM
- Arg count: 0
- Result count: 0
- Usage: files ⤶ [files]⥱Console

### ceval
- Description: Conditionally evaluate x if result flag is 1
- Arg count: 1
- Result count: 0
- Usage: ⤒<sequence> | ceval ⤶

### *
- Description: Multiply y by x
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ * ⤶ ⤒12

### /
- Description: Divide y by x
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ / ⤶ ⤒3

### apply
- Description: Evalue x against all entries in y
- Arg count: 2
- Result count: 1
- Usage: apply ⤶

### end
- Description: Return from function
- Arg count: 0
- Result count: 0
- Usage: end ⤶

### purge
- Description: Purge x
- Arg count: 1
- Result count: 0
- Usage: 'a ⤶ purge ⤶ undefined⥗a

### drop
- Description: Drop x
- Arg count: 1
- Result count: 0
- Usage: drop ⤶

### >
- Description: Define sequence
- Arg count: 0
- Result count: 0
- Usage: > ⤶

### collect
- Description: Collect stack into x
- Arg count: 1
- Result count: 1
- Usage: 1 ⤶ 2 ⤶ collect ⤶ ⤒[2]:1,2

### <=
- Description: Set the result flag to 1 if y <= x
- Arg count: 2
- Result count: 0
- Usage: 

### +
- Description: Add x and y
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ + ⤶ ⤒8

### mod
- Description: y modulus by x
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ / ⤶ ⤒0

### store
- Description: Store y into x
- Arg count: 2
- Result count: 1
- Usage: 2 ⤶ 'a ⤶ store ⤶ ⤒2; y⥗a

### expand
- Description: Expand x into the stack
- Arg count: 1
- Result count: -1
- Usage: ⤒[2]:1,2 | expand ⤶ ⤒1 ⤒2

### status
- Description: Display status
- Arg count: 0
- Result count: 0
- Usage: 

### -
- Description: Subtract x from y
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ - ⤶ ⤒4

### setloop
- Description: Set loop counter to x
- Arg count: 1
- Result count: 0
- Usage: 5 ⤶ setloop ⤶

### dec
- Description: Decrement the loop register
- Arg count: 0
- Result count: 0
- Usage: dec

### halt
- Description: Halt execution
- Arg count: 0
- Result count: 0
- Usage: halt ⤶

### >=
- Description: Set the result flag to 1 if y >= x
- Arg count: 2
- Result count: 0
- Usage: 

### eval
- Description: Evaluate x
- Arg count: 1
- Result count: 0
- Usage: 

### this
- Description: Refer to the current sequence
- Arg count: 0
- Result count: 1
- Usage: this ⤶

### enter
- Description: Enter function
- Arg count: 0
- Result count: 0
- Usage: enter ⤶

### put
- Description: Put y into x
- Arg count: 2
- Result count: 0
- Usage: 2 ⤶ 'a ⤶ store ⤶ y⥗a

### print
- Description: Print x
- Arg count: 1
- Result count: 0
- Usage: 'Hello world ⤶ print ⤶ Hello world⥱Console

### repeat
- Description: Execute x repeatedly
- Arg count: 4
- Result count: 0
- Usage: 0 ⤶ < ⤶'f ⤶ repeat ⤶

