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

### mod
- Description: y modulus by x
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ / ⤶ ⤒0

### enter
- Description: Enter function
- Arg count: 0
- Result count: 0
- Usage: 

### store
- Description: Store y into x
- Arg count: 2
- Result count: 1
- Usage: 2 ⤶ 'a ⤶ store ⤶ ⤒2; y⥗a

### drop
- Description: Drop x
- Arg count: 1
- Result count: 0
- Usage: drop ⤶

### >=
- Description: Set the result flag to 1 if y >= x
- Arg count: 2
- Result count: 0
- Usage: 

### status
- Description: Display status
- Arg count: 0
- Result count: 0
- Usage: 

### exchange
- Description: Exchange y and the value in var x
- Arg count: 2
- Result count: 1
- Usage: 3 ⤶ 'a ⤶ exchange ⤶ ⤒a 3⥗a

### purge
- Description: Purge x
- Arg count: 1
- Result count: 0
- Usage: 'a ⤶ purge ⤶ undefined⥗a

### print
- Description: Print x
- Arg count: 1
- Result count: 0
- Usage: 'Hello world ⤶ print ⤶ Hello world⥱Console

### graph
- Description: Graph a sequence
- Arg count: 3
- Result count: 0
- Usage: 

### <=
- Description: Set the result flag to 1 if y <= x
- Arg count: 2
- Result count: 0
- Usage: 

### -
- Description: Subtract x from y
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ - ⤶ ⤒4

### put
- Description: Put y into x
- Arg count: 2
- Result count: 0
- Usage: 2 ⤶ 'a ⤶ store ⤶ y⥗a

### ==
- Description: Set the result flag to 1 if x = y
- Arg count: 2
- Result count: 0
- Usage: 

### ceval
- Description: Conditionally evaluate x if result flag is 1
- Arg count: 1
- Result count: 0
- Usage: 

### loop
- Description: Execute x if the loop counter is not zero
- Arg count: 0
- Result count: 0
- Usage: 

### >
- Description: Define sequence
- Arg count: 0
- Result count: 0
- Usage: 

### this
- Description: Refer to the current sequence
- Arg count: 0
- Result count: 1
- Usage: 

### produce
- Description: Pop from this stack and push to previous
- Arg count: 1
- Result count: 0
- Usage: 

### swap
- Description: Swap x and y
- Arg count: 2
- Result count: 2
- Usage: swap ⤶ ⤒x,y

### mmap
- Description: Map a file to RAM
- Arg count: 1
- Result count: 0
- Usage: 'rom/file.raw ⤶ mmap ⤶ file.byes⥱RAM

### collect
- Description: Collect stack into x
- Arg count: 1
- Result count: 1
- Usage: 

### recall
- Description: Recall x
- Arg count: 1
- Result count: 1
- Usage: 'a ⤶ recall ⤶ ⤒a

### +
- Description: Add x and y
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ + ⤶ ⤒8

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

### <
- Description: Define sequence
- Arg count: 0
- Result count: 0
- Usage: 

### render
- Description: Render RAM as buffer
- Arg count: 0
- Result count: 0
- Usage: 

### clear
- Description: Clear stack
- Arg count: 0
- Result count: 0
- Usage: clear ⤶

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

### eval
- Description: Evaluate x
- Arg count: 1
- Result count: 0
- Usage: 

### consume
- Description: Pop from previous stack and push to current
- Arg count: 0
- Result count: 1
- Usage: 

### end
- Description: Return from function
- Arg count: 0
- Result count: 0
- Usage: 

### repeat
- Description: Execute x repeatedly
- Arg count: 4
- Result count: 0
- Usage: 0 ⤶ < ⤶'f ⤶ repeat ⤶

### setloop
- Description: Set loop counter to x
- Arg count: 1
- Result count: 0
- Usage: 

