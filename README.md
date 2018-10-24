<img width="500" src="https://raw.githubusercontent.com/rodrigo-brito/gocity/master/logo.png" alt="GoCity" />

[![Build Status](https://travis-ci.org/rodrigo-brito/gocity.svg?branch=master)](https://travis-ci.org/rodrigo-brito/gocity)
[![Go Report Card](https://goreportcard.com/badge/github.com/rodrigo-brito/gocity)](https://goreportcard.com/report/github.com/rodrigo-brito/gocity)
[![GoDoc](https://godoc.org/github.com/rodrigo-brito/gocity?status.svg)](https://godoc.org/github.com/rodrigo-brito/gocity)
<a href="https://opensource.org/licenses/MIT">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square" alt="License MIT">
</a>
<hr />
 
 Available here: http://gocity.rodrigobrito.net/
 
GoCity is a web tool to visualize Go source code. This project is inspired on Code City project created by [Richard Wettel](https://twitter.com/richardwettel) and [JS City](https://github.com/ASERG-UFMG/JSCity/wiki/JSCITY) created by [Marcus Viana](https://github.com/malilovick).

This tool makes part of a scientific research developed at Federal University of Minas Gerais (UFMG) and was advised by [Marco Túlio](https://homepages.dcc.ufmg.br/~mtov/)

## Structures Information

 - Folders are districts
 - Files are buildings
 - Structs are represented as buildings on the top of their files.

## Structures Characteristics

 - The Number of Lines of Source Code (LOC) represents the build color (high values makes the building dark)
 - The Number of Variables (NOV) in a function correlates to the building's base size.
 - The Number of methods (NOM) correlates to the building height.
