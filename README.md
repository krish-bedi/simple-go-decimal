# simple-go-decimal

## The Problem

Floats use base-2 which struggle at representing certain values such as 0.1

Just like how our everyday base-10 system struggles to represent 1/3 = 0.333333...

In case of monetary values, money is naturally base-10 which is why we store it in base-10 or decimal.

## The Solution

Store fractional values as integers to the 10th power. 
This lets us store numbers in base-10 and prevent rounding errors that would be present with using floats.

simple-go-decimal uses a precision of 4 i.e., up to 0.0001 can be represented perfectly
