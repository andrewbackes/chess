[![Go Report Card](https://goreportcard.com/badge/github.com/andrewbackes/chess)](https://goreportcard.com/report/github.com/andrewbackes/chess) [![GoDoc](https://godoc.org/github.com/andrewbackes/chess?status.svg)](https://godoc.org/github.com/andrewbackes/chess) [![Build Status](https://travis-ci.org/andrewbackes/chess.svg?branch=master)](https://travis-ci.org/andrewbackes/chess) [![Coverage Status](https://coveralls.io/repos/github/andrewbackes/chess/badge.svg?branch=master)](https://coveralls.io/github/andrewbackes/chess?branch=master)

# chess
Multipurpose chess package for Go/Golang.

###What does it do?
This package provides tools for working with chess games. You can:
- Play games
- Detect checks, checkmates, and draws (stalemate, threefold, 50 move rule, insufficient material)
- Open PGN files or strings and filter them
- Open EPD files or strings
- Load FENs
- Generate legal moves from any position.
- and more

For details you can visit the [godoc](https://godoc.org/github.com/andrewbackes/chess)

##How to get it
If you have your GOPATH set in the recommended way ([](https://golang.org/doc/code.html#GOPATH)):

```go get github.com/andrewbackes/chess```

otherwise you can clone the repo.