package handler

import (
	"encoding/json"
	"net/http"
	"simple_todo/model"
	"simple_todo/repo"
	"strconv"
)

type Handler struct {
	repo *repo.Repo
}

func New(repo *repo.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.repo.GetAll(r.Context())
	if err != nil {
		writeJson(w, model.ErrorResponseAPI{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	writeJson(w, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJson(w, model.ErrorResponseAPI{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	result, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeJson(w, model.ErrorResponseAPI{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	writeJson(w, result)
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	var p model.TodoAddPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJson(w, model.ErrorResponseAPI{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	if err := h.repo.Add(r.Context(), p.Task); err != nil {
		writeJson(w, model.ErrorResponseAPI{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJson(w, model.ErrorResponseAPI{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeJson(w, model.ErrorResponseAPI{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJson(w, model.ErrorResponseAPI{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	var p model.TodoUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJson(w, model.ErrorResponseAPI{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	if p.Done != nil {
		if err := h.repo.EditStatus(r.Context(), id, *p.Done); err != nil {
			writeJson(w, model.ErrorResponseAPI{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			})
		}
	}

	if p.Task != nil {
		if err := h.repo.EditTask(r.Context(), id, *p.Task); err != nil {
			writeJson(w, model.ErrorResponseAPI{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			})
		}
	}
}

func writeJson(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
