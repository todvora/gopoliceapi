package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func toInteger(component string) int {
	val, err := strconv.Atoi(component)
	if err != nil {
		log.Fatalf("Failed to parse component %s", val)
	}
	return val
}

func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return version[1:]
	}
	return version
}

func increment(component string, version string) string {
	parts := strings.Split(normalizeVersion(version), ".")

	if len(parts) != 3 {
		log.Fatalf("Provided version doesn't correspond to a valid semver value: %s", version)

	}
	major := toInteger(parts[0])
	minor := toInteger(parts[1])
	patch := toInteger(parts[2])


	switch component {
	case "major":
		{
			major += 1
			minor = 0
			patch = 0

		}
	case "minor":
		{
			minor += 1
			patch = 0

		}
	case "patch":
		{
			patch += 1
		}

	default:
		log.Fatalf("Unknown component: %s", component)
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)

}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Required exactly 2 arguments in format major|minor|patch X.Y.Z, %d given: %s", len(os.Args)-1, os.Args[1:])

	}
	component := os.Args[1]
	version := os.Args[2]
	fmt.Println(increment(component, version))
}
