# Decode string

The decode function receive string as input and return slice of integers according to these rules:
- L means the following number is decrement.
- R means the following number is increment.
- = means repeating the last element

## Assumption
1. The input string must be 'L', 'R' or '='
2. The output cannot be negative
3. The '=' letter cannot be at the beginning?