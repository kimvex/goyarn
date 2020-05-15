package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	init := flag.Bool("init", false, "Initializa program")
	install := flag.Bool("install", false, "install programs")
	uninstall := flag.Bool("uninstall", false, "uninstall programs")
	version := flag.Bool("version", false, "Version of goyarn")
	depInstallAll := flag.Bool("i", false, "Install all dependencies of package.json")

	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filePackage := fmt.Sprintf("%v/package.json", pwd)

	if *init == true {
		_, errIf := os.Open(filePackage)

		if errIf != nil {
			arrSplitProject := strings.Split(pwd, "/")
			projectName := fmt.Sprintf("%v", arrSplitProject[len(arrSplitProject)-1])
			in := make(map[string]interface{})
			depsCreate := make(map[string]interface{})
			in["version"] = "v0.0.1"
			in["name"] = projectName
			in["dependencies"] = depsCreate
			packageEdit, _ := json.MarshalIndent(in, "", " ")
			_ = ioutil.WriteFile("package.json", packageEdit, 0644)
		}
		return
	}

	file, errIf := os.Open(filePackage)

	b1 := make([]byte, 32*1024)
	n1, e := file.Read(b1)

	var result map[string]interface{}
	json.Unmarshal([]byte(string(b1[:n1])), &result)

	if len(result) < 1 {
		result["version"] = "v0.0.1"
		interfaceDepsNew := make(map[string]interface{})
		result["dependencies"] = interfaceDepsNew
		packageNew, _ := json.MarshalIndent(result, "", " ")
		_ = ioutil.WriteFile("package.json", packageNew, 0644)
		return
	}

	if e != nil {
		log.Fatal(e)
	}

	if errIf != nil {
		fmt.Println("package.json not found")
		return
	}

	if *depInstallAll == true {
		Package(result)
	}

	if *install == true {
		if len(flag.Args()) < 1 {
			fmt.Println("not package empty")
			return
		}

		packageInstall := fmt.Sprintf("go get -u %v", flag.Args()[0])
		exeCommand(packageInstall)

		arrSplit := strings.Split(flag.Args()[0], "/")
		lastText := fmt.Sprintf("%v", arrSplit[len(arrSplit)-1])
		fmt.Println(result["dependencies"], result)
		if result["dependencies"] == nil {
			interfaceDeps := make(map[string]interface{})
			result["dependencies"] = interfaceDeps
			deps := result["dependencies"].(map[string]interface{})
			deps[lastText] = flag.Args()[0]
		} else {
			deps, _ := result["dependencies"].(map[string]interface{})
			deps[lastText] = flag.Args()[0]
		}
		fmt.Println(result)

		packageEdit, _ := json.MarshalIndent(result, "", " ")
		_ = ioutil.WriteFile("package.json", packageEdit, 0644)
	}

	if *uninstall == true {
		if result["dependencies"] == nil {
			fmt.Println("Not found dependencies")
			return
		}

		FoundPackage(result, flag.Args()[0])
	}

	if *version == true {

	}
}

/*Package install pakcages of a package.json*/
func Package(result map[string]interface{}) {
	for PrincipalKey, PrincipalValue := range result {
		fmt.Println(PrincipalKey)
		switch PrincipalKey {
		case "dependencies":
			deps := PrincipalValue.(map[string]interface{})
			for keydeps, valuedeps := range deps {
				fmt.Println(keydeps, valuedeps)
				packageInstall := fmt.Sprintf("go get -u %v", valuedeps)
				exeCommand(packageInstall)
			}
		}
	}
}

func FoundPackage(result map[string]interface{}, packageNameValue string) {
	for PrincipalKey, PrincipalValue := range result {
		switch PrincipalKey {
		case "dependencies":
			deps := PrincipalValue.(map[string]interface{})
			for keydeps, valuedeps := range deps {
				if packageNameValue == keydeps {
					packageUninstall := fmt.Sprintf("go clean -i -n %v", valuedeps)
					exeCommand(packageUninstall)
					delete(deps, keydeps)
				}
			}
		}
	}

	packageEdit, _ := json.MarshalIndent(result, "", " ")
	_ = ioutil.WriteFile("package.json", packageEdit, 0644)
}

func exeCommand(command string) {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	errorCmd := cmd.Run()

	if errorCmd != nil {
		fmt.Println(errorCmd, os.Stderr)
		return
	}
}
