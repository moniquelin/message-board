package main

import (
	"errors"
	"fmt"
	"message-board/internal/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	messages, err := app.messages.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Messages = messages
	// Pass the data to the render() helper as normal.
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) messageView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	message, err := app.messages.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// And do the same thing again here...
	data := app.newTemplateData(r)
	data.Message = message
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) messageCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// Use the http.Error() function to send a 405 status code and "Method Not
		// Allowed" string as the response body.
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// Pass the data to the MessageModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.messages.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the message.
	http.Redirect(w, r, fmt.Sprintf("/message/view?id=%d", id), http.StatusSeeOther)
}
