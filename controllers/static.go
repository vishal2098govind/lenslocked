package controllers

import (
	"html/template"
	"net/http"
)

func FAQ(t Template) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		questions := []struct {
			Question string
			Answer   template.HTML
		}{
			{Question: "Question q1", Answer: "Answer a1"},
			{Question: "Question q2", Answer: "Answer a2"},
			{Question: "Question q3", Answer: "Answer a3"},
			{Question: "Question q4", Answer: "Answer a4"},
			{
				Question: "Question q5",
				Answer:   `Answer a5 <a href="mailto:vishal.govind@gmail.com">vishal.govind@gmail.com</a>`,
			},
		}
		t.Execute(w, r, questions)
	}

}

func StaticHandler(t Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, r, nil)
	}
}
