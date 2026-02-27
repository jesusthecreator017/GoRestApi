package main

import (
	"net/http"

	"github.com/jesusthecreator017/fswithgo/cmd/api/helpers"
)

func (app *application) adminStatsHandler(w http.ResponseWriter, req *http.Request) {
	stats, err := app.store.Admin.GetStats(req.Context())
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get admin stats")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"stats": stats})
}
