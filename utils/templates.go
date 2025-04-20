// File: utils/templates.go
package utils

import (
	"fmt"
	"regexp"
)

// GetJobTemplates returns available job templates as JSON
func GetJobTemplates() (string, error) {
	// In a real implementation, this would read templates from a directory or database
	// templates := map[string]string{
	// 	"service": "Basic service job template",
	// 	"batch":   "Basic batch job template",
	// 	"system":  "Basic system job template",
	// }

	// Convert to JSON string
	return `{
		"templates": [
			{
				"name": "service",
				"description": "Basic service job template"
			},
			{
				"name": "batch",
				"description": "Basic batch job template"
			},
			{
				"name": "system",
				"description": "Basic system job template"
			}
		]
	}`, nil
}

// GetJobTemplate returns a specific job template
func GetJobTemplate(name string) (string, error) {
	// In a real implementation, this would read from a file or database

	switch name {
	case "service":
		return `job "example-service" {
  datacenters = ["dc1"]
  type = "service"

  group "web" {
    count = 2

    network {
      port "http" {
        to = 8080
      }
    }

    task "server" {
      driver = "docker"

      config {
        image = "nginx:latest"
        ports = ["http"]
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}`, nil
	case "batch":
		return `job "example-batch" {
  datacenters = ["dc1"]
  type = "batch"

  group "batch-group" {
    count = 1

    task "batch-task" {
      driver = "docker"

      config {
        image = "alpine:latest"
        command = "/bin/sh"
        args = ["-c", "echo 'Processing data' && sleep 5"]
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}`, nil
	case "system":
		return `job "example-system" {
  datacenters = ["dc1"]
  type = "system"

  group "system-group" {
    task "system-task" {
      driver = "docker"

      config {
        image = "consul:latest"
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}`, nil
	default:
		return "", fmt.Errorf("template not found: %s", name)
	}
}

// ExtractTemplateNameFromURI extracts the template name from the URI
func ExtractTemplateNameFromURI(uri string) string {
	re := regexp.MustCompile(`^nomad://templates/(.+)$`)
	matches := re.FindStringSubmatch(uri)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}
