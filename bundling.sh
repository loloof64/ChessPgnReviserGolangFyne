#!/bin/bash
fyne bundle reverse.svg > bundled.go
fyne bundle -append start.svg >> bundled.go
fyne bundle -append stop.svg >> bundled.go
fyne bundle -append first.svg >> bundled.go
fyne bundle -append last.svg >> bundled.go
fyne bundle -append previous.svg >> bundled.go
fyne bundle -append next.svg >> bundled.go
fyne bundle -append chess.svg >> bundled.go
fyne bundle -append FreeSerif.ttf >> bundled.go