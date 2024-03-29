package session

import (
	"math"
	"path"
	"sort"
	"strings"

	"github.com/altid/libs/markup"
	"github.com/altid/libs/service/controller"
	"rsc.io/pdf"
)

func parsePdfBody(c controller.Controller, docname string, r *pdf.Reader) error {
	numPages := r.NumPage()

	w, err := c.MainWriter(docname)
	if err != nil {
		return err
	}

	body := markup.NewCleaner(w)
	defer body.Close()
	for i := 1; i <= numPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		parsePage(body, page)
	}

	return nil
}

func parsePdfTitle(c controller.Controller, docname string, r *pdf.Reader) error {
	w, err := c.TitleWriter(docname)
	if err != nil {
		return err
	}

	title := markup.NewCleaner(w)
	defer title.Close()

	tstring := r.Outline().Title
	if tstring == "" {
		tstring = docname
	}
	title.WriteStringEscaped(tstring)

	return nil
}

func parsePdfSidebar(c controller.Controller, docname string, r *pdf.Reader) {
	entries := make(chan entry)
	defer close(entries)
	go writeOutline(c, docname, entries)
	// Skip the document title - do the first walk here
	for _, item := range r.Outline().Child {
		walkPdfOutline(item, 0, entries)
	}
}

func walkPdfOutline(r pdf.Outline, n int, entries chan entry) {
	if r.Title != "" {
		entries <- entry{
			len: n,
			msg: []byte(r.Title),
			url: []byte(r.Title),
		}
	}
	n++
	for _, item := range r.Child {
		walkPdfOutline(item, n, entries)
	}
}

func parsePdf(c controller.Controller, newfile string) error {
	docname := path.Base(newfile)
	if ! c.HasBuffer(docname) {
		c.CreateBuffer(docname)
	}

	w, err := c.StatusWriter(docname)
	if err != nil {
		return err
	}

	status := markup.NewCleaner(w)
	defer status.Close()

	pages, err := pdf.Open(newfile)
	if err != nil {
		return err
	}

	status.WriteString("Parsing file...")
	parsePdfTitle(c, docname, pages)
	parsePdfSidebar(c, docname, pages)

	if e := parsePdfBody(c, docname, pages); e != nil {
		return e
	}

	return c.Remove(docname, "status")
}

func parsePage(body *markup.Cleaner, page pdf.Page) {
	content := page.Content()

	var text []pdf.Text
	text = append(text, content.Text...)
	text = findWords(text)
	for _, t := range text {
		body.WriteStringEscaped(t.S + " ")
	}
}

// Adapted from golang.org/x/arch/arm/armspec
// Copyright 2014 The Go Authors
func findWords(chars []pdf.Text) (words []pdf.Text) {
	const nudge = 1
	sort.Sort(pdf.TextVertical(chars))
	old := -100000.0
	for i, c := range chars {
		if c.Y != old && math.Abs(old-c.Y) < nudge {
			chars[i].Y = old
		} else {
			old = c.Y
		}
	}
	sort.Sort(pdf.TextVertical(chars))
	for i := 0; i < len(chars); {
		j := i + 1
		for j < len(chars) && chars[j].Y == chars[i].Y {
			j++
		}
		var end float64
		// Split line into phrases
		for k := i; k < j; {
			ck := &chars[k]
			s := ck.S
			end = ck.X + ck.W
			charSpace := ck.FontSize / 6
			wordSpace := ck.FontSize * 2 / 3
			l := k + 1
			for l < j {
				// Grow word
				cl := &chars[l]
				if sameFont(cl.Font, ck.Font) && math.Abs(cl.FontSize-ck.FontSize) < 0.1 && cl.X <= end+charSpace {
					s += cl.S
					end = cl.X + cl.W
					l++
					continue
				}
				// Add space to phrase before next word
				if sameFont(cl.Font, ck.Font) && math.Abs(cl.FontSize-ck.FontSize) < 0.1 && cl.X <= end+wordSpace {
					s += " " + cl.S
					end = cl.X + cl.W
					l++
					continue
				}
				break
			}
			f := ck.Font
			f = strings.TrimSuffix(f, ",Italic")
			f = strings.TrimSuffix(f, "-Italic")
			words = append(words, pdf.Text{Font: f, FontSize: ck.FontSize, X: ck.X, Y: ck.Y, W: end - ck.X, S: s})
			k = l
		}
		words[len(words)-1].S += "\n"
		i = j
	}
	return words
}

func sameFont(f1, f2 string) bool {
	f1 = strings.TrimSuffix(f1, ",Italic")
	f1 = strings.TrimSuffix(f1, "-Italic")
	f2 = strings.TrimSuffix(f2, ",Italic")
	f2 = strings.TrimSuffix(f2, "-Italic")
	return strings.TrimSuffix(f1, ",Italic") == strings.TrimSuffix(f2, ",Italic") || f1 == "Symbol" || f2 == "Symbol" || f1 == "TimesNewRoman" || f2 == "TimesNewRoman"
}
