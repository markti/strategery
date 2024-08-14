package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Coordinate struct {
	X int
	Y int
}

type Territory struct {
	Name        string `xml:"name,attr"`
	Coordinates []Coordinate
}

type Game struct {
	Territories []Territory `xml:"map>territory"`
}

func main() {
	fmt.Print("Hello, World!")

	args := os.Args

	if len(args) < 2 {
		fmt.Print("No command-line arguments provided")
		os.Exit(1)
	}

	// Get the folder path and filename from the command-line arguments
	rootFolderPath := os.Args[1]
	fileName := os.Args[2]

	fmt.Printf("Root Folder path: %s\n", rootFolderPath)
	fmt.Printf("File name: %s\n", fileName)

	mapFolderPath := filepath.Join(rootFolderPath, "map")
	gamesFolderPath := filepath.Join(mapFolderPath, "games")
	gameFilePath := filepath.Join(gamesFolderPath, fileName)

	fmt.Printf("Game File name: %s\n", gameFilePath)

	// Parse the XML file.
	game, err := parseXMLFile(gameFilePath)
	if err != nil {
		fmt.Printf("Error parsing XML file: %v\n", err)
		return
	}

	polygonsFilePath := filepath.Join(mapFolderPath, "polygons.txt")
	err = processPolygonsFile(polygonsFilePath, game)
	if err != nil {
		fmt.Printf("Error processing polygons file: %v\n", err)
		return
	}

	// Print the parsed territories
	for _, territory := range game.Territories {
		fmt.Println(territory.Name, " (", len(territory.Coordinates), ")")
		//fmt.Println("Coordinates:")
		/*for _, coord := range territory.Coordinates {
			fmt.Printf("(%d, %d) ", coord.X, coord.Y)
		}*/
		//fmt.Println("\n")
	}
}

// parseXMLFile reads and parses the XML file into a Game struct.
func parseXMLFile(fileName string) (*Game, error) {
	// Open the XML file.
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read the file's content.
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Unmarshal the XML data into the Game struct.
	var game Game
	err = xml.Unmarshal(data, &game)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %v", err)
	}

	return &game, nil
}

func processPolygonsFile(fileName string, game *Game) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening polygons file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into the territory name and the list of coordinates.
		parts := strings.SplitN(line, "<", 2)
		if len(parts) != 2 {
			continue
		}

		territoryName := strings.TrimSpace(parts[0])
		coordinatesPart := strings.TrimSpace(parts[1])

		//fmt.Println("Territory:", territoryName)
		//fmt.Println("Coordinates:", coordinatesPart)

		// Clean up the coordinatesPart by removing '>', and splitting by space.
		coordinatesPart = strings.TrimSuffix(coordinatesPart, ">")
		coordinatesStrs := strings.Fields(coordinatesPart)

		// Parse the coordinates and create a slice of Coordinate objects.
		var coordinates []Coordinate
		for _, coordStr := range coordinatesStrs {
			coordStr = strings.Trim(coordStr, "()")
			xy := strings.Split(coordStr, ",")
			if len(xy) != 2 {
				continue
			}

			x, err := strconv.Atoi(xy[0])
			if err != nil {
				return fmt.Errorf("error parsing X coordinate: %v", err)
			}

			y, err := strconv.Atoi(xy[1])
			if err != nil {
				return fmt.Errorf("error parsing Y coordinate: %v", err)
			}

			coordinates = append(coordinates, Coordinate{X: x, Y: y})
		}

		// Find the matching territory and assign the coordinates.
		for i := range game.Territories {
			if game.Territories[i].Name == territoryName {
				game.Territories[i].Coordinates = coordinates
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading polygons file: %v", err)
	}

	return nil
}
