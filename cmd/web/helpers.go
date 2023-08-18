package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
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

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err 
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
	// If we try to use an invalid target destination, the Decode() method // will return an error with the type *form.InvalidDecoderError.We use // errors.As() to check for this and raise a panic rather than returning // the error.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) { 
			panic(err)
		}
		// For all other errors, we return them as normal.
		return err 
	}
	return nil
}