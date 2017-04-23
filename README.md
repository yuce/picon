# PiCon

## Build

Although there is inital support for Glide, this project depends on a development branch of Pilosa Go client which is not accessible to Glide.

PiCon can be installed to a custom $GOPATH as follows:
```
$ export GOPATH=$HOME/picon
$ mkdir -p $GOPATH/src/github.com/pilosa && cd $GOPATH/src/github.com/pilosa
$ git clone git@github.com:yuce/go-client-pilosa.git
$ cd go-client-pilosa && git checkout v2
$ mkdir -p $GOPATH/src/bitbucket.org/yuce && cd $GOPATH/src/bitbucket.org/yuce
$ git clone git@bitbucket.org:yuce/picon.git && cd picon
$ go get bitbucket.org/yuce/picon/cmd/picon
$ go build bitbucket.org/yuce/picon/cmd/picon
$ ls picon
```

## Usage

You can run the console with `./picon`. To get a list of commands, hit `:` and then `Tab`. To exit, you can type `:exit` or hit `Ctrl+C`.

- Commands start with `:`.
- Notes start with `#`.
- Queries can be run directly.
- In order to enter multiline commands/queries, finish a line with backslash (`\`).
- Up/down arrow keys can be used to access the history.
- `:use` command supports index name completion.
- If a command is made up of more than one word, they can be autocompleted.

Sample workflow:

```
> :connect :10101
> :ensure index myindex
> :ensure frame myframe
> SetBit(id=1, frame='myframe', col_id=100)
> Bitmap(id=1, frame='myframe')
... Some output
```

