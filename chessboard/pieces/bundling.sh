#!/bin/bash
fyne bundle Chess_plt45.svg > bundled.go
fyne bundle -append Chess_nlt45.svg >> bundled.go
fyne bundle -append Chess_blt45.svg >> bundled.go
fyne bundle -append Chess_rlt45.svg >> bundled.go
fyne bundle -append Chess_qlt45.svg >> bundled.go
fyne bundle -append Chess_klt45.svg >> bundled.go

fyne bundle -append Chess_pdt45.svg >> bundled.go
fyne bundle -append Chess_ndt45.svg >> bundled.go
fyne bundle -append Chess_bdt45.svg >> bundled.go
fyne bundle -append Chess_rdt45.svg >> bundled.go
fyne bundle -append Chess_qdt45.svg >> bundled.go
fyne bundle -append Chess_kdt45.svg >> bundled.go