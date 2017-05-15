package main

import (
	"os"
	"fmt"
	"bytes"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/fatih/structs"
)

// printTask prints a single task in a Key/Value format
// using olekukonko's tablewriter package.
func printTask(task Task) {

	if task.Id != 0 {

		data := [][]string{}

		s := structs.New(task)

		names := s.Names()[1:]

		for _, name := range names {

			// Get the Field struct for the "Name" field
			field := s.Field(name)

			switch t := field.Value().(type) {
	    	case int:
	    		data = append(data, []string{name, strconv.Itoa(t)})
	     	case string:
	     		data = append(data, []string{name, t})
	     	case []string:
	     		var stringBuffer bytes.Buffer

	     		if len(t) > 0 {
		     		for _, str := range t {

		     			stringBuffer.WriteString(str + ", ")
		     		}

		     		data = append(data, []string{name, stringBuffer.String()[:len(stringBuffer.String()) - 2]})
				} else {
					data = append(data, []string{name, ""})
				}
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetRowLine(true)

		for _, v := range data {
		   	table.Append(v)
		}

		table.Render()
	} else {
		fmt.Println("Task does not exist!")
	}
}

// printTasks prints a list of tasks in tabular format using
// olekukonko's tablewriter package.
func printTasks(tasks []Task) {

	data := [][]string{}

	if(len(tasks) > 0) {

		for _, task := range tasks {

				v := structs.Values(task)
				s := make([]string, len(v))

				for i := range v {

					switch t := v[i].(type) {
		        	case int:
		          		s[i] = strconv.Itoa(t)
		         	case string:
		         		if len(t) > 30 {

		         			s[i] = t[:29]
		         		} else {

		         			s[i] = t
		         		}
		         	case []string:
		         		var buffer bytes.Buffer

		         		for _, str := range t {
		         			buffer.WriteString(str + ", ")
		         		}

		         		strs := buffer.String()[:len(buffer.String()) - 2]

		         		if len(strs) > 30 {

		         			s[i] = strs[:29]
		         		} else {

		         			s[i] = strs
		         		}
		         	default:
		         		panic("unhandled type")
		         	}
				}

				data = append(data, s)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(structs.Names(tasks[0]))

		for _, v := range data {
	    	table.Append(v)
		}

		table.Render()
	} else {

		fmt.Println("There are no tasks to show!")
	}
}