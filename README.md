# PiCon

<a href="https://github.com/pilosa"><img src="https://img.shields.io/badge/pilosa-v0.4.0-blue.svg"></a>
<img src="https://c1.staticflickr.com/9/8754/16788993048_af85d47b1b_z.jpg" style="float: right" align="right" height="180">

Simple console for [Pilosa](https://www.pilosa.com/) high performance distributed bitmap index.

This app uses JSON API of Pilosa directly for queries and the official [Go client](https://github.com/pilosa/go-pilosa) for everything else (e.g., creating indexes, frames, ...).

## Change Log

* **0.2.0** (2017-06-17)
    * Initial public version

## Dependencies

* [Pilosa Go Client](https://github.com/pilosa/go-pilosa)
* [Readline](https://github.com/chzyer/readline)
* [Go Pretty JSON](github.com/hokaccha/go-prettyjson)

## Build

```
go get github.com/yuce/picon/cmd/picon && go build github.com/yuce/picon/cmd/picon
```

## Usage

You can run the console with `picon`.

- Commands start with `:`.
- To get a list of commands, hit `:` and then `Tab`.
- To exit, you can type `:exit` or hit `Ctrl+D`.
- Notes start with `#`.
- Queries can be run directly.
- In order to enter multiline commands/queries, finish a line with backslash (`\`).
- Up/down arrow keys can be used to access the history.
- Hit tab for PQL call completion.
- `:use` command supports index name completion.
- `:connect` command supports address completion using the history.
- If a command is made up of more than one word, they can be autocompleted.

Sample workflow:

```
> :connect :10101
> :ensure index myindex col=col_id
> :ensure frame myframe inverse=true
> SetBit(frame='myframe', rowID=1, col_id=100)
> Bitmap(frame='myframe', rowID=1)
... Some output
```

### Available commands

* `:connect`: Connect to the Pilosa server. Usage: `:connect pilosa-address`.
* `:create`: Create an index or a frame. Usage: `:create {index | frame} name [option1=value1, ...]`.
* `:delete`: Delete an index or a frame. Usage: `:delete {index | frame} name1, ...`.
* `:ensure`: Ensure that an index or a frame exists. Usage: `:ensure {index | frame} name [option1=value1, ...]`.
* `:http`: Send a raw HTTP request to the server. See: [API Documentation](https://www.pilosa.com/docs/api-reference/). Usage: `:http method path [data]`.
* `:schema`: Display the scheme (indexes and frames) on the server. Usage: `:schema`.
* `:use`: Open an index. Usage: `:use index-name`.

`:create index` and `:ensure index` commands support the following options:
* `column_label`, `columnLabel`, `col`, `c`
* `time_quantum`, `timeQuantum`, `time`, `t`

`:create frame` and `:ensure frame` commands support the following options:
* `row_label`, `rowLabel`, `row`, `r`
* `time_quantum`, `timeQuantum`, `time`, `t`
* `inverse_enabled`, `inverseEnabled`, `inverse`, `i`.

You can use `true`, `t` or `1` for true value; `false`, `f` or `0` false value.

### Queries

Any valid PQL query can be executed directly. See: [PQL Documentation](https://www.pilosa.com/docs/query-language/)

## License

```
Copyright 2017 Yuce Tekol

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived
from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH
DAMAGE.
```
