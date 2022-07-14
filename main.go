package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	if args["operation"] == "add" {
		if args["item"] == "" {
			return fmt.Errorf("-item flag has to be specified")
		}
		bytes, _ := ioutil.ReadFile(args["fileName"])
		users := []User{}
		json.Unmarshal(bytes, &users)
		item := User{}
		json.Unmarshal([]byte(args["item"]), &item)
		for i := 0; i < len(users); i++ {
			if users[i].Id == item.Id {
				writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", item.Id)))
				return nil
			}
		}
		users = append(users, item)
		bytes, _ = json.Marshal(&users)
		ioutil.WriteFile(args["fileName"], bytes, 0644)
	} else if args["operation"] == "list" {
		bytes, _ := ioutil.ReadFile(args["fileName"])
		writer.Write(bytes)
	} else if args["operation"] == "findById" {
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		bytes, _ := ioutil.ReadFile(args["fileName"])
		users := []User{}
		json.Unmarshal(bytes, &users)
		for i := 0; i < len(users); i++ {
			if users[i].Id == args["id"] {
				bytes, _ = json.Marshal(users[i])
				writer.Write(bytes)
				return nil
			}
		}
	} else if args["operation"] == "remove" {
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		bytes, _ := ioutil.ReadFile(args["fileName"])
		users := []User{}
		json.Unmarshal(bytes, &users)
		for i := 0; i < len(users); i++ {
			if users[i].Id == args["id"] {
				users = append(users[:i], users[i+1:]...)
				bytes, _ = json.Marshal(&users)
				ioutil.WriteFile(args["fileName"], bytes, 0644)
				return nil
			}
		}
		writer.Write([]byte(fmt.Sprintf("Item with id %s not found", args["id"])))
	} else {
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
	return nil
}

func parseArgs() Arguments {
	operation := flag.String("operation", "", "operation")
	item := flag.String("item", "", "item")
	fileName := flag.String("fileName", "", "fileName")
	flag.Parse()

	return Arguments{
		"operation": *operation,
		"item":      *item,
		"fileName":  *fileName,
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
