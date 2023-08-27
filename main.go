package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// set the active Avatar implementation
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

type TemplateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if cookieValue, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(cookieValue.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The address of the application.")
	flag.Parse() // parse the flag
	myApp := new(app)
	
	app_env, err := myApp.ReadEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	// set up gomniauth
	gomniauth.SetSecurityKey("mysecretopenhash")
	gomniauth.WithProviders(
		facebook.New(
			app_env["fbClientId"],
			app_env["fbClientSecret"],
			app_env["fbCallback"],
		),
		github.New(
			app_env["githubClientId"],
			app_env["githubClientSecret"],
			app_env["githubCallback"],
		),
		google.New(
			app_env["googleClientId"],
			app_env["googleClientSecret"],
			app_env["googleCallback"],
		),
	)

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&TemplateHandler{filename: "chat.html"}))
	http.Handle("/login", &TemplateHandler{filename: "login.html"})
	http.Handle("/upload", &TemplateHandler{filename: "upload.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/uploader", uploadHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))
	http.Handle("/room", r)

	// GET the room started
	go r.run()

	// START the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
