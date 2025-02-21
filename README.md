# GOCK3-parse - PDXScript Tools for Crusader Kings 3

[![Go](https://github.com/unLomTrois/gock3/actions/workflows/go.yml/badge.svg)](https://github.com/unLomTrois/gock3/actions/workflows/go.yml)

**GOCK3** is a collection of tools written in Go for tokenizing, parsing, and validating PDXScript files used in [Crusader Kings 3](https://www.crusaderkings.com/). This project aims to assist mod developers by providing utilities that can analyze PDXScript code, catch errors, and improve code quality.

> This project was inspired by [ck3-tiger](https://github.com/amtep/ck3-tiger). Some concepts and code structures have been adapted with gratitude.

## Features

The project consists of three main components:

- **Lexer**: Tokenizes PDXScript code and catches lexical errors, such as unknown tokens (e.g., invalid operators like `!=`).
- **Parser**: Constructs an Abstract Syntax Tree (AST) from the token stream and catches syntax errors (e.g., unclosed curly braces).
