# About
28z is a simple VM loosely inspired by the RPL environment on the HP-28 calculator.

# Building
go build .

# Running
./28z

## Data types

### Floating point
All values are assumed to be floating point by default

### String
String values are identified by a single preceeding quotation mark (')

## Supported instructions

### >>
- Description: Reduce function
- Arg count: 0
- Result count: 0
- Usage: 

### enter
- Description: Enter function
- Arg count: 1
- Result count: 0
- Usage: 

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

### eval
- Description: Evaluate x
- Arg count: 1
- Result count: 0
- Usage: 

### -
- Description: Subtract x from y
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ - ⤶ ⤒4

### store
- Description: Store y into x
- Arg count: 2
- Result count: 1
- Usage: 2 ⤶ 'a ⤶ store ⤶ ⤒2; y⥗a

### recall
- Description: Recall x
- Arg count: 1
- Result count: 1
- Usage: 'a ⤶ recall ⤶ ⤒a

### print
- Description: Print x
- Arg count: 1
- Result count: 0
- Usage: 'Hello world ⤶ print ⤶ Hello world⥱Console

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

### clear
- Description: Clear stack
- Arg count: 0
- Result count: 0
- Usage: clear ⤶

### put
- Description: Put y into x
- Arg count: 2
- Result count: 0
- Usage: 2 ⤶ 'a ⤶ store ⤶ y⥗a

### swap
- Description: Swap x and y
- Arg count: 2
- Result count: 2
- Usage: swap ⤶ ⤒x,y

### status
- Description: Display status
- Arg count: 0
- Result count: 0
- Usage: 

### /
- Description: Divide y by x
- Arg count: 2
- Result count: 1
- Usage: 6 ⤶ 2 ⤶ / ⤶ ⤒3

### <<
- Description: Define function
- Arg count: 1
- Result count: 0
- Usage: 

### end
- Description: Return from function
- Arg count: 1
- Result count: 0
- Usage: 

### exchange
- Description: Exchange y and the value in var x
- Arg count: 2
- Result count: 1
- Usage: 3 ⤶ 'a ⤶ exchange ⤶ ⤒a 3⥗a

### halt
- Description: Halt execution
- Arg count: 0
- Result count: 0
- Usage: halt ⤶

### collect
- Description: Collect stack into x
- Arg count: 1
- Result count: 1
- Usage: 

### graph
- Description: Render graph
- Arg count: 0
- Result count: 0
- Usage: 

