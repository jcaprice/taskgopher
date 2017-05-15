package main

type Task struct {
	Id int `fauna:"id"`
	Title string `fauna:"title"`
	Project string `fauna:"project"`
	Status string `fauna:"status"`
	//Users []string `fauna:"users"`
}

// SetId assigns the ID used to look up the task.
func (t *Task) SetId(id int) {
	t.Id = id
}

// SetStatus assigns the status to the task.
func (t *Task) SetStatus(status string) {
	t.Status = status
}

// getTask
func getTask(db Database, id int) Task {

	return db.LookupTask(id)
}

// addTask adds a new task to the Database
// Task ID is assigned to the task before writing to the database
// and is used for future single-task operations.
func addTask(db Database, task Task) Task {

	task.SetStatus("open")

	return db.PutTask(task)
}

// listTasks returns a list of tasks based on the Match Parameters.
func listTasks(db Database, matchParams []MatchParam) []Task {

	return db.SearchTasks(matchParams)
}

// editTask takes a Task object and updates the database with it's new values.
func editTask(db Database, task map[string]interface{}) Task {

	return db.UpdateTask(task)
}

//undoTask reverts 
func undoTask(db Database, id int) Task {

	return db.UndoTask(id)
}