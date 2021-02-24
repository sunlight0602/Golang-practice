package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
	"regexp"
	"errors"
)

type Page struct {
	Title string
	Body []byte // []byte means "a byte slice"
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
	path := "data/"
	filename := p.Title + ".txt"

	return ioutil.WriteFile(path + filename, p.Body, 0600)
}
// This is a method named save that takes as its receiver p, a pointer to Page.
// It takes no parameters, and returns a value of type error.
// If all goes well, Page.save() will return nil

func loadPage(title string) (*Page, error) {
	path := "data/"
	filename := title + ".txt"

	body, error := ioutil.ReadFile(path + filename)
	if error != nil{
		return nil, error
	}

	return &Page{Title: title, Body: body}, error
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here we will extract the page title from the Request,
		// and call the provided handler 'fn'
		m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
			// http.NotFound(w, r)
			fmt.Println("enter makeHandler")
			http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
            return
        }
        fn(w, r, m[2])
	}
}
// The returned function is called a closure because it encloses values defined outside of it.
// In this case, the variable fn (the single argument to makeHandler) is enclosed by the closure.
// The variable fn will be one of our save, edit, or view handlers.

func viewHandler(w http.ResponseWriter, r *http.Request, title string)  {
	// title := r.URL.Path[len("/view/"):] // read title from r
	// // the path will invariably begin with "/view/"
	// // which is not part of the page's title.
	// title, err := getTitle(w, r)
    // if err != nil {
    //     return
	// }
	
	fmt.Println("enter viewHandler")
	p, err := loadPage(title)
	if err != nil{
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string){
	// title := r.URL.Path[len("/edit/"):]
    // title, err := getTitle(w, r)
    // if err != nil {
    //     return
    // }

    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
	}
	
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title := r.URL.Path[len("/save/"):]
    // title, err := getTitle(w, r)
    // if err != nil {
    //     return
    // }
	body := r.FormValue("body")
	
    p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func rootHandler(w http.ResponseWriter, r *http.Request, title string){
	fmt.Println("enter rootHandler")
	// http.Redirect(w, r, "/view/root", http.StatusFound)
}

func main()  {
	// // Create a page
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a Sample Page.")}
	// p1.save()

	// // Load the page
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))

	// http.HandleFunc("/view/", viewHandler)
    // http.HandleFunc("/edit/", editHandler)
    // http.HandleFunc("/save/", saveHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))
	
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/", makeHandler(rootHandler)) // root Handler is not executed
    log.Fatal(http.ListenAndServe(":8080", nil))
}