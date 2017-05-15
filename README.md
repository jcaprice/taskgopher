taskgopher
===

taskgopher is a task manager written in Go that uses FaunaDB for persistence.

## Get Taskgopher

`go get github.com/jcaprice/taskgopher`

## Setup

### FaunaDB

A Setup script will be created in the future. In the meantime, taskgopher requires the following classes and indexes (provided in Javascript for easy setup using the FaunaDB Query Console)

**Create your taskgopher Database**

Create a database named "taskgopher" under the root ("/") database using the Query Console in the FaunaDB dashboard.

```
q.CreateDatabase({ name: "taskgopher" })
```

**Create a server key for taskgopher**

Once the taskgopher database has been created, refresh the dashboard. Once refreshed, select the taskgopher database and create a server key for the newly created database. Save this server key for the configuration step below.

```
q.CreateKey({ database: q.Database("taskgopher"), role: "server" })
```

**Create Classes**

taskgopher requires two classes (for now), tasks and global. Global is used to store a counter that is incremented when tasks are added to assign an ID to tasks for easy lookup.

```
q.CreateClass({ name: "tasks" })
q.CreateClass({ name: "global" })
```

**Create Indexes**

taskgopher requires 4 indexes (with more in store) for effectively storing and retrieving tasks from FaunaDB.

The first index, all_tasks, indexes each task in the tasks class.

```
q.CreateIndex({ name: "all_tasks", source: q.Ref("classes/tasks") })
```

The second index, tasks_by_id, is used to retrieve task instances using an integer task ID that is incremented each time a task is added.

```
q.CreateIndex(
    {
      name: "tasks_by_id",
      source: q.Class("tasks"),
      unique: true,
      terms: [{ field: ["data", "id"] }]
    })
```
The index is unique as no two tasks should ever share the same integer ID.

The third index, tasks_by_status, is used to retrieve tasks based on their status (open or closed.)

```
q.CreateIndex(
    {
      name: "tasks_by_status",
      source: q.Class("tasks"),
      terms: [{ field: ["data", "status"] }]
    })
```

The fourth and final index, tasks_by_project, is used to retrieve tasks based on their project.

```
q.CreateIndex(
    {
      name: "tasks_by_project",
      source: q.Class("tasks"),
      terms: [{ field: ["data", "project"] }]
    })
```

**Create Counter Instance**

```
q.Create(q.Ref(q.Class("global"), "1"),{ data: { id: 0 } })
```

### Configuration

taskgopher uses viper (github.com/spf13/viper) for configuration, and the configuration file can be found in config/taskgopher.toml.

Only one configuration parameter needs to be changed at this time, authentication.secret. Replace the "SECRET" placeholder with your taskgopher server secret from the previous FaunaDB section.

## Installation

Install taskgopher by running `go install` from the taskgopher directory.

## Usage

taskgopher currently supports:

- Viewing Tasks
- Adding Tasks
- Listing Tasks
- Editing Tasks
- Opening and Closing Tasks
- Undoing Changes on Tasks

**Add Your First Task**

`> taskgopher add-task "This is your very first task" -project "taskgopher"`
```
+---------+------------------------------+
| Title   | This is your very first task |
+---------+------------------------------+
| Project | taskgopher                   |
+---------+------------------------------+
| Status  | open                         |
+---------+------------------------------+
```

**Add Your Second Task**

`> taskgopher add-task "This is your SECOND task" -project "taskgopher"`
```
+---------+--------------------------+
| Title   | This is your SECOND task |
+---------+--------------------------+
| Project | taskgopher               |
+---------+--------------------------+
| Status  | open                     |
+---------+--------------------------+
```

**List Your Tasks**

`> taskgopher list-tasks`
```
+----+------------------------------+------------+--------+
| ID |            TITLE             |  PROJECT   | STATUS |
+----+------------------------------+------------+--------+
|  1 | This is your very first task | taskgopher | open   |
|  2 | This is your SECOND task     | taskgopher | open   |
+----+------------------------------+------------+--------+
```

**Edit A Task**

`> taskgopher edit-task 1 -title "This is a new Title"`
```
+---------+---------------------+
| Title   | This is a new Title |
+---------+---------------------+
| Project | taskgopher          |
+---------+---------------------+
| Status  | open                |
+---------+---------------------+
```

`> taskgopher list-tasks`
```
+----+--------------------------+------------+--------+
| ID |          TITLE           |  PROJECT   | STATUS |
+----+--------------------------+------------+--------+
|  1 | This is a new Title      | taskgopher | open   |
|  2 | This is your SECOND task | taskgopher | open   |
+----+--------------------------+------------+--------+
```

**View a Task**

`> taskgopher 1`
```
+---------+---------------------+
| Title   | This is a new Title |
+---------+---------------------+
| Project | taskgopher          |
+---------+---------------------+
| Status  | open                |
+---------+---------------------+
```

**Undo Previous Edit**

`> taskgopher undo-task 1`
```
+---------+------------------------------+
| Title   | This is your very first task |
+---------+------------------------------+
| Project | taskgopher                   |
+---------+------------------------------+
| Status  | open                         |
+---------+------------------------------+
```

**Close a Task**

`> taskgopher close-task 1`

```
+---------+------------------------------+
| Title   | This is your very first task |
+---------+------------------------------+
| Project | taskgopher                   |
+---------+------------------------------+
| Status  | closed                       |
+---------+------------------------------+
```

**List Open Tasks**

`> taskgopher list-tasks -status "open"`

```
+----+--------------------------+------------+--------+
| ID |          TITLE           |  PROJECT   | STATUS |
+----+--------------------------+------------+--------+
|  2 | This is your SECOND task | taskgopher | open   |
+----+--------------------------+------------+--------+
```

**Re-Open a Task**

`> taskgopher open-task 1`

```
+---------+------------------------------+
| Title   | This is your very first task |
+---------+------------------------------+
| Project | taskgopher                   |
+---------+------------------------------+
| Status  | open                         |
+---------+------------------------------+
```
