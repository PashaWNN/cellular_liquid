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


## How it works

In few words, all algorithm is built on top of only three rules: falling, spreading and decompressing.

Each cell can have limited mass of liquid and may be under limited pressure.

* First of all, liquid is *falling down* until it reaches solid surface or already compressed cell of liquid.
* Then second rule is applied: right after falling down, cell starts to *spread* around. So it will *spread over the surface*.
* Finally, compressed cells starts to flow back upwards.

Most sensitive part of this algorithm is getStableState function which calculates how much cells will compress and decompress. This function allows fluid to have realistic physics and have the effect of communicating vessels.
