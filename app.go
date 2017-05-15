package main

import (
	// Golang Packages
	"fmt"
	"os"
	"strconv"

	// Community Packages
	"github.com/urfave/cli"
	"github.com/spf13/viper"
)

func main() {

	// Config 
	viper.SetConfigName("taskgopher")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	  if err != nil {
	   	fmt.Println("There was an error reading your configuration file.")
	  } else {

	  	user := viper.GetString("authentication.user")
	  	secret := viper.GetString("authentication.secret")
	  	backend := viper.GetString("general.backend")

	  	db := InitDB(backend, 
	  			map[string]string{

	  				"secret": secret,
	  			},
	  		)

	  	// Command Line
		app := cli.NewApp()

		app.Name = "taskgopher"
		app.Usage = "Task Manager written in Go"

		// Command Construction
		app.Commands = []cli.Command{
			addCommand(db, user),
			listCommand(db, user),
			editCommand(db),
			closeCommand(db),
			openCommand(db),
			undoCommand(db),
		}

		app.Action = func(c *cli.Context) error {

    		if c.NArg() > 0 {
      			
      			id, err := strconv.Atoi(c.Args()[0])

      			if err != nil {
      				panic(err)
   				}

				task := getTask(db, id)
				printTask(task)
    		}

    		return nil
    	}

		app.Run(os.Args)
	}
}



