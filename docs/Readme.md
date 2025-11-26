# optional

[![Go Reference](https://pkg.go.dev/badge/sweetkennedy.net/optional.svg)](https://pkg.go.dev/sweetkennedy.net/optional)

The _optional_ package provides an _option type_, which can either be empty or
hold a value. In that respect, it's very similar to an ordinary pointer type,
except it has methods that make its possible emptiness more explicit.

The interface of the `optional.Value` type is modeled on Java's
[`java.util.Optional`][java-opt] type and C++'s [`std::optional`][std-opt].

# Usage

```bash
go get sweetkennedy.net/optional
```

```go
import "sweetkennedy.net/optional"
```

```go
full := optional.New(42)
empty := optional.Value[int]{}

fmt.Printf("full: %v\n", full)  // Output: 42
fmt.Printf("empty: %v\n", empty) // Output: None
```

[java-opt]: https://docs.oracle.com/javase/8/docs/api/java/util/Optional.html
[std-opt]: https://en.cppreference.com/w/cpp/utility/optional
