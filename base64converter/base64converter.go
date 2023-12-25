package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func encodeFile(inputPath, outputPath string) error {
	inputData, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	byted := []byte(inputData)
	hash := sha256.Sum256(byted)
	str := hex.EncodeToString(hash[:])
	fmt.Println(str)
	encodedData := base64.StdEncoding.EncodeToString(inputData)

	if outputPath == "" {
		outputPath = inputPath + ".out"
	}

	return ioutil.WriteFile(outputPath, []byte(encodedData), 0644)
}

func decodeFile(inputPath, outputPath string) (error, string) {
	inputData, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err, ""
	}

	decodedData, err := base64.StdEncoding.DecodeString(string(inputData))
	if err != nil {
		return err, ""
	}

	byted := []byte(decodedData)
	hash := sha256.Sum256(byted)
	str := hex.EncodeToString(hash[:])
	fmt.Println(str)

	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = strings.TrimSuffix(inputPath, ext)
		if ext == "" {
			ext = ".out"
		}
		outputPath += "-decoded" + ext
	}

	return ioutil.WriteFile(outputPath, []byte(decodedData), 0644), str
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: encoder64 <command> [options]")
		fmt.Println("Commands: encode, decode")
		fmt.Println("Please, enter full path to the source files")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "encode":
		encodeCmd := flag.NewFlagSet("encode", flag.ExitOnError)
		inputFlag := encodeCmd.String("i", "", "Input file")
		outputFlag := encodeCmd.String("o", "", "Output file")
		encodeCmd.Parse(os.Args[2:])

		inputPath := *inputFlag
		outputPath := *outputFlag

		if inputPath == "" {
			inputPath = os.Args[2]
		}

		if outputPath == "" {
			ext := filepath.Ext(inputPath)
			outputPath = strings.TrimSuffix(inputPath, ext)
			if ext == "" {
				ext = ".out"
			}
			outputPath += ".out"
		}

		err := encodeFile(inputPath, outputPath)
		if err != nil {
			fmt.Println("Error encoding file:", err)
			os.Exit(1)
		}
		fmt.Println("File encoded successfully.")
	case "decode":
		decodeCmd := flag.NewFlagSet("decode", flag.ExitOnError)
		inputFlag := decodeCmd.String("i", "", "Input file")
		outputFlag := decodeCmd.String("o", "", "Output file")
		hashFlag := decodeCmd.String("h", "", "Hash")
		decodeCmd.Parse(os.Args[2:])

		inputPath := *inputFlag
		outputPath := *outputFlag
		yourHash := *hashFlag

		if inputPath == "" {
			inputPath = os.Args[2]
		}

		if outputPath == "" {

			ext := filepath.Ext(inputPath)
			outputPath = strings.TrimSuffix(inputPath, ext)
			if ext == "" {
				ext = ".out"
			}
			outputPath += ".out"
		}

		err, hash_d := decodeFile(inputPath, outputPath)

		if yourHash == "" {
			yourHash = os.Args[3]
			fmt.Println(yourHash)
			fmt.Println(hash_d)
		}

		if yourHash == hash_d {
			fmt.Println("Hashes matched")
		}

		if err != nil {
			fmt.Println("Error decoding file:", err)
			os.Exit(1)
		}
		fmt.Println("File decoded successfully.")
	default:
		fmt.Println("Unknown command. Please, use 'encode' or 'decode'.")
		os.Exit(1)
	}
}
