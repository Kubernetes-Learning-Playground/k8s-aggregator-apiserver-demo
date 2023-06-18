package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// kubectl 命令调用方式： kubectl --kubeconfig ./resources/config_local --insecure-skip-tls-verify=true get myingresses test

func main() {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	args := []string{"--kubeconfig", "./resources/config_local", "--insecure-skip-tls-verify=true"}
	fmt.Println(args)
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}
	fmt.Println("kubectl", args)
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

}
