# godepviz

Visualize Go imports in Graphviz DOT format. 

Support visualize package from [https://pkg.go.dev](https://pkg.go.dev/)


## Installation

```bash
$ go get github.com/chuongtrh/godepviz
```

## Usage

```bash
$ godepviz github.com/gofiber/fiber/v2
```
The output is a graph in [Graphviz](http://graphviz.org/) dot format. 

```base
digraph G {
         rankdir="LR";
    pad=.25;
    ratio=fill;
    dpi=360;
    node [shape=box];
        "github.com/gofiber/fiber/v2" -> "github.com/valyala/fasthttp";
        "github.com/valyala/fasthttp" -> "github.com/andybalholm/brotli";
        "github.com/valyala/fasthttp" -> "github.com/klauspost/compress/flate";
        "github.com/valyala/fasthttp" -> "github.com/klauspost/compress/gzip";
        "github.com/valyala/fasthttp" -> "github.com/klauspost/compress/zlib";
        "github.com/valyala/fasthttp" -> "github.com/valyala/bytebufferpool";
        "github.com/gofiber/fiber/v2" -> "github.com/valyala/fasthttp/reuseport";
        "github.com/valyala/fasthttp/reuseport" -> "github.com/valyala/tcplisten";
        "github.com/gofiber/fiber/v2"  [style=filled];
}
```

You can view a output graph with tool [https://dreampuf.github.io/GraphvizOnline](https://dreampuf.github.io/GraphvizOnline)
or [http://viz-js.com](http://viz-js.com)

If you have the **graphviz** tools installed you can render it by piping the output to dot:

```bash
$ godepviz github.com/gofiber/fiber/v2 | dot -Tpng -o godep.png
```

Export SVG file with a large graph have many nodes of dependencies
```bash
$ godepviz github.com/gofiber/fiber/v2 | dot -Gdpi=0 -T svg -o godep.svg
```
## Example

- Package `net/http`

![Example net/http](./screenshots/godep.png)


- Package `github.com/labstack/echo/v4`

![Example github.com/labstack/echo/v4](./screenshots/godep1.png)

- Package `github.com/gofiber/fiber/v2`

![Example github.com/gofiber/fiber/v2](./screenshots/godep.svg)
