package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type Timeline interface {
	GetDeputyEvents(sirius.Context, int) (sirius.DeputyEvents, error)
}

type timelineVars struct {
	DeputyEvents sirius.DeputyEvents
	AppVars
}

type TimelineHandler struct {
	router
}

func (h *TimelineHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	deputyEvents, err := h.Client().GetDeputyEvents(ctx, v.DeputyId())
	if err != nil {
		return err
	}

	v.PageName = "Timeline"

	vars := timelineVars{
		DeputyEvents: deputyEvents,
		AppVars:      v,
	}

	return h.execute(w, r, vars, v)
}
