package main

import (

	"fmt"
	"strconv"
	"github.com/urfave/cli"
)

// addCommand adds a new Task to the Database.
func addCommand(db Database, user string) cli.Command{

	command :=  cli.Command{
				    Name:    "add-task",
				    Aliases: []string{"a"},
				    Usage:   "add-task \"TASK_TITLE\" [FLAGS]. Adds a new Task to the Database.",
				    Flags: []cli.Flag{projectFlag()},
				    Action:  func(c *cli.Context) error {
				        	   	res := addTask(db, Task{
									Title:	c.Args().First(),
									//Users: []string{user},
									Project: c.String("project"),
								})
								printTask(res)
				        		return nil
			      			 },
		    	}

	return command
}

// listCommand lists all Tasks in the specified Project for the current User.
// If no Project is provided, only tasks without a Project will be listed.
// All tasks for a User, regardless of Project, can be listed by specifying "all"
// as the project.
func listCommand(db Database, user string) cli.Command{

	command :=  cli.Command{		    
				    Name:    "list-tasks",
				    Aliases: []string{"l"},
				    Usage:   "list-tasks [FLAGS]. Lists tasks for the current user that matches the provided FLAGS.",
				    Flags: []cli.Flag{projectFlag(), statusFlag()},
				    Action:  func(c *cli.Context) error {

				    			var matchParams []MatchParam

				    			//matchParams = append(matchParams, MatchParam{"tasks_by_user", user})

				    			if (len(c.String("status")) > 0) {
				    				matchParams = append(matchParams, MatchParam{"tasks_by_status", c.String("status")})
				    			}

				    			if (len(c.String("project")) > 0) {
				    				matchParams = append(matchParams, MatchParam{"tasks_by_project", c.String("project")})
				    			}
								
								res := listTasks(db, matchParams)
								printTasks(res)
				        		return nil
				     		 },
		    	}

	return command
}

// editCommand edits the specified Task with the provided flags.
func editCommand(db Database) cli.Command{

	command :=  cli.Command{
				    Name:    "edit-task",
				    Aliases: []string{"e"},
				    Usage:   "edit-task `ID` [FLAGS]. Update Task `ID` with the provided FLAGS.",
				    Flags: []cli.Flag{titleFlag(), projectFlag()},
				    Action:  func(c *cli.Context) error {
				      			id, _ := strconv.Atoi(c.Args().First())

				      			task := make(map[string]interface{})

				      			if (len(c.String("title")) > 0) {
				    				task["title"] = c.String("title")
				    			}

				    			if (len(c.String("project")) > 0) {
				    				task["project"] = c.String("project")
				    			}

				    			if(len(task) > 0) {

				    				task["id"] = id
				      				res := editTask(db, task)
				      				printTask(res)
				        	 	} else {

				        	 		fmt.Println("No edits specified!")
				        	 	}
				        	 	return nil
				      		 },
				}

	return command
}

// openCommand marks the Task as open.
func openCommand(db Database) cli.Command{

	command :=  cli.Command{
				    Name:    "open-task",
				    Aliases: []string{"o"},
				    Usage:   "open-task `ID`. Mark Task `ID` as open.",
				    Action:  func(c *cli.Context) error {
				      			id, _ := strconv.Atoi(c.Args().First())

				      			task := make(map[string]interface{})
				      			task["id"] = id
				      			task["status"] = "open"

				      			res := editTask(db, task)

								printTask(res)
				        	 	return nil
				      		 },
				}

	return command
}

// closeCommand marks the Task as closed.
func closeCommand(db Database) cli.Command{

	command :=  cli.Command{
				    Name:    "close-task",
				    Aliases: []string{"c"},
				    Usage:   "close-task `ID`. Mark task `ID` as closed.",
				    Action:  func(c *cli.Context) error {
				      			id, _ := strconv.Atoi(c.Args().First())

				      			task := make(map[string]interface{})
				      			task["id"] = id
				      			task["status"] = "closed"

				      			res := editTask(db, task)

								printTask(res)
				        		return nil
				    		 },
				}

    return command
}

// undoCommand removes the last event from a Task, or removes the
// Task if there is only one event.
func undoCommand(db Database) cli.Command{

	command :=  cli.Command{
				    Name:    "undo-task",
				    Aliases: []string{"u"},
				    Usage:   "undo-task `ID`. Revert the previous edit to Task `ID`, or delete it if there have been no edits.",
				    Action:  func(c *cli.Context) error {
				      			id, _ := strconv.Atoi(c.Args().First())

				      			res := undoTask(db, id)

								printTask(res)
				        	 	return nil
				      		 },
				}

	return command
}


