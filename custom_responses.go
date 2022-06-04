package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

type (
	CustomResponseFolder struct {
		Path   string
		Routes map[string][]byte
	}
	CustomResponse struct {
		Path       string `json:"path"`
		SourceFile string `json:"source_file"`
	}
	CustomResponseData struct {
		Path    string
		Content []byte
	}
)

func GetCustomResponseFolder(folderName string) (crf *CustomResponseFolder, err error) {
	if stat, err := os.Stat(folderName); err != nil || !stat.IsDir() {
		return nil, fmt.Errorf("Custom response folder %s does not exist", folderName)
	}
	folderName, _ = filepath.Abs(folderName)
	routesFile := filepath.Join(folderName, "routes.json")
	if _, err := os.Stat(routesFile); err != nil {
		return nil, fmt.Errorf("Expected file %s not found", routesFile)
	}
	content, err := os.ReadFile(routesFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read routes.json from %s - %v", routesFile, err)
	}
	var routes []CustomResponse
	if err = json.Unmarshal(content, &routes); err != nil {
		return nil, fmt.Errorf("Failed to parse routes.json from %s - %v", routesFile, err)
	}
	crf = &CustomResponseFolder{
		Path:   folderName,
		Routes: make(map[string][]byte),
	}

	for _, route := range routes {
		if len(route.Path) == 0 {
			log.Printf("Invalid custom response: empty path - %v", route)
			continue
		}
		route.SourceFile = filepath.Join(folderName, filepath.Base(route.SourceFile))
		if stat, err := os.Stat(route.SourceFile); err != nil || stat.IsDir() {
			log.Printf("Invalid custom response: source file %s does not exist", route.SourceFile)
			continue
		}
		content, err = os.ReadFile(route.SourceFile)
		if err != nil {
			log.Printf("Failed to read file %s - %v", route.SourceFile, err)
			continue
		}
		crf.Routes[route.Path] = content
	}
	if len(crf.Routes) == 0 {
		return nil, fmt.Errorf("No valid custom responses found in %s", folderName)
	}
	return crf, nil
}

func GetCustomResponses(configFile string) (responses []CustomResponseData, err error) {
	log.Printf("Reading custom responses from %s", configFile)
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %s - %v", configFile, err)
	}
	var config []CustomResponse
	err = json.Unmarshal(content, &config)
	if err != nil {
		return
	}
	for _, response := range config {
		if response.Path == "" {
			log.Printf("Invalid response path - %s", response.Path)
			continue
		}
		content, err = os.ReadFile(response.SourceFile)
		if err != nil {
			log.Printf("Failed to read file %s - %v", response.SourceFile, err)
			continue
		}
		responses = append(responses, CustomResponseData{
			Path:    response.Path,
			Content: content,
		})
		log.Printf("Adding custom response for path %s -> %s", response.Path, response.SourceFile)
	}
	if len(responses) == 0 {
		err = fmt.Errorf("No valid responses found in %s", configFile)
	}
	return
}

func GetArgResponsesFolder() (customResponsesFolder string) {
	customResponsesFolder = os.Getenv("CUSTOM_RESPONSES_FOLDER")
	source := ""
	if len(customResponsesFolder) > 0 {
		source = "env CUSTOM_RESPONSES_FOLDER=" + customResponsesFolder
	} else {
		getNext := false
		for _, arg := range os.Args {
			if arg == "-c" || arg == "--config" {
				getNext = true
			} else if getNext {
				customResponsesFolder = arg
				source = "command line '--config " + customResponsesFolder + "'"
				break
			}
		}
	}

	if len(customResponsesFolder) > 0 {
		log.Printf("Using custom responses config folder from %s", source)
		cfgFile, err := filepath.Abs(customResponsesFolder)
		if err != nil {
			log.Printf("Failed to get absolute path for %s - %v", customResponsesFolder, err)
		} else {
			if _, err := os.Stat(cfgFile); err != nil {
				log.Printf("Failed to find custom responses config folder %s - %v", cfgFile, err)
				customResponsesFolder = ""
			} else {
				customResponsesFolder = cfgFile
				log.Printf("Custom responses config folder: %s", customResponsesFolder)
			}
		}
	}

	return
}

func SetupHttpCustomHandlers() (routes []string) {
	responsesFolder := GetArgResponsesFolder()
	if len(responsesFolder) == 0 {
		return
	}
	crf, err := GetCustomResponseFolder(responsesFolder)
	if err != nil {
		log.Printf("Failed to get custom responses - %v", err)
		return
	}
	for route, response := range crf.Routes {
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			w.Write(response)
		})
	}

	return GetRoutes()
}

func GetRoutes() (result []string) {
	routes := reflect.ValueOf(http.DefaultServeMux).Elem().FieldByName("m").MapRange()
	for routes.Next() {
		result = append(result, routes.Key().String())
	}
	return
}

func listFiles() {
	files, err := filepath.Glob("*")
	if err != nil {
		log.Printf("Failed to list files - %v", err)
		return
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
