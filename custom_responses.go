package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
		Path:       folderName,
		RoutesData: make(map[string]CustomResponseData),
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
		method := "GET"
		if len(route.Method) > 0 && strings.Contains("GET POST PUT PATCH DELETE", strings.ToUpper(route.Method)) {
			method = strings.ToUpper(route.Method)
		}
		statusCode := uint(200)
		if route.StatusCode > 0 {
			statusCode = route.StatusCode
		}
		contentType := "application/json"
		if len(route.ContentType) == 0 {
			contentType, err = GetFileContentType(content)
			if err != nil {
				log.Printf("Failed to get content type for file %s - %s", route.SourceFile, err)
			} else {
				log.Printf("Detected content type for file %s - %s", route.SourceFile, contentType)
				route.ContentType = contentType
			}
		}

		crf.RoutesData[route.Path] = CustomResponseData{
			Path:        route.Path,
			Content:     content,
			Method:      method,
			StatusCode:  statusCode,
			ContentType: route.ContentType,
		}

		log.Printf("Custom response: %s %s %d (%s) %s", route.Path, method, statusCode, filepath.Base(route.SourceFile), route.ContentType)
	}
	if len(crf.RoutesData) == 0 {
		return nil, fmt.Errorf("No valid custom responses found in %s", folderName)
	}
	return crf, nil
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

	if len(customResponsesFolder) == 0 {
		routesFile, err := filepath.Abs("./custom_responses/routes.json")
		if err == nil {
			if _, err := os.Stat(routesFile); err == nil {
				customResponsesFolder = filepath.Dir(routesFile)
				source = "default"
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

func SetupHttpCustomHandlersRouter(router *http.ServeMux) (routes []string) {
	responsesFolder := GetArgResponsesFolder()
	if len(responsesFolder) == 0 {
		return
	}
	crf, err := GetCustomResponseFolder(responsesFolder)
	if err != nil {
		log.Printf("Failed to get custom responses - %v", err)
		return
	}
	for route, response := range crf.RoutesData {
		router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			requestCount.Inc()
			w.Header().Set("Content-Type", response.ContentType)
			w.WriteHeader(int(response.StatusCode))
			w.Write(response.Content)
		})
	}

	return GetRoutesRouter(router)
}

func GetRoutesRouter(router *http.ServeMux) (result []string) {
	routes := reflect.ValueOf(router).Elem().FieldByName("m").MapRange()
	for routes.Next() {
		result = append(result, routes.Key().String())
	}
	return
}
