// A CLI tool for prioritizing tasks in Things 3.
package main

func main() {
	Logger.Info("Starting prioritizer-terminal")
	todos, err := getThingsTodos()
	if err != nil {
		Logger.Fatal(err)
	}
	Logger.Debugf("Got %d todos", len(todos))
	Logger.Debugf("Todos: %+v", todos)
}
