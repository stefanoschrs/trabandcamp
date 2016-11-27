<p align="center">
	<img src="http://res.cloudinary.com/dkxp3eifs/image/upload/c_scale,w_200/v1465057926/go-bc-logo_ofgay7.png"/>
</p>

# Trabandcamp [![Build Status](https://img.shields.io/travis/stefanoschrs/trabandcamp/master.svg?style=flat-square)](https://travis-ci.org/stefanoschrs/trabandcamp) [![License](https://img.shields.io/github/license/stefanoschrs/trabandcamp.svg?style=flat-square)]()

Download tracks from bandcamp **GO** style

> *Full documentation can be found at the Trabandcamp's [Wiki Page](https://github.com/stefanoschrs/trabandcamp/wiki)*

Installation
-
Download the latest binary from the [releases](https://github.com/stefanoschrs/trabandcamp/releases) page

Usage
-
- `./trabandcamp-<os>-<arch>[.exe] <Band Name>`    
*e.g for https://dopethrone.bandcamp.com/ on a 64bit Linux you should run* `./trabandcamp-linux-amd64 dopethrone`
- If you want to change the download directory you can add a `.trabandcamprc` file (there is already a sample one for you to copy)

Development
-
If you want to build the binary yourself  
`(export GOOS=<Operating System>; export GOARCH=<Architecture>; go build -o build/trabandcamp-$GOOS-$GOARCH trabandcamp.go)`  
*if you want to target ARM architecture you should also add the GOARM variable*
