package pgnLoader

import (
	"bufio"
	"os"
	"strings"
)

// Loader holds all games of a loaded PGN file.
type Loader struct {
	Games []string
}

// LoadPgnFile tries to load all games from a PGN file.
func LoadPgnFile(path string) (*Loader, error) {
	fileContent, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fileContent.Close()
	scanner := bufio.NewScanner(fileContent)

	games := []string{}
	currentGame := ""
	processingGameHeaders := false

	for scanner.Scan() {
		line := scanner.Text()
		weMayNeedToRegisterGameAndStartNewOne := strings.HasPrefix(line, "[")
		if weMayNeedToRegisterGameAndStartNewOne {
			weDoWantToRegisterGameAndStartNewOne :=
				(len(games) == 0 || currentGame != "") &&
					!processingGameHeaders
			if weDoWantToRegisterGameAndStartNewOne {
				if currentGame != "" {
					games = append(games, currentGame)
				}
				currentGame = ""
				processingGameHeaders = true
			}
		} else {
			processingGameHeaders = false
		}
		currentGame += line
		currentGame += "\n"
	}

	if currentGame != "" {
		games = append(games, currentGame)
	}

	return &Loader{Games: games}, nil
}
