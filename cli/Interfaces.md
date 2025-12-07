# Interfaces

#go #dev 


sdfasdf

lala
sdfsdf
[link test](https://google.com)!
another test sdfsdf
Interfaces are **named collections of method signatures**.
Some *italic text too just to see*.

 How about a footnote[^1]

[^1]: My reference is cool!

> A quote! Wauw. So good.

Now for a list:

- one
- two
- three
    - indent
    - another indent

And not for an ordered list!

1) Something
1) Some other thing
1) Some even other thing

Now for a TODO:

- [ ] something
- [ ] another thing
- [x] a finished thing!


> [!faq] My title
> My first line
> My second line
> My third line
> My fourth line
> My fifth line
> My sixth line
> My seventh line
> My eighth line
> My ninth `line` 
> ```py
> def test():
>     return None
> ```
> something *else*


We define it like so:

```go
type geometry interface {
    area() float64
    perim() float64
}
```

To implement this interface on these 2 structs:

```go
type rect struct {
    width, height float64
}
type circle struct {
    radius float64
}
```

We would just implement all the method signatures of the interface on them:

```go
func (r rect) area() float64 {
    return r.width * r.height
}
func (r rect) perim() float64 {
    return 2*r.width + 2*r.height
}

func (c circle) area() float64 {
    return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
    return 2 * math.Pi * c.radius
}
```

That's it!

Now we can use the interface to make this function work on all types that implement the `geometry` interface:

```go
func measure(g geometry) {
    fmt.Println(g)
    fmt.Println(g.area())
    fmt.Println(g.perim())
}
```

## Empty Interfaces

```go
var a interface{}
```

If we would do the following:

```go
fmt.Println(a)
```

We would see `nil`.

This empty interface can accept any number of values of any type:

```go
a = "Hello World"
a = 5
```

They can be used for function to accept parameters of any type:

```go
func displayValue(i interface{}) {
  fmt.Println(i)
}
```

The `any` keyword is exactly the same as `interface{}`.

### Type Assertions

To help resolve the type of values stores by empty interfaces, we can use **type assertions**.

```go
var a interface{}
a = 10
myint := a.(int)  // <-
```

This `.(int)` will check if the value of `a` is an integer or not, and if it is, return it.
`.(int)` would panic if it isn't an integer. To avoid this panic, the type assertion actually returns a second boolean value that tells us whether or not the assertion succeeded:

```go
mystr, ok := a.(string)
```

If we read it, no panic will occur.

### Type Switches

In the same vein as type assertions, we can use `.(type)` in a switch statement called a **type switch**:

```go
switch v := i.(type) {
case T:
    // here v has type T
case S:
    // here v has type S
default:
    // no match; here v has the same type as i
}
```

## Stringers

One of the most ubiquitous interfaces is [`Stringer`](https://go.dev/pkg/fmt/#Stringer) defined by the [`fmt`](https://go.dev/pkg/fmt/) package:

```go
type Stringer interface {
    String() string
}
```

A stringer is a type that can describe itself as a string. As you can see, all types that implement the `String()` method are Stringers. The `fmt` package (and many others) look for this interface to print values in various places.
