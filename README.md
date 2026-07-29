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
The precision is set to a **default of four**. The precision can be changed via `const precision = 4` at the top of `decimal.go`

```text
1.23 is stored as 12300
12300 / 10000 = 1.23
```
## Benefits & Caveats

### Performance

| Operation            | Implementation                                              |
| -------------------- | ---------------------------------------------------------- |
| Addition, Subtraction | Native `int64` arithmetic with overflow checks             |
| Multiplication, Division | `big.Int` for the transient intermediate, result stored back as `int64` |

Addition and subtraction operate directly on the scaled `int64`, so they stay
allocation-free and take one CPU instruction. 

Multiplication and division briefly use a `big.Int` to hold an
intermediate value that can exceed `int64` before being scaled back down. This
avoids overflow

### Caveats

- **Range:** values must fit in an `int64` after scaling. With the default
  precision of 4, the range is roughly
  `-922,337,203,685,477.5808` to `922,337,203,685,477.5807`. Operations that
  exceed this return `ErrOverflow`.
- **Rounding:** multiplication and division truncate, so digits
  beyond the configured precision are dropped (e.g. `2 / 3` becomes `0.6666`,
  not `0.6667`).
