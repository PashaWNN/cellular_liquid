# Liquid simulation cellular automata

Implementation of cellular automata that simulated liquid on top of HTML5 canvas using WebAssembly.

[Live demo](https://pashawnn.github.io/cellular_liquid/)


## Compilation

```
GOOS=js GOARCH=wasm go build -o main.wasm 
```

Note that you can't just open index.html from local filesystem. You need web-server which sets correct mime type (`mime .wasm application/wasm`) to run WebAssembly. Simpliest solution is to download Caddy server and just run:
```
caddy
```
from project root.
