# simple-go-decimal

## The Problem

Binary floating-point values cannot represent many base-10 fractions exactly. For
example, `0.1` has no finite binary representation, much like `1/3` has no
finite decimal representation.

This can make `float64` a poor choice for values such as money, where exact
base-10 arithmetic matters.

## The Solution

`simple-go-decimal` stores each value as a scaled `int64` instead of a binary
floating-point number.
With a precision of four decimal places, which is enough to cover most common monetary cases

```text
1.23 is stored as 12300
12300 / 10000 = 1.23
```
