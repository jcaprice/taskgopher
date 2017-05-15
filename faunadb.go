package main

import (

	f "github.com/fauna/faunadb-go/faunadb"
)

type FaunaDatabase struct {

	Client *f.FaunaClient
}

// InitFaunaDB: Initialize a FaunaDB instance of the Database interface.
func InitFaunaDB(secret string) Database {

	return &FaunaDatabase{

		Client: f.NewFaunaClient(secret),
	}
}

// LookupTask: Look up a Task in FaunaDB with the given Task ID.
func (fDB *FaunaDatabase) LookupTask (id int) Task {

	res, _ := fDB.Client.Query(
    			f.Select(f.Arr{"data"},
    				f.Get(
        				f.MatchTerm(
            				f.Index("tasks_by_id"), id))))

	if res == nil {
		return Task{}
	}

	var task Task

	res.Get(&task)

	return task
}

// PutTask: Put a new Task to FaunaDB.
func (fDB *FaunaDatabase) PutTask (task Task) Task {

	task.Id = IncrementCounter(fDB, "classes/global/1")

	res, err := fDB.Client.Query(
					f.Select(f.Arr{"data"},
						f.Create(
							f.Class("tasks"),
							f.Obj{"data": task})))
	if err != nil {
		panic(err)
	}

	var result Task

	res.Get(&result)

	return result
}

// UpdateTask: Update a Task in FaunaDB with the provided changes.
func (fDB *FaunaDatabase) UpdateTask (task map[string]interface{}) Task {

	ref := LookupTaskRef(fDB, "tasks_by_id", task["id"].(int))

	update, err := fDB.Client.Query(
								f.Select(f.Arr{"data"},
									f.Update(ref, f.Obj{ "data" : task})))

	if err != nil {
		panic(err)
	}

	var result Task

	update.Get(&result)

	return result
}

// SearchTasks: Query FaunaDB using the provided list of parameters for matching. 
// If more than one match parameter is provided, the resulting matches are 
// intersected.
func (fDB *FaunaDatabase) SearchTasks (matchParams []MatchParam) []Task {
	
	var list []Task

	res, err := fDB.Client.Query(
					f.Map(
						f.Paginate(
								matchBuilder(matchParams)),
						f.Lambda("data", 
							f.Select(f.Arr{"data"}, 
								f.Get(f.Var("data"))))))
	if err != nil {
		panic(err)
	}

	if err := res.At(f.ObjKey("data")).Get(&list); err != nil {
		panic(err)
	}

	return list
}

// UndoTask: Revert a previous change to a task. If the task was just added,
// or all events but the initial creation of the task are removed,
// the task will be deleted.
func (fDB *FaunaDatabase) UndoTask (id int) Task {

	taskObject := LookupTaskObj(fDB, "tasks_by_id", id)

	ts, _ := fDB.Client.Query(f.Select("ts", taskObject))
	ref, _ := fDB.Client.Query(f.Select("ref", taskObject))

	_, err := fDB.Client.Query(f.Remove(ref, ts, f.ActionCreate))

	if err != nil {
		panic(err)
	}

	return fDB.LookupTask(id)
}

// IncrementCounter: Increment a Counter in FaunaDB.
func IncrementCounter (fDB *FaunaDatabase, counter string) int {

	var nextId int

	nextIdResult, err := fDB.Client.Query(
					f.Select(f.Arr{"data", "id"}, 
						f.Update(
							f.Ref(counter), f.Obj{ "data" : 
								f.Obj{ "id" : 
									f.Add(
										f.Select(f.Arr{"data", "id"}, 
											f.Get(f.Ref(counter))), 1) }})))

	if err != nil {
			panic(err)
	}

	if err := nextIdResult.Get(&nextId); err != nil {
		panic(err)
	}

	return nextId
}

/////////////////////
// FaunaDB Helpers //
/////////////////////

// LookupTaskRef: Look up a Task ref in FaunaDB using the given index and id.
func LookupTaskRef(fDB *FaunaDatabase, index string, id int) f.Value {

	res, err := fDB.Client.Query(
    				f.Select("ref", 
    						f.Get(
        					f.MatchTerm(
            					f.Index(index), id))))

	if err != nil {
		panic(err)
	}

	return res
}

// LookupTaskObj: Look up a Task Object in FaunaDB using the given index and id.
func LookupTaskObj(fDB *FaunaDatabase, index string, id int) f.Value {

	res, err := fDB.Client.Query(
    					f.Get(
        					f.MatchTerm(
            					f.Index(index), id)))

	if err != nil {
		panic(err)
	}

	return res
}

// MatchParam: Match Parameters for FaunaDB
type MatchParam struct{
	Index string
	Term interface{}
}

// matchBuilder: Build matches based on the number of MatchParams received.
func matchBuilder(matchParams []MatchParam) f.Expr {
	
	var matches []f.Expr

	if(len(matchParams) > 0) {

		for _, matchParam := range matchParams {

			matches = append(matches, f.MatchTerm(f.Index(matchParam.Index), matchParam.Term))
		}

		if (len(matches) > 1) {

			return f.Intersection(matches)

		} else {

			return matches[0]
		}
	} else {

		return f.Match(f.Index("all_tasks"))
	}
}

