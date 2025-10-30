package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("Converting yaml translation files to json...")

	data, err := os.ReadFile("internal/translations/en.yaml")
	if err != nil {
		panic(err)
	}

	translations := map[string]string{}
	err = yaml.Unmarshal(data, &translations)
	if err != nil {
		panic(err)
	}

	// Regex: {{ .Var }} → {{Var}}
	re := regexp.MustCompile(`\{\{\s*\.(\w+)\s*\}\}`)

	for k, v := range translations {
		translations[k] = re.ReplaceAllString(v, "{{$1}}")
	}

	out, err := json.MarshalIndent(translations, "", "	")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("ui/js/translations", 0755)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("ui/js/translations/en.json", out, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Converted yaml translation files to json.")
}
