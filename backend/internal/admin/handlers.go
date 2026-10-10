// Package admin is the HTML management UI at /admin: login, settings editor,
// generic CRUD for every schema.Resource, the message inbox and the media
// library. Templates are embedded from web/.
package admin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"portfolio/backend/internal/media"
	"portfolio/backend/internal/models"
	"portfolio/backend/internal/schema"
	"portfolio/backend/internal/store"
)

// Handler holds the admin's dependencies.
type Handler struct {
	store    *store.Store
	sessions *Sessions
	media    *media.Library
	log      *slog.Logger
	tmpl     *renderer
}

// New constructs the admin handler.
func New(st *store.Store, sessions *Sessions, lib *media.Library, log *slog.Logger) (*Handler, error) {
	r, err := newRenderer(log)
	if err != nil {
		return nil, err
	}
	return &Handler{store: st, sessions: sessions, media: lib, log: log, tmpl: r}, nil
}

// Routes registers every admin endpoint on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /admin/assets/{file}", staticHandler())
	mux.HandleFunc("GET /admin/login", h.loginForm)
	mux.HandleFunc("POST /admin/login", h.login)
	mux.HandleFunc("POST /admin/logout", h.auth(h.logout))

	mux.HandleFunc("GET /admin", h.auth(h.dashboard))
	mux.HandleFunc("GET /admin/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /admin/settings", h.auth(h.settingsForm))
	mux.HandleFunc("POST /admin/settings", h.auth(h.settingsSave))
	mux.HandleFunc("GET /admin/password", h.auth(h.passwordForm))
	mux.HandleFunc("POST /admin/password", h.auth(h.passwordSave))

	mux.HandleFunc("GET /admin/messages", h.auth(h.messages))
	mux.HandleFunc("GET /admin/messages/{id}", h.auth(h.message))
	mux.HandleFunc("POST /admin/messages/{id}/read", h.auth(h.messageToggleRead))
	mux.HandleFunc("POST /admin/messages/{id}/delete", h.auth(h.messageDelete))

	mux.HandleFunc("GET /admin/media", h.auth(h.mediaList))
	mux.HandleFunc("POST /admin/media/upload", h.auth(h.mediaUpload))
	mux.HandleFunc("POST /admin/media/delete", h.auth(h.mediaDelete))

	// One explicit route set per resource. Literal segments keep the patterns
	// unambiguous for ServeMux (a `{res}` wildcard would clash with
	// /admin/messages/{id} and /admin/assets/{file}).
	for i := range schema.Resources {
		res := &schema.Resources[i]
		p := "/admin/" + res.Key
		mux.HandleFunc("GET "+p, h.auth(h.resourceList(res)))
		mux.HandleFunc("GET "+p+"/new", h.auth(h.resourceForm(res)))
		mux.HandleFunc("POST "+p+"/new", h.auth(h.resourceSave(res)))
		mux.HandleFunc("GET "+p+"/{id}", h.auth(h.resourceForm(res)))
		mux.HandleFunc("POST "+p+"/{id}", h.auth(h.resourceSave(res)))
		mux.HandleFunc("POST "+p+"/{id}/delete", h.auth(h.resourceDelete(res)))
		mux.HandleFunc("POST "+p+"/{id}/move", h.auth(h.resourceMove(res)))
	}
}

// --- auth & shared plumbing ------------------------------------------------

type ctxKey int

const userKey ctxKey = 1

// auth requires a valid session and, on POST, a matching CSRF token.
func (h *Handler) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, cookie, ok := h.sessions.current(r)
		if !ok {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		if r.Method == http.MethodPost {
			// Multipart (uploads) is parsed by the handler; everything else here.
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				if err := r.ParseForm(); err != nil {
					http.Error(w, "bad form", http.StatusBadRequest)
					return
				}
			} else if err := r.ParseMultipartForm(media.MaxUploadBytes); err != nil {
				http.Error(w, "upload too large or malformed", http.StatusBadRequest)
				return
			}
			if r.FormValue("_csrf") != h.sessions.CSRF(cookie) {
				http.Error(w, "invalid CSRF token; reload the page and try again", http.StatusForbidden)
				return
			}
		}
		ctx := context.WithValue(r.Context(), userKey, user)
		next(w, r.WithContext(ctx))
	}
}

// page builds the common view for the current request.
func (h *Handler) page(w http.ResponseWriter, r *http.Request, title, active string, data any) view {
	v := view{Title: title, Active: active, Resources: schema.Resources, Data: data}
	if u, ok := r.Context().Value(userKey).(string); ok {
		v.User = u
	}
	if _, cookie, ok := h.sessions.current(r); ok {
		v.CSRF = h.sessions.CSRF(cookie)
	}
	v.FlashKind, v.Flash = h.sessions.takeFlash(w, r)
	if v.User != "" {
		v.Unread, _ = h.store.UnreadMessages(r.Context())
		if s, err := h.store.Settings(r.Context()); err == nil {
			v.Handle = s["handle"]
		}
	}
	return v
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request, to, kind, msg string) {
	if msg != "" {
		h.sessions.setFlash(w, kind, msg)
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
}

func (h *Handler) serverError(w http.ResponseWriter, r *http.Request, what string, err error) {
	h.log.Error("admin", "what", what, "err", err, "path", r.URL.Path)
	http.Error(w, "Something went wrong: "+what, http.StatusInternalServerError)
}

func pathID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id, err == nil && id > 0
}

// --- login -----------------------------------------------------------------

func (h *Handler) loginForm(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := h.sessions.current(r); ok {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	h.tmpl.render(w, http.StatusOK, "login", h.page(w, r, "Log in", "", nil))
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	user, err := h.store.UserByUsername(r.Context(), username)
	// Always run bcrypt so a missing user takes as long as a wrong password.
	hash := user.PasswordHash
	if err != nil {
		hash = "$2a$10$7EqJtq98hPqEX7fNZaFWoOhi5XcG8yXcq6o9F8T9jE6RrQ0q0qK7K"
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil || err != nil {
		time.Sleep(300 * time.Millisecond)
		v := h.page(w, r, "Log in", "", map[string]string{"Username": username})
		v.Flash, v.FlashKind = "Wrong username or password.", "error"
		h.tmpl.render(w, http.StatusUnauthorized, "login", v)
		return
	}
	h.sessions.Set(w, user.Username)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	h.sessions.Clear(w)
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// --- dashboard -------------------------------------------------------------

type dashboardData struct {
	Counts   []resourceCount
	Messages []messageRow
	Unread   int
}

type resourceCount struct {
	Res   schema.Resource
	Count int
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	var d dashboardData
	for i := range schema.Resources {
		n, err := h.store.CountRecords(r.Context(), &schema.Resources[i])
		if err != nil {
			h.serverError(w, r, "counting content", err)
			return
		}
		d.Counts = append(d.Counts, resourceCount{Res: schema.Resources[i], Count: n})
	}
	msgs, err := h.store.Messages(r.Context(), 5)
	if err != nil {
		h.serverError(w, r, "loading messages", err)
		return
	}
	d.Messages = messageRows(msgs)
	d.Unread, _ = h.store.UnreadMessages(r.Context())
	h.tmpl.render(w, http.StatusOK, "dashboard", h.page(w, r, "Dashboard", "dashboard", d))
}

// --- settings --------------------------------------------------------------

type settingsData struct {
	Groups []groupView
	Media  []string
}

type groupView struct {
	G      schema.Group
	Fields []fieldView
}

func (h *Handler) settingsView(values map[string]string, errs map[string]string) settingsData {
	var d settingsData
	for _, g := range schema.SettingsGroups {
		gv := groupView{G: g}
		for _, f := range g.Fields {
			gv.Fields = append(gv.Fields, fieldView{F: f, Value: values[f.Key], Error: errs[f.Key], Checked: map[string]bool{}})
		}
		d.Groups = append(d.Groups, gv)
	}
	if files, err := h.media.List(); err == nil {
		d.Media = mediaOptions(files)
	}
	return d
}

func (h *Handler) settingsForm(w http.ResponseWriter, r *http.Request) {
	values, err := h.store.Settings(r.Context())
	if err != nil {
		h.serverError(w, r, "loading settings", err)
		return
	}
	h.tmpl.render(w, http.StatusOK, "settings", h.page(w, r, "Settings", "settings", h.settingsView(values, nil)))
}

func (h *Handler) settingsSave(w http.ResponseWriter, r *http.Request) {
	values := map[string]string{}
	errs := map[string]string{}
	for _, f := range schema.SettingFields() {
		raw := formValue(r, f)
		v, problem := f.Clean(raw)
		if problem != "" {
			errs[f.Key] = problem
			values[f.Key] = raw
			continue
		}
		values[f.Key] = v
	}
	if len(errs) > 0 {
		v := h.page(w, r, "Settings", "settings", h.settingsView(values, errs))
		v.Flash, v.FlashKind = "Please fix the highlighted fields.", "error"
		h.tmpl.render(w, http.StatusUnprocessableEntity, "settings", v)
		return
	}
	if err := h.store.SaveSettings(r.Context(), values); err != nil {
		h.serverError(w, r, "saving settings", err)
		return
	}
	h.redirect(w, r, "/admin/settings", "ok", "Settings saved.")
}

// formValue reads a field from the posted form, joining multi-valued inputs
// (Tags checkboxes) with commas.
func formValue(r *http.Request, f schema.Field) string {
	if f.Kind == schema.Tags {
		return strings.Join(r.Form[f.Key], ",")
	}
	return r.FormValue(f.Key)
}

// --- password --------------------------------------------------------------

func (h *Handler) passwordForm(w http.ResponseWriter, r *http.Request) {
	h.tmpl.render(w, http.StatusOK, "password", h.page(w, r, "Change password", "password", nil))
}

func (h *Handler) passwordSave(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(string)
	current, next, confirm := r.FormValue("current"), r.FormValue("new"), r.FormValue("confirm")

	u, err := h.store.UserByUsername(r.Context(), user)
	if err != nil {
		h.serverError(w, r, "loading account", err)
		return
	}
	fail := func(msg string) {
		v := h.page(w, r, "Change password", "password", nil)
		v.Flash, v.FlashKind = msg, "error"
		h.tmpl.render(w, http.StatusUnprocessableEntity, "password", v)
	}
	switch {
	case bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(current)) != nil:
		fail("Current password is wrong.")
		return
	case len(next) < 10:
		fail("New password must be at least 10 characters.")
		return
	case next != confirm:
		fail("New passwords do not match.")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(next), bcrypt.DefaultCost)
	if err != nil {
		h.serverError(w, r, "hashing password", err)
		return
	}
	if err := h.store.UpdatePassword(r.Context(), user, string(hash)); err != nil {
		h.serverError(w, r, "saving password", err)
		return
	}
	h.redirect(w, r, "/admin", "ok", "Password changed.")
}

// --- generic resources -----------------------------------------------------

type listData struct {
	Res     schema.Resource
	Columns []schema.Field
	Rows    []listRow
}

type listRow struct {
	ID    int64
	Cells []string
	First bool
	Last  bool
}

// The resource handlers are built per resource in Routes, so each closure
// knows its schema.Resource and only the row id comes from the URL.

func (h *Handler) resourceList(res *schema.Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { h.listResource(res, w, r) }
}

func (h *Handler) listResource(res *schema.Resource, w http.ResponseWriter, r *http.Request) {
	records, err := h.store.ListRecords(r.Context(), res)
	if err != nil {
		h.serverError(w, r, "loading "+res.Plural, err)
		return
	}
	d := listData{Res: *res}
	for _, f := range res.Fields {
		if f.List {
			d.Columns = append(d.Columns, f)
		}
	}
	for i, rec := range records {
		row := listRow{ID: rec.ID, First: i == 0, Last: i == len(records)-1}
		for _, f := range d.Columns {
			row.Cells = append(row.Cells, displayValue(f, rec.Values[f.Key]))
		}
		d.Rows = append(d.Rows, row)
	}
	h.tmpl.render(w, http.StatusOK, "list", h.page(w, r, res.Plural, res.Key, d))
}

// displayValue formats a stored value for the list table.
func displayValue(f schema.Field, v string) string {
	switch f.Kind {
	case schema.Bool:
		if v == "1" {
			return "Yes"
		}
		return "—"
	case schema.Lines:
		return strings.Join(schema.SplitLines(v), ", ")
	case schema.Tags:
		return strings.Join(schema.SplitComma(v), ", ")
	}
	return v
}

type formData struct {
	Res    schema.Resource
	ID     int64
	Fields []fieldView
	Media  []string
}

func (h *Handler) formView(res *schema.Resource, id int64, values, errs map[string]string) formData {
	d := formData{Res: *res, ID: id}
	for _, f := range res.Fields {
		fv := fieldView{F: f, Value: values[f.Key], Error: errs[f.Key], Checked: map[string]bool{}}
		if f.Kind == schema.Tags {
			for _, t := range schema.SplitComma(fv.Value) {
				fv.Checked[t] = true
			}
		}
		d.Fields = append(d.Fields, fv)
	}
	if files, err := h.media.List(); err == nil {
		d.Media = mediaOptions(files)
	}
	return d
}

func (h *Handler) resourceForm(res *schema.Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { h.formResource(res, w, r) }
}

func (h *Handler) formResource(res *schema.Resource, w http.ResponseWriter, r *http.Request) {
	values := map[string]string{}
	var id int64
	if r.PathValue("id") != "" {
		var ok bool
		if id, ok = pathID(r); !ok {
			http.NotFound(w, r)
			return
		}
		rec, err := h.store.GetRecord(r.Context(), res, id)
		if errors.Is(err, store.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			h.serverError(w, r, "loading "+res.Singular, err)
			return
		}
		values = rec.Values
	} else {
		// Sensible defaults for a new row.
		if f, ok := res.Field("status"); ok && len(f.Options) > 0 {
			values["status"] = f.Options[0]
		}
		if _, ok := res.Field("featured"); ok {
			values["featured"] = "1"
		}
	}
	title := "New " + strings.ToLower(res.Singular)
	if id != 0 {
		title = "Edit " + strings.ToLower(res.Singular)
	}
	h.tmpl.render(w, http.StatusOK, "form", h.page(w, r, title, res.Key, h.formView(res, id, values, nil)))
}

func (h *Handler) resourceSave(res *schema.Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { h.saveResource(res, w, r) }
}

func (h *Handler) saveResource(res *schema.Resource, w http.ResponseWriter, r *http.Request) {
	var id int64
	if r.PathValue("id") != "" {
		var ok bool
		if id, ok = pathID(r); !ok {
			http.NotFound(w, r)
			return
		}
	}

	values := map[string]string{}
	errs := map[string]string{}
	for _, f := range res.Fields {
		raw := formValue(r, f)
		v, problem := f.Clean(raw)
		if problem != "" {
			errs[f.Key] = problem
			values[f.Key] = raw
			continue
		}
		values[f.Key] = v
	}

	var err error
	if len(errs) == 0 {
		if id == 0 {
			_, err = h.store.CreateRecord(r.Context(), res, values)
		} else {
			err = h.store.UpdateRecord(r.Context(), res, id, values)
		}
		switch {
		case errors.Is(err, store.ErrDuplicate):
			key := "slug"
			if _, ok := res.Field(key); !ok {
				key = res.Fields[0].Key
			}
			errs[key] = "This value is already used by another entry."
		case errors.Is(err, store.ErrNotFound):
			http.NotFound(w, r)
			return
		case err != nil:
			h.serverError(w, r, "saving "+res.Singular, err)
			return
		}
	}
	if len(errs) > 0 {
		title := "New " + strings.ToLower(res.Singular)
		if id != 0 {
			title = "Edit " + strings.ToLower(res.Singular)
		}
		v := h.page(w, r, title, res.Key, h.formView(res, id, values, errs))
		v.Flash, v.FlashKind = "Please fix the highlighted fields.", "error"
		h.tmpl.render(w, http.StatusUnprocessableEntity, "form", v)
		return
	}
	h.redirect(w, r, "/admin/"+res.Key, "ok", res.Singular+" saved.")
}

func (h *Handler) resourceDelete(res *schema.Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { h.deleteResource(res, w, r) }
}

func (h *Handler) deleteResource(res *schema.Resource, w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if err := h.store.DeleteRecord(r.Context(), res, id); err != nil && !errors.Is(err, store.ErrNotFound) {
		h.serverError(w, r, "deleting "+res.Singular, err)
		return
	}
	h.redirect(w, r, "/admin/"+res.Key, "ok", res.Singular+" deleted.")
}

func (h *Handler) resourceMove(res *schema.Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { h.moveResource(res, w, r) }
}

func (h *Handler) moveResource(res *schema.Resource, w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	dir := r.FormValue("dir")
	if dir != "up" && dir != "down" {
		http.Error(w, "dir must be up or down", http.StatusBadRequest)
		return
	}
	if err := h.store.MoveRecord(r.Context(), res, id, dir); err != nil && !errors.Is(err, store.ErrNotFound) {
		h.serverError(w, r, "reordering "+res.Plural, err)
		return
	}
	h.redirect(w, r, "/admin/"+res.Key, "", "")
}

// --- messages --------------------------------------------------------------

type messageRow struct {
	ID      int64
	Name    string
	Email   string
	Body    string
	Read    bool
	Created string
}

func messageRows(msgs []models.Message) []messageRow {
	out := make([]messageRow, len(msgs))
	for i, m := range msgs {
		created := m.CreatedAt
		if t, err := time.Parse(time.RFC3339, m.CreatedAt); err == nil {
			created = t.Local().Format("2 Jan 2006, 15:04")
		}
		out[i] = messageRow{ID: m.ID, Name: m.Name, Email: m.Email, Body: m.Body, Read: m.Read, Created: created}
	}
	return out
}

func (h *Handler) messages(w http.ResponseWriter, r *http.Request) {
	msgs, err := h.store.Messages(r.Context(), 200)
	if err != nil {
		h.serverError(w, r, "loading messages", err)
		return
	}
	h.tmpl.render(w, http.StatusOK, "messages", h.page(w, r, "Messages", "messages", messageRows(msgs)))
}

func (h *Handler) message(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	m, err := h.store.Message(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		h.serverError(w, r, "loading message", err)
		return
	}
	if !m.Read {
		_ = h.store.MarkMessageRead(r.Context(), id, true)
		m.Read = true
	}
	h.tmpl.render(w, http.StatusOK, "message", h.page(w, r, "Message from "+m.Name, "messages", messageRows([]models.Message{m})[0]))
}

func (h *Handler) messageToggleRead(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	read := r.FormValue("read") == "1"
	if err := h.store.MarkMessageRead(r.Context(), id, read); err != nil {
		h.serverError(w, r, "updating message", err)
		return
	}
	h.redirect(w, r, "/admin/messages", "", "")
}

func (h *Handler) messageDelete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if err := h.store.DeleteMessage(r.Context(), id); err != nil && !errors.Is(err, store.ErrNotFound) {
		h.serverError(w, r, "deleting message", err)
		return
	}
	h.redirect(w, r, "/admin/messages", "ok", "Message deleted.")
}

// --- media -----------------------------------------------------------------

func (h *Handler) mediaList(w http.ResponseWriter, r *http.Request) {
	files, err := h.media.List()
	if err != nil {
		h.serverError(w, r, "listing uploads", err)
		return
	}
	h.tmpl.render(w, http.StatusOK, "media", h.page(w, r, "Media", "media", files))
}

func (h *Handler) mediaUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		h.redirect(w, r, "/admin/media", "error", "Choose a file to upload.")
		return
	}
	defer file.Close()
	saved, err := h.media.Save(header.Filename, file)
	if err != nil {
		h.redirect(w, r, "/admin/media", "error", "Upload failed: "+err.Error())
		return
	}
	h.redirect(w, r, "/admin/media", "ok", fmt.Sprintf("Uploaded %s — use it as %s", saved.Name, saved.URL))
}

func (h *Handler) mediaDelete(w http.ResponseWriter, r *http.Request) {
	if err := h.media.Delete(r.FormValue("name")); err != nil {
		h.redirect(w, r, "/admin/media", "error", "Could not delete: "+err.Error())
		return
	}
	h.redirect(w, r, "/admin/media", "ok", "File deleted.")
}
