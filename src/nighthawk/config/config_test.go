package config

import (
	"fmt"
	"testing"
) 
  


func TestConfig(t *testing.T) {
	fmt.Println("===== BEFORE LOADING CONFIG =====")
	ShowHashSetConfig()
	ShowStackDbConfig()
	ShowStandaloneConfig()

	err := LoadNighthawkConfig()
	if err != nil {
		fmt.Println(err)
	}	

	fmt.Println("===== AFTER LOADING CONFIG =====")
	ShowHashSetConfig()
	ShowStackDbConfig()
	ShowStandaloneConfig()
}


func HashSetConfigTest() {
	ShowHashSetConfig()
	SetConfigFile("/etc/nighthawk.json")
	err := LoadNighthawkConfig()
	if err != nil {
		fmt.Println(err)
	}
	ShowHashSetConfig()
}


func ShowHashSetConfig() {
	fmt.Println("HashSetIndex(): ", HashSetIndex())
	fmt.Println("HashSetEnabled(): ", HashSetEnabled())
	fmt.Println("HashSetAvailable(): ", HashSetAvailable())
	fmt.Println("HashSetChecked(): ", HashSetChecked())
}



func StackDbConfigTest() {
	ShowStackDbConfig()
}

func ShowStackDbConfig() {
	fmt.Println("StackDbIndex(): ", StackDbIndex())
	fmt.Println("StackDbPath(): ", StackDbPath())
	fmt.Println("StackDbEnabled(): ", StackDbEnabled())
	fmt.Println("StackDbAvailable(): ", StackDbAvailable())
	fmt.Println("StackDbChecked(): ", StackDbChecked())
}


func ShowStandaloneConfig() {
	fmt.Println("IsStandalone(): ", IsStandalone())
}