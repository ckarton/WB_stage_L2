package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("shell> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		commands := strings.Split(input, "|")
		if len(commands) > 1 {
			err := executePipeline(commands)
			if err != nil {
				fmt.Println("Error:", err)
			}
			continue
		}

		args := strings.Fields(commands[0])
		switch args[0] {
		case "cd":
			if len(args) < 2 {
				fmt.Println("cd: missing argument")
			} else {
				err := os.Chdir(args[1])
				if err != nil {
					fmt.Println(err)
				}
			}
		case "pwd":
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(dir)
			}
		case "echo":
			fmt.Println(strings.Join(args[1:], " "))
		case "kill":
			if len(args) < 2 {
				fmt.Println("kill: missing argument")
			} else {
				pid, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("kill: invalid PID")
				} else {
					proc, err := os.FindProcess(pid)
					if err != nil {
						fmt.Println(err)
					} else {
						err := proc.Kill()
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		case "ps":
			cmd := exec.Command("ps")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		case "\\quit":
			return
		default:
			err := executeCommand(args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	}
}

func executeCommand(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func executePipeline(commands []string) error {
	var cmds []*exec.Cmd

	for _, command := range commands {
		args := strings.Fields(command)
		if len(args) == 0 {
			continue
		}
		cmds = append(cmds, exec.Command(args[0], args[1:]...))
	}

	if len(cmds) == 0 {
		return nil
	}

	for i := 0; i < len(cmds)-1; i++ {
		cmds[i+1].Stdin, _ = cmds[i].StdoutPipe()
		cmds[i].Stderr = os.Stderr
	}

	cmds[len(cmds)-1].Stdout = os.Stdout
	cmds[len(cmds)-1].Stderr = os.Stderr

	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			return err
		}
	}

	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
