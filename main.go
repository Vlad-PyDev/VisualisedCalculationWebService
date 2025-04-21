package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/Knetic/govaluate"
)

type PageData struct {
	Title string
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "IP address to listen on")
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{Title: "Добро пожаловать"}
		renderTemplate(w, "start.html", data)
	})

	http.HandleFunc("/calculator", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{Title: "Калькулятор"}
		renderTemplate(w, "calculator.html", data)
	})

	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			expression := r.FormValue("expression")
			result := evaluate(expression)
			fmt.Fprintf(w, result)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	address := fmt.Sprintf("%s:%s", *ip, *port)
	fmt.Printf("Сервер запущен на http://%s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func evaluate(expression string) string {
	expression = sanitizeInput(expression)
	if expression == "" {
		return "Ошибка"
	}

	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "Ошибка"
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		return "Ошибка"
	}

	return fmt.Sprintf("%v", result)
}

func sanitizeInput(input string) string {
	allowedChars := "0123456789+-*/(). "
	cleaned := ""
	for _, char := range input {
		if strings.ContainsRune(allowedChars, char) {
			cleaned += string(char)
		}
	}
	return cleaned
}
