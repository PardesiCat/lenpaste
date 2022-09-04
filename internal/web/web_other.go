// Copyright (C) 2021-2022 Leonid Maslakov.

// This file is part of Lenpaste.

// Lenpaste is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// Lenpaste is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero Public License for more details.

// You should have received a copy of the GNU Affero Public License along with Lenpaste.
// If not, see <https://www.gnu.org/licenses/>.

package web

import (
	"html/template"
	"net/http"
)

type jsTmpl struct {
	Translate func(string, ...interface{}) template.HTML
}

func (data Data) StyleCSSHand(rw http.ResponseWriter, req *http.Request) {
	data.Log.HttpRequest(req)

	rw.Header().Set("Content-Type", "text/css")
	rw.Write(*data.StyleCSS)
}

func (data Data) MainJSHand(rw http.ResponseWriter, req *http.Request) {
	data.Log.HttpRequest(req)

	rw.Header().Set("Content-Type", "application/javascript")
	rw.Write(*data.MainJS)
}

func (data Data) CodeJSHand(rw http.ResponseWriter, req *http.Request) {
	data.Log.HttpRequest(req)

	rw.Header().Set("Content-Type", "application/javascript")
	data.CodeJS.Execute(rw, jsTmpl{Translate: data.Locales.findLocale(req).translate})
}

func (data Data) HistoryJSHand(rw http.ResponseWriter, req *http.Request) {
	data.Log.HttpRequest(req)

	rw.Header().Set("Content-Type", "application/javascript")
	data.HistoryJS.Execute(rw, jsTmpl{Translate: data.Locales.findLocale(req).translate})
}
