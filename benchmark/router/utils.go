package router

import (
	"code.byted.org/lang/gg/gslice"
)

// getOptionalPaths returns all possible paths when the original path
// has optional arguments
func getOptionalPaths(path string) []string {
	paths := make([]string, 0)

	start := 0
walk:
	for {
		if start >= len(path) {
			return paths
		}

		c := path[start]
		start++

		if c != '{' {
			continue
		}

		newPath := ""
		hasRegex := false
		questionMarkIndex := -1

		brackets := 0

		for end, c := range []byte(path[start:]) {
			switch c {
			case '{':
				brackets++

			case '}':
				if brackets > 0 {
					brackets--
					continue
				} else if questionMarkIndex == -1 {
					continue walk
				}

				end++
				newPath += path[questionMarkIndex+1 : start+end]

				path = path[:questionMarkIndex] + path[questionMarkIndex+1:] // remove '?'
				paths = append(paths, newPath)
				start += end - 1

				continue walk

			case ':':
				hasRegex = true

			case '?':
				if hasRegex {
					continue
				}

				questionMarkIndex = start + end
				newPath += path[:questionMarkIndex]

				if len(path[:start-2]) == 0 {
					// include the root slash because the param is in the first segment
					paths = append(paths, "/")

				} else if !gslice.Contains(paths, path[:start-2]) {
					// include the path without the wildcard
					// -2 due to remove the '/' and '{'
					paths = append(paths, path[:start-2])
				}
			}
		}
	}
}
