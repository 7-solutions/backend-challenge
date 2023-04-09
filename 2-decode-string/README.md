# Decode string

The decode function receive string as input and return slice of integers according to these rules:
- L means the following number is decrement.
- R means the following number is increment.
- = means repeating the last element

## Assumption
1. The input string must be 'L', 'R' or '='
2. The output cannot be negative
3. The '=' letter cannot be at the beginning?

## Pseudocode
1. Read elements from input string
2. Translate string to number and append to stack.
3. If the element is 'L', the stack is polled from last to first index to the result e.g., stack{0,1,2} => result{2,1,0}.
4. If the element is 'R', the stack is polled from first to last element to the result e.g., stack{0,1,2} => result{0,1,2}.
5. If the lement is '=', the number append to the result corresponding to the last element in the result.