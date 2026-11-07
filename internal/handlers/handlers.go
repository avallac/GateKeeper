package handlers

import (
	"html/template"
	"net/http"
	"regexp"
	"slices"

	"github.com/avallac/gatekeeper/internal/service/ipguard"
)

type Handlers struct {
	App  *ipguard.Client
	Tmpl *template.Template
}

type LoginPageData struct {
	Error string
	Name  string
	Phone string
}

type HomePageData struct {
	Sites []ipguard.AppSite
	Log   []ipguard.AppLog
}

type ValidateCodePageData struct {
	Error string
	Code  string
	Phone string
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	var err error
	if h.App.IsConnected() {
		log := slices.Clone(h.App.Log)
		slices.Reverse(log)
		data := HomePageData{
			Sites: h.App.Sites,
			Log:   log,
		}
		err = h.Tmpl.ExecuteTemplate(w, "index.html", data)
	} else if h.App.IsNeededPhone() {
		data := LoginPageData{
			Name:  h.App.UserConfig.Name,
			Phone: h.App.UserConfig.Phone,
		}
		err = h.Tmpl.ExecuteTemplate(w, "login.html", data)
	} else if h.App.IsNeededCode() {
		data := LoginPageData{
			Phone: h.App.UserConfig.Phone,
		}
		err = h.Tmpl.ExecuteTemplate(w, "code.html", data)
	} else if h.App.IsWebWaiting() {
		err = h.Tmpl.ExecuteTemplate(w, "waiting.html", nil)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func ingressBase(r *http.Request) string {
	return r.Header.Get("X-Ingress-Path")
}

func (h *Handlers) checkDoor(id string) bool {
	for _, site := range h.App.Sites {
		for _, board := range site.Boards {
			if board.Identity == id {
				return true
			}
		}
	}
	return false
}

func (h *Handlers) OpenDoor(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	if !h.checkDoor(id) {
		http.Error(w, "Bad door", http.StatusBadRequest)
		return
	}
	h.App.OpenDoor(id)
	http.Redirect(w, r, ingressBase(r)+"/", http.StatusFound)
}

func (h *Handlers) Check(w http.ResponseWriter, r *http.Request) {
	h.App.CheckUser()
}

func (h *Handlers) Override(w http.ResponseWriter, r *http.Request) {
	h.App.Override()
}

func (h *Handlers) SendCode(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	phone := r.FormValue("phone")
	data := ValidateCodePageData{
		Code:  code,
		Phone: phone,
	}
	var codeRe = regexp.MustCompile(`^\d{6}$`)
	if !codeRe.MatchString(code) {
		data.Error = "Invalid code"
	}

	var phoneRe = regexp.MustCompile(`^\+?[0-9]{7,30}$`)
	if !phoneRe.MatchString(phone) {
		data.Error = "Invalid phone"
	}
	h.App.ValidateCode(code)
	http.Redirect(w, r, ingressBase(r)+"/", http.StatusFound)
}

func (h *Handlers) AskCode(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	data := LoginPageData{
		Name:  name,
		Phone: phone,
	}

	var nameRe = regexp.MustCompile(`^[a-zA-Z ]+$`)
	if !nameRe.MatchString(name) {
		data.Error = "Invalid name"
	}

	var phoneRe = regexp.MustCompile(`^\+?[0-9]{7,30}$`)
	if !phoneRe.MatchString(phone) {
		data.Error = "Invalid phone"
	}

	if data.Error != "" {
		err = h.Tmpl.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	h.App.UpdateUserInfo(name, phone)
	h.App.AskCode()
	h.Tmpl.ExecuteTemplate(w, "code.html", data)
}
