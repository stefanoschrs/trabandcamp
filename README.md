<p align="center">
	<img src="http://res.cloudinary.com/dkxp3eifs/image/upload/c_scale,w_200/v1465057926/go-bc-logo_ofgay7.png"/>
</p>

# Trabandcamp [![Build Status](https://travis-ci.org/stefanoschrs/trabandcamp.svg?branch=master)](https://travis-ci.org/stefanoschrs/trabandcamp)
Download tracks from bandcamp **GO** style

Installation
-
Download the latest binary from the [releases](https://github.com/stefanoschrs/trabandcamp/releases) page

Usage
-
`./trabandcamp-<os>-<arch>[.exe] <Band Name>`    
*e.g for https://dopethrone.bandcamp.com/ on a 64bit Linux you should run* `./trabandcamp-linux-amd64 dopethrone`

Development
-
If you want to build the binary yourself  
`(export GOOS=<Operating System>; export GOARCH=<Architecture>; go build -o build/trabandcamp-$GOOS-$GOARCH trabandcamp.go)`  
*if you want to target ARM architecture you should also add the GOARM variable*
