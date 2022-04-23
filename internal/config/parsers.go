package config

import (
	"fmt"
	"regexp"
	"strings"
)

func parseServices(input string) ([]Service, error) {
	if input == "" {
		return []Service{}, nil
	}

	serializedServices := strings.Split(input, ",")
	services := make([]Service, 0, len(serializedServices))

	alphaNumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`)

	for _, service := range serializedServices {
		majorParts := strings.Split(service, "@")
		if len(majorParts) != 2 {
			return nil, fmt.Errorf(
				"expected only one occurrence of '@', got %d in %s",
				len(majorParts)-1,
				service,
			)
		}

		minorParts := strings.Split(majorParts[0], ":")
		if len(minorParts) > 2 {
			return nil, fmt.Errorf(
				"expected at most one occurrence of ':', got %d in %s",
				len(minorParts)-1,
				service,
			)
		}

		if !alphaNumeric.MatchString(minorParts[0]) {
			return nil, fmt.Errorf("unexpected non-alphanumeric characters in %s", service)
		}

		secret := ""
		if len(minorParts) > 1 {
			if !alphaNumeric.MatchString(minorParts[1]) {
				return nil, fmt.Errorf("unexpected non-alphanumeric characters in %s", service)
			}
			secret = minorParts[1]
		}

		services = append(services, Service{
			Label:    minorParts[0],
			Secret:   secret,
			Endpoint: majorParts[1],
		})
	}

	return services, nil
}
