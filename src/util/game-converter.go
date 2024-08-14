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
	Name            string `xml:"name,attr"`
	Polygon         []Coordinate
	Placements      []Coordinate
	IsCapitol       bool
	CapitolLocation *Coordinate
	CenterLocation  *Coordinate
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
	placementsFilePath := filepath.Join(mapFolderPath, "place.txt")
	capitolsFilePath := filepath.Join(mapFolderPath, "capitols.txt")
	//centersFilePath := filepath.Join(mapFolderPath, "centers.txt")

	for i, territory := range game.Territories {
		polygonCoordinates, err := processPolygonsFile(polygonsFilePath, territory.Name)
		placementCoordinates, err := processPlaceFile(placementsFilePath, territory.Name)
		CapitolCoordinate, err := processCapitolsFile(capitolsFilePath, territory.Name)
		//CenterCoordinate, err := processCapitolsFile(centersFilePath, territory.Name)

		if err != nil {
			fmt.Printf("Error processing polygons file: %v\n", err)
			return
		}
		game.Territories[i].Polygon = polygonCoordinates
		game.Territories[i].Placements = placementCoordinates
		//game.Territories[i].CenterLocation = CenterCoordinate

		/*if CenterCoordinate != nil {
			fmt.Printf("Center Coordinate nil: %s\n", territory.Name)
			return
		}*/

		if CapitolCoordinate != nil {
			game.Territories[i].IsCapitol = true
			game.Territories[i].CapitolLocation = CapitolCoordinate
		}
	}

	// Print the parsed territories
	for _, territory := range game.Territories {

		// Prepare the output string
		output := fmt.Sprintf(
			"Name: %s, # of Polygons: %d, # of Placements: %d, Is Capitol: %v",
			territory.Name,
			len(territory.Polygon),
			len(territory.Placements),
			territory.IsCapitol,
		)

		// If it's a capital, add the capital location
		if territory.IsCapitol {
			output += fmt.Sprintf(", Capitol Location: (%d, %d)", territory.CapitolLocation.X, territory.CapitolLocation.Y)
		}

		// Print the output
		fmt.Println(output)
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

func processPolygonsFile(fileName string, territoryName string) ([]Coordinate, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening polygons file: %v", err)
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

		currentTerritoryName := strings.TrimSpace(parts[0])
		coordinatesPart := strings.TrimSpace(parts[1])

		if territoryName != currentTerritoryName {
			continue
		}

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
				return nil, fmt.Errorf("error parsing X coordinate: %v", err)
			}

			y, err := strconv.Atoi(xy[1])
			if err != nil {
				return nil, fmt.Errorf("error parsing Y coordinate: %v", err)
			}

			coordinates = append(coordinates, Coordinate{X: x, Y: y})
		}

		return coordinates, nil
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading polygons file: %v", err)
	}

	return []Coordinate{}, nil
}

func processPlaceFile(fileName string, territoryName string) ([]Coordinate, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening polygons file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into the territory name and the list of coordinates.
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue // Skip lines that do not have enough parts
		}

		currentTerritoryName := parts[0]

		if territoryName != currentTerritoryName {
			continue
		}

		var coordinates []Coordinate

		// Parse each coordinate pair
		for _, coordStr := range parts[1:] {
			coordStr = strings.Trim(coordStr, "()")
			xy := strings.Split(coordStr, ",")
			if len(xy) != 2 {
				continue // Skip invalid coordinate pairs
			}

			x, err := strconv.Atoi(xy[0])
			if err != nil {
				return nil, fmt.Errorf("error parsing X coordinate: %v", err)
			}

			y, err := strconv.Atoi(xy[1])
			if err != nil {
				return nil, fmt.Errorf("error parsing Y coordinate: %v", err)
			}

			coordinates = append(coordinates, Coordinate{X: x, Y: y})
		}

		// Store the parsed coordinates in the map with the territory name as the key
		return coordinates, nil
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading polygons file: %v", err)
	}

	return []Coordinate{}, nil
}

func processCapitolsFile(filePath string, territoryName string) (*Coordinate, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening capitols file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into the territory name and the coordinate part.
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue // Skip lines that do not have enough parts
		}

		capitolName := parts[0]
		coordStr := strings.Trim(parts[1], "()")

		// Normalize both names by removing spaces for comparison
		normalizedCapitolName := strings.ReplaceAll(capitolName, " ", "")
		normalizedInputName := strings.ReplaceAll(territoryName, " ", "")

		// If the territory name does not match, skip to the next line
		if normalizedCapitolName != normalizedInputName {
			continue
		}

		// Parse the coordinate
		xy := strings.Split(coordStr, ",")
		if len(xy) != 2 {
			return nil, fmt.Errorf("invalid coordinate format: %s", coordStr)
		}

		x, err := strconv.Atoi(xy[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing X coordinate: %v", err)
		}

		y, err := strconv.Atoi(xy[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing Y coordinate: %v", err)
		}

		// Return true for is_capitol and the coordinate
		return &Coordinate{X: x, Y: y}, nil
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading capitols file: %v", err)
	}

	// If no match is found, return false and an empty coordinate
	return nil, nil
}
