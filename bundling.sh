#!/bin/bash
fyne bundle reverse.svg > bundled.go
fyne bundle -append start.svg >> bundled.go
fyne bundle -append agreement.svg >> bundled.go