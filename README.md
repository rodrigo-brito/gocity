<img width="500" src="https://raw.githubusercontent.com/rodrigo-brito/gocity/master/logo.png" alt="GoCity" />

[![Build Status](https://travis-ci.org/rodrigo-brito/gocity.svg?branch=master)](https://travis-ci.org/rodrigo-brito/gocity)
[![Go Report Card](https://goreportcard.com/badge/github.com/rodrigo-brito/gocity)](https://goreportcard.com/report/github.com/rodrigo-brito/gocity)
[![GoDoc](https://godoc.org/github.com/rodrigo-brito/gocity?status.svg)](https://godoc.org/github.com/rodrigo-brito/gocity)
<a href="https://opensource.org/licenses/MIT">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square" alt="License MIT">
</a>
<hr />
 
Available here: https://go-city.github.io
 
GoCity is an implementation of the Code City metaphor for visualizing source code. GoCity represents a Go program as a city, as follows: 

 - Folders are districts
 - Files are buildings
 - Structs are represented as buildings on the top of their files.

## Structures Characteristics

 - The Number of Lines of Source Code (LOC) represents the build color (high values makes the building dark)
 - The Number of Variables (NOV) correlates to the building's base size.
 - The Number of methods (NOM) correlates to the building height.
 
## How To Run API in Localhost

To start the API in localhost
- Install Go 1.9+
- Clone repository
- Run: `go run main.go`
 
## UI / Front-end

The UI is built with React and uses [babylon.js](https://www.babylonjs.com/) to plot 3D structures. The front-end source code is available in the [front-end](https://github.com/rodrigo-brito/gocity/tree/front-end) branch. 
 
### Related Works
- [Code City](https://wettel.github.io/codecity.html) by [Richard Wettel](https://twitter.com/richardwettel)
- [JS City](https://github.com/ASERG-UFMG/JSCity/wiki/JSCITY) by [Marcus Viana](https://github.com/malilovick).

This tool makes part of a scientific research developed at Federal University of Minas Gerais (UFMG)<br/>
Student: [Rodrigo Brito](https://github.com/rodrigo-brito)<br/>
Advisor: [Marco Tulio Valente](https://homepages.dcc.ufmg.br/~mtov/)
