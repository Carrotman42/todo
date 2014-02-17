package main

import (
	"fmt"
	"net/http"
	"io"
	"strconv"
	"os"
	"encoding/gob"
)

func write(w io.Writer, s...interface{}) {
	fmt.Fprint(w, s...)
}

const saveDoc = "saveDoc"

func main() {
	fmt.Println("This is your to-do list")
	todos = load(saveDoc)
	http.HandleFunc("/", MainPage)
	http.HandleFunc("/new", PostNew)
	http.HandleFunc("/markDone", markDone)
	http.HandleFunc("/debug", TestPost)
	http.HandleFunc("/resources/", HandleResource)
	http.ListenAndServe(":16005", nil)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	write(w, CALENDAR_REQ)
	write(w, NEW_TODO_HTML)
	write(w, POST_FORM)
	
	for i,v := range todos {
		if v.Done == true {
			write(w,
				`<input checked type=checkbox onclick="submitDone(`,
				i,
				`)"><del>`,
				v.Name,
				"</del></input>",
				"<br>\n")
		} else {
			write(w,
				`<input type=checkbox onclick="submitDone(`,
				i,
				`)">`,
				v.Name,
				"</input>",
				"<br>\n")
		}
	}
}

const CALENDAR_REQ = `
<head>
<link type="text/css" rel="stylesheet" href="resources/dhtmlgoodies_calendar/dhtmlgoodies_calendar.css" media="screen"></LINK>
<SCRIPT type="text/javascript" src="resources/dhtmlgoodies_calendar/dhtmlgoodies_calendar.js"></script>
</head>
`

// These two must match
const JS_TIME_FORMAT = `yyyy-mm-dd hh:ii`
const GO_TIME_FORMAT = `2006-01-02 15:04`

const NEW_TODO_HTML = `
<form method="post" action="new">
New Task Name: <input name="name" />
<input type=submit value="Add" />
<input type="text" readonly name="dueDate" />
<input type="button" value="Cal" onclick="displayCalendar(dueDate,'` + JS_TIME_FORMAT + `',this,true)" />
</form>
`

const POST_FORM = `
<form id="doneform" action="markDone" method="POST">
<input type="hidden" name="data" value="default value">
</form>
<script>
function submitDone(index) {
	var theForm = document.getElementById('doneform')
	theForm.data.value = index
	theForm.submit()
}
</script>
`

func HandleResource(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Path[1:]
	http.ServeFile(w, r, resource)
}

func markDone (w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	index := r.Form["data"][0]
	i, _ := strconv.Atoi(index)
	todos[i].Done = !todos[i].Done
	
	http.Redirect(w, r, "/", 303)
	save(todos, saveDoc)
}

func TestPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	write(w, r.Form)
}

func PostNew(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	
	name := r.Form["name"][0]
	//TODO: due := r.Form["dueDate"][0]
	
	newTodo := Todo{
		Name: name,
		Done: false,
	}
	todos = append(todos, newTodo)
	http.Redirect(w, r, "/", 303)
	save(todos, saveDoc)
}

var todos []Todo

type Todo struct {
	Name string
	Done bool
}

func save (x []Todo, saveTo string) {
	saveFile, err := os.Create(saveTo)
	if err != nil {
		panic(err)
	}
	defer saveFile.Close()
	
	enc := gob.NewEncoder(saveFile)
	err = enc.Encode(x)
	if err != nil {
		panic(err)
	}
}

func load(toLoad string) []Todo {
	loadFile, err := os.Open(toLoad)
	if err != nil {
		return nil
	}
	defer loadFile.Close()
	
	dec := gob.NewDecoder(loadFile)
	var ret []Todo
	err = dec.Decode(&ret)
	if err != nil {
		panic(err)
	}
	return ret
}