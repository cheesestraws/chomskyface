package main

import (
	"log"

	"github.com/hironobu-s/go-corenlp"
	"github.com/hironobu-s/go-corenlp/connector" 
	
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"	
	"github.com/fogleman/gg"
)

func newDC(width int, height int) *gg.Context {
	// set up a graphics context
	g := gg.NewContext(width, height)
	
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 12})
	g.SetFontFace(face)
	
	g.SetRGB(1,1,1)
	g.Clear()
	
	return g
}

func main() {
	con := connector.NewHTTPClient(nil, "http://127.0.0.1:9000/")
	
	sentence := "I will put my pyjamas in the drawer marked pyjamas."
	doc, err := corenlp.Annotate(con, sentence)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	
	if len(doc.Sentences) != 1 {
		log.Printf("error: not a single sentence")
		return	
	}

	p, _, err := doc.Sentences[0].Parse()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
		
	t := TreeFromParse(p)
	
	g := newDC(5000, 5000)
	t.PrepareToDraw(g)
	
	g = newDC(t.GWidth(), t.GHeight())
	t.Draw(g)
	g.SavePNG("res.png")
}
