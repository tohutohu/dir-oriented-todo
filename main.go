package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

type Config struct {
	SaveDir  string
	FilePath string
}

type Todo struct {
	Path string `json:"path"`
	Body string `json:"body"`
}

type TodoList map[string]*Todo

func (config *Config) load() (TodoList, error) {
	dir := filepath.Join(os.Getenv("HOME"), ".config", "todo")
	file := filepath.Join(dir, "todoData.json")
	_, err := os.Stat(file)

	config.SaveDir = dir
	config.FilePath = file

	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		os.MkdirAll(dir, 0700)
		err := ioutil.WriteFile(file, []byte("{}"), os.FileMode(0600))

		if err != nil {
			return nil, err
		}
	}

	var todoList TodoList
	body, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &todoList)
	if err != nil {
		return nil, err
	}

	return todoList, nil
}

func (todoList *TodoList) Save(path string) error {
	bytes, err := json.Marshal(todoList)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, os.FileMode(0600))
}

func cmdMain(c *cli.Context) {
	var config Config
	todoList, err := config.load()
	if err != nil {
		panic(err)
	}
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println("Hello")

	body, ok := todoList[curDir]
	if ok {
		fmt.Println(body)
	} else {
		todo := new(Todo)
		fmt.Print("what todo?> ")
		fmt.Scan(&(todo.Body))
		todo.Path = curDir
		todoList[curDir] = todo
		todoList.Save(config.FilePath)
		fmt.Println("Saved!")
	}
}

func test(c *cli.Context) {
	fmt.Println("pjoij")
	os.Chdir("/home/to-hutohu")
}

func main() {
	app := cli.NewApp()

	app.Name = "Todo"
	app.Usage = "Check todo current directory"

	app.Author = "to-hutohu"
	app.Email = "tohu.soy@gmail.com"
	app.Commands = []cli.Command{
		{
			Name:   "test",
			Action: test,
		},
	}

	app.Action = cmdMain

	app.Run(os.Args)
}
