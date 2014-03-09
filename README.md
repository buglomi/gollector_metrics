# gollector\_metrics: A Linux system metrics collection library

This package is a collection of ... collectors which gather metrics from
various Linux subsystems. To build it, you must have a C compiler available and
glibc. Running the monitors themselves requires a linux-compliant /proc
filesystem for at least a Linux 3.0 kernel.

This package will be imported by gollector for gathering its system metrics. At
this point it should be considered alpha quality as it has no known consumers.
It also needs significant refactors.

There are many monitors available and you can find the API documentation at
[godoc.org](http://godoc.org/github.com/gollector/gollector_metrics).

## Usage

```
go get -u -d https://github.com/gollector/gollector_metrics
```

```go
import (
  gm "github.com/gollector/gollector_metrics"
)
```

Since this package will change frequently, it is strongly advised you use a
tool like [godep](https://github.com/kr/godep) to manage versions.

## License

MIT, (c) 2014 Erik Hollensbe. Earlier editions may be found in the
[gollector](https://github.com/gollector/gollector) source which is (c) 2013
Erik Hollensbe, also MIT licensed.

## Author

Erik Hollensbe <erik+github@hollensbe.org>
