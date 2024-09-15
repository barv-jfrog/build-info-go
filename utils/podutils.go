package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Pod struct {
	Name         string `json:"name"`
	Version      string `json:"version,omitempty"`
	Dependencies []Pod  `json:"dependencies,omitempty"`
}

func GetPodDependencies() (map[string][]string, []string, error) {
	podsSection, err := getPodsSection("Podfile.lock")
	if err != nil {
		return nil, nil, err
	}
	depGraph, depDirect := extractDependencies(podsSection)
	return depGraph, depDirect, nil
}

func extractDependencies(podsSection string) (map[string][]string, []string) {
	graph := make(map[string][]string)
	var directDependencies []string
	lines := strings.Split(podsSection, "\n")
	var currentPod string
	for index, _ := range lines {
		line := strings.TrimSpace(lines[index])
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, " ", 3)
			currentPod = parts[1]
			directDependencies = append(directDependencies, currentPod)
			graph[currentPod] = []string{}
			index += 1
			line = strings.TrimSpace(lines[index])
			for !strings.Contains(line, ":") && index < len(lines) {
				parts = strings.SplitN(line, " ", 3)
				subPod := parts[1]
				graph[currentPod] = append(graph[currentPod], subPod)
				index += 1
				line = strings.TrimSpace(lines[index])
			}
		}
	}
	return graph, directDependencies
}

func getPodsSection(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	var podsData string
	// Read the file line by line.
	scanner := bufio.NewScanner(file)
	capture := false
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line is "PODS:" to start capturing.
		if strings.HasPrefix(line, "PODS:") {
			capture = true
			continue
		}
		if capture {
			podsData += line + "\n"
		}
		if strings.HasPrefix(line, "DEPENDENCIES:") {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return podsData, nil
}

//
//func parseDependenciesToGraph(packages []Pod) (map[string][]string, []string, error) {
//	// Create packages map.
//	packagesMap := map[string][]string{}
//	allSubPackages := map[string]bool{}
//	for _, pkg := range packages {
//		var subPackages []string
//		for _, subPkg := range pkg.Dependencies {
//			subPkgFullName := subPkg.Key + ":" + subPkg.InstalledVersion
//			subPackages = append(subPackages, subPkgFullName)
//			allSubPackages[subPkgFullName] = true
//		}
//		packagesMap[pkg.Package.Key+":"+pkg.Package.InstalledVersion] = subPackages
//	}
//
//	var topLevelPackagesList []string
//	for pkgName := range packagesMap {
//		if !allSubPackages[pkgName] {
//			topLevelPackagesList = append(topLevelPackagesList, pkgName)
//		}
//	}
//	return packagesMap, topLevelPackagesList, nil
//}
