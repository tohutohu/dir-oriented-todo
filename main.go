package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
)

type Config struct {
	SaveDir  string
	FilePath string
}

type Todo struct {
	Body string    `json:"body"`
	Time time.Time `json:"time"`
}

type TodoDir struct {
	Path  string `json:"path"`
	Todos []Todo `json:"todos"`
}

type TodoList map[string]*TodoDir

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

	if c.Bool("a") {
		for _, todos := range todoList {
			if len(todos.Todos) > 0 {
				fmt.Println(todos.Path)
				printTodos(todos.Todos)
				fmt.Print("\n")
			}
		}
		return
	}

	if err != nil {
		panic(err)
	}
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	todoDir, ok := todoList[curDir]
	if ok && len(todoDir.Todos) > 0 {
		printTodos(todoDir.Todos)
	} else {
		fmt.Println("Current directory has no todo!")
	}
}

func cmdAdd(c *cli.Context) {
	var config Config
	todoList, err := config.load()
	if err != nil {
		panic(err)
	}
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var todoDir *TodoDir
	todoDir, ok := todoList[curDir]
	if !ok {
		todoDir = new(TodoDir)
	} else if len(todoDir.Todos) > 0 {
		printTodos(todoDir.Todos)
		fmt.Println()
	}
	todo := Todo{}
	fmt.Print("what todo?> ")
	sc := bufio.NewScanner(os.Stdin)
	if sc.Scan() {
		todo.Body = sc.Text()
	}
	todo.Time = time.Now()
	todoDir.Todos = append(todoDir.Todos, todo)
	todoDir.Path = curDir
	todoList[curDir] = todoDir
	todoList.Save(config.FilePath)
	fmt.Println("Saved!")
}

func cmdDelete(c *cli.Context) {
	var config Config
	todoList, err := config.load()
	if err != nil {
		panic(err)
	}
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	todoDir, ok := todoList[curDir]
	if !ok || len(todoDir.Todos) == 0 {
		fmt.Println("Current directory has no todo!")
		return
	}
	printTodos(todoDir.Todos)

	fmt.Print("ID> ")
	var ID int
	fmt.Scan(&ID)
	if ID-1 < len(todoDir.Todos) {
		todoDir.Todos = append(todoDir.Todos[:(ID-1)], todoDir.Todos[ID:]...)
		todoList[curDir] = todoDir
		todoList.Save(config.FilePath)
		fmt.Println("Deleted!")
	} else {
		fmt.Println("Invalid ID")
	}
}

func printTodos(todos []Todo) {
	fmt.Println(" ID             Body                                                         Date")
	for i, todo := range todos {
		fmt.Printf("%3d             %-*s %s\n", i+1, 60-(len(todo.Body)-len([]rune(todo.Body)))/2, todo.Body, todo.Time.Format("01/2 15:04"))
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
		{
			Name:   "add",
			Action: cmdAdd,
		},
		{
			Name:   "delete",
			Action: cmdDelete,
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "show all todo",
		},
	}

	app.Action = cmdMain

	app.Run(os.Args)
}
