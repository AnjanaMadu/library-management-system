package html

import (
	"os"
	"strings"
)

func MakeTable(title string, columns []string, rows [][]string) string {
	f, _ := os.ReadFile("./public/table.html")
	s := string(f)

	// Replace title
	s = strings.ReplaceAll(s, "{{title}}", title)

	// Make columns
	tr := createElement("tr")
	for i, c := range columns {
		th := createElement("th")
		th.setAttribute("scope", "col")
		th.addClass("p-3")
		if i == 0 {
			th.addClass("pl-5")
		}
		th.appendHTML(c)
		tr.appendHTML(th.getString())
	}
	s = strings.ReplaceAll(s, "{{columns}}", tr.getString())

	// Make rows
	trs := ""
	for ir, r := range rows {
		tr := createElement("tr")
		for i, c := range r {
			td := createElement("td")
			td.addClass("p-2 text-slate-400")
			if i == 0 {
				td.addClass("pl-8")
			} else {
				td.addClass("pl-3")
			}
			if ir != len(rows)-1 {
				td.addClass("border-b border-slate-700 ")
			}
			td.appendHTML(c)
			tr.appendHTML(td.getString())
		}
		trs += tr.getString()
	}
	s = strings.ReplaceAll(s, "{{rows}}", trs)
	return s
}
