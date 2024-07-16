package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("$$$", os.Getenv("GOFLAGS"))

	fmt.Println(extractLdflagsFromGoFlags())

	fmt.Println("%%%", os.Getenv("GOFLAGS"))
}

func extractLdflagsFromGoFlags() string {
	goFlags := os.Getenv("GOFLAGS")

	ldFlags, goFlags := extractLdFlags(goFlags)
	fmt.Printf("!!! [%s]\n", ldFlags)
	fmt.Printf("@@@ [%s]\n", goFlags)
	if goFlags != "" {
		os.Setenv("GOFLAGS", goFlags)
	} else {
		os.Unsetenv("GOFLAGS")
	}

	return ldFlags
}

func extractLdFlags(goFlags string) (string, string) {
	if goFlags == "" {
		return "", ""
	}

	flags := strings.Fields(goFlags)
	ldflags := ""
	newGoFlags := ""

	for _, flag := range flags {
		if strings.HasPrefix(flag, "-ldflags=") {
			ldflags += strings.TrimPrefix(flag, "-ldflags=") + " "
		} else {
			newGoFlags += flag + " "
		}
	}

	ldflags = strings.TrimSpace(ldflags)
	newGoFlags = strings.TrimSpace(newGoFlags)

	return ldflags, newGoFlags
}
