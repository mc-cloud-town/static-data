package utils

import "strings"

func GetPortFromURL(urlString string) string {
	// Check if the URL contains a colon, indicating a port number
	if strings.Contains(urlString, ":") {
		// Split the URL by colon
		parts := strings.Split(urlString, ":")
		// Check if there's a port number after the colon
		if len(parts) > 1 {
			// Return the part after the colon, which should be the port number
			return parts[len(parts)-1]
		}
	}
	// Return an empty string if no port number is found
	return ""
}
