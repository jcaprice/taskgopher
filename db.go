package main

type Database interface {

	LookupTask(id int) Task
	PutTask(task Task) Task
	UpdateTask(task map[string]interface{}) Task
	UndoTask(id int) Task
	SearchTasks(matchParams []MatchParam) []Task
}

func InitDB(name string, conf map[string]string) Database {

	var db Database

	if(name == "faunadb") {

		db = InitFaunaDB(conf["secret"])
	} 

	return db
}