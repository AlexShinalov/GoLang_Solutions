package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type StoryNode struct {
	Description string
	Option      string
	Options     map[int]string // ключ - номер варианта, значение - файл с описанием действия
	Option_TT   map[int]string
}

func readNodeFromFile(filename string) (*StoryNode, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var nodeID int
	var description string
	var descriptionBuilder strings.Builder
	var optionBulder strings.Builder
	options := make(map[int]string)
	opt := make(map[int]string)

	var readingOptions bool

	// Чтение содержимого файла
	for scanner.Scan() {
		line := scanner.Text()

		if line == "&" {
			readingOptions = true
			continue
		}

		if !readingOptions {
			descriptionBuilder.WriteString(line + "\n")
			if len(line) == 0 {
				continue
			}

			if nodeID == 0 {
				nodeID, err = strconv.Atoi(line)
				if err != nil {
					return nil, fmt.Errorf("failed to parse node ID in file %s", filename)
				}
			} else if description == "" {
				description = line
			}
		} else {
			parts := strings.Split(line, "@")
			optionBulder.WriteString(parts[0] + " " + parts[1] + "\n")
			if len(parts) < 2 {
				return nil, fmt.Errorf("invalid format in file %s", filename)
			}

			optionNum, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("failed to parse option number in file %s", filename)
			}

			options[optionNum] = parts[1]
			opt[optionNum] = parts[2]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &StoryNode{Description: descriptionBuilder.String(), Option: optionBulder.String(), Options: options, Option_TT: opt}, nil
}

func main() {
	nodes := make(map[int]*StoryNode)

	// Чтение файлов и создание узлов
	for i := 1; i <= 12; i++ {
		filename := fmt.Sprintf("%d.txt", i)
		node, err := readNodeFromFile(filename)
		if err != nil {
			fmt.Println("Error reading node:", err)
			continue
		}
		nodes[i] = node
	}

	currentNodeID := 1

	// Начало истории
	for {
		currentNode := nodes[currentNodeID]
		fmt.Println(currentNode.Description)

		if len(currentNode.Options) == 0 {
			fmt.Println("Конец истории.")
			break
		}

		fmt.Println("Выберите действие:")
		/*for id, text := range currentNode.Options {
			fmt.Printf("%d. %s\n", id, text)
		}
		*/
		fmt.Println(currentNode.Option)
		var choice int
		fmt.Scanln(&choice)

		nextNodeFile, exists := currentNode.Option_TT[choice]
		if !exists {
			fmt.Println("Некорректный выбор. Попробуйте снова.")
			continue
		}

		newNode, err := readNodeFromFile(nextNodeFile)

		if err != nil {
			fmt.Println("Ошибка при чтении узла:", err, nextNodeFile)
			continue
		}

		// Обновляем ID текущего узла

		for id, node := range nodes {
			if node == newNode {

				currentNodeID = id
				fmt.Println(currentNodeID, " ")
				break
			}

		}
		id := strings.Split(nextNodeFile, ".")
		idd, _ := strconv.Atoi(id[0])
		currentNodeID = idd
	}
}
