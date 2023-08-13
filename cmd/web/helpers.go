package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError helper writes error message and stack trace to errorLog
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// sends a specific status code and description to user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// 404 not found
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// initialize buffer
	buf := new(bytes.Buffer)

	// Write the template to buffer instead of straight to http writer, this is done so that if there's a run time error, the page won't be
	// partially loaded with incorrect response code if there's a runtime error
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	// Write out the provided HTTP status code ('200 OK', '400 Bad Request' // etc) if there's no errors when writing to buffer.
	w.WriteHeader(status)

	// Write contents of buffer to http.ResponseWriter
	buf.WriteTo(w)

}
