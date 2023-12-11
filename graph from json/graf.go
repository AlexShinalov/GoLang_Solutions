package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Graph struct {
	Nodes []struct {
		ID int `json:"id"`
	} `json:"nodes"`
	Edges []struct {
		Source int `json:"source"`
		Target int `json:"target"`
	} `json:"edges"`
}

func main() {

	file, err := os.Open("graf.json") // Укажите имя вашего JSON-файла
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	var jsonData map[string]Graph
	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	graph := jsonData["graph"]

	// Указываем начальный и конечный узлы для поиска пути
	var startNode int
	var endNode int
 fmt.print('Введите начальную ноду')
 fmt.Scan(&startNode)
 fmt.print('Введите конечную ноду')
 fmt.Scan(&endNode)

	
	path := findPath(graph, startNode, endNode, []int{})
	if len(path) > 0 {
		fmt.Printf("Найден путь от узла %d к узлу %d: %v\n", startNode, endNode, path)
	} else {
		fmt.Printf("Путь от узла %d к узлу %d не найден\n", startNode, endNode)
	}
}

// DFS
func findPath(graph Graph, currentNode int, endNode int, visited []int) []int {
	if currentNode == endNode {
		return []int{currentNode}
	}

	visited = append(visited, currentNode)

	for _, edge := range graph.Edges {
		if edge.Source == currentNode {
			if !isVisited(visited, edge.Target) {
				newPath := findPath(graph, edge.Target, endNode, visited)
				if len(newPath) > 0 {
					return append([]int{currentNode}, newPath...)
				}
			}
		}
	}

	return []int{}
}

func isVisited(visited []int, node int) bool {
	for _, v := range visited {
		if v == node {
			return true
		}
	}
	return false
}
