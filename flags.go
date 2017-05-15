package main

import (

	"github.com/urfave/cli"
)

func projectFlag() cli.Flag{

	flag := cli.StringFlag{
		Name: "project, p",
		Value: "",
		Usage: "The project this task belongs to.",
	}

	return flag
}

func statusFlag() cli.Flag{

	flag := cli.StringFlag{
		Name: "status, s",
		Value: "",
		Usage: "The status of the task (open, closed, active.)",
	}

	return flag
}

func titleFlag() cli.Flag{

	flag := cli.StringFlag{
		Name: "title, t",
		Value: "",
		Usage: "The title of the task (Only used in editting tasks.)",
	}

	return flag
}