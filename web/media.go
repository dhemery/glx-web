package web

import "github.com/genealogix/glx/go-glx"

type Media struct {
	a  *Archive
	g  *glx.Media
	ID string
}

func newMedia(id string, gm *glx.Media, a *Archive) *Media {
	return &Media{
		a:  a,
		g:  gm,
		ID: id,
	}
}
