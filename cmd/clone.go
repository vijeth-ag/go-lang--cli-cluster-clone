/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

type AnnotationStruct struct {
	Config string `json:"kubectl.kubernetes.io/last-applied-configuration"`
}

type Annotations struct {
	Annotations AnnotationStruct `json:"annotations"`
}

type ItemStruct struct {
	Metadata Annotations `json:"metadata"`
}

type Items struct {
	Items []ItemStruct `json:"items"`
}

const ShellToUse = "bash"

func getKubeConfigPath() string {
	os := runtime.GOOS

	var kubeConfigBasePath string

	switch os {
	case "windows":
		kubeConfigBasePath = `%USERPROFILE%\.kube\`
	case "darwin":
		kubeConfigBasePath = `$HOME/.kube/`
	case "linux":
		kubeConfigBasePath = `~/.kube/`
	default:
		fmt.Printf("%s.\n", os)
	}

	return kubeConfigBasePath
}

func getConfigPath(config string) string {
	var configPath string
	if strings.Contains(config, "/") {
		// Absolute path
		configPath = config
	} else {
		// get OS specifc  kube base path and append config file name
		configPath = getKubeConfigPath() + config
	}

	return configPath
}

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func getSourceResources(kubeConfig string) error {

	cmd := "kubectl get deployment -o json --kubeconfig=" + kubeConfig + "> source_resources.json"

	err, _, _ := Shellout(cmd)
	if err != nil {
		log.Printf("error: %v\n", err)
		return err
	}
	return nil
}

func getLastAppliedConfig() error {
	file, err := ioutil.ReadFile("source_resources.json")

	if err != nil {
		log.Fatal(err)
		return err
	}

	res := Items{}

	json.Unmarshal([]byte(file), &res)
	data := res.Items[0].Metadata.Annotations.Config

	writeErr := os.WriteFile("source_resources.json", []byte(data), 0644)
	if writeErr != nil {
		log.Println(err)
		return err
	}

	return nil
}

func applyConfig(kubeConfig string) error {
	cmd := "kubectl apply -f source_resources.json --kubeconfig=" + kubeConfig

	err, _, _ := Shellout(cmd)
	if err != nil {
		log.Printf("error: %v\n", err)
		return err
	}
	return nil
}

var sourceKubeConfig string
var destKubeConfig string

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "A brief description of your command",
	Long:  `A longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Missing kube config file")
			fmt.Println("Usage: clone <sourceKubeConfig>  <destinationKubeConfig>")
			return
		}

		sourceKubeConfig = getConfigPath(args[0])
		destKubeConfig = getConfigPath(args[1])

		err := getSourceResources(sourceKubeConfig)
		if err != nil {
			log.Println("error getting source resources")
		}
		configReadErr := getLastAppliedConfig()
		if configReadErr != nil {
			log.Println("error parsing resource file")
		}
		applyErr := applyConfig(destKubeConfig)
		if applyErr != nil {
			log.Println("error applying config to destination")
		}
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
