// Copyright (C) 2021-2023 Leonid Maslakov.

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
	"embed"
	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/netshare"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"html/template"
	textTemplate "text/template"
)

//go:embed data/*
var embFS embed.FS

type Data struct {
	DB  storage.DB
	Log logger.Logger

	RateLimitNew *netshare.RateLimitSystem
	RateLimitGet *netshare.RateLimitSystem

	Lexers      *[]string
	Locales     Locales
	LocalesList LocalesList
	Themes      Themes
	ThemesList  ThemesList

	StyleCSS       *textTemplate.Template
	ErrorPage      *template.Template
	Main           *template.Template
	MainJS         *[]byte
	HistoryJS      *textTemplate.Template
	CodeJS         *textTemplate.Template
	PastePage      *template.Template
	PasteJS        *textTemplate.Template
	PasteContinue  *template.Template
	Settings       *template.Template
	About          *template.Template
	TermsOfUse     *template.Template
	Authors        *template.Template
	License        *template.Template
	SourceCodePage *template.Template

	Docs        *template.Template
	DocsApiV1   *template.Template
	DocsApiLibs *template.Template

	EmbeddedPage     *template.Template
	EmbeddedHelpPage *template.Template

	Version *string

	TitleMaxLen *int
	BodyMaxLen  *int
	MaxLifeTime *int64

	ServerAbout      *string
	ServerRules      *string
	ServerTermsExist *bool
	ServerTermsOfUse *string

	AdminName *string
	AdminMail *string

	RobotsDisallow *bool

	LenPasswdFile *string

	UiDefaultLifeTime *string
	UiDefaultTheme    *string
}

func Load(db storage.DB, cfg config.Config) (*Data, error) {
	var data Data
	var err error

	// Setup base info
	data.DB = db
	data.Log = cfg.Log

	data.RateLimitNew = cfg.RateLimitNew
	data.RateLimitGet = cfg.RateLimitGet

	data.Version = &cfg.Version

	data.TitleMaxLen = &cfg.TitleMaxLen
	data.BodyMaxLen = &cfg.BodyMaxLen
	data.MaxLifeTime = &cfg.MaxLifeTime
	data.UiDefaultLifeTime = &cfg.UiDefaultLifetime
	data.UiDefaultTheme = &cfg.UiDefaultTheme
	data.LenPasswdFile = &cfg.LenPasswdFile

	data.ServerAbout = &cfg.ServerAbout
	data.ServerRules = &cfg.ServerRules
	data.ServerTermsOfUse = &cfg.ServerTermsOfUse

	serverTermsExist := false
	if cfg.ServerTermsOfUse != "" {
		serverTermsExist = true
	}
	data.ServerTermsExist = &serverTermsExist

	data.AdminName = &cfg.AdminName
	data.AdminMail = &cfg.AdminMail

	data.RobotsDisallow = &cfg.RobotsDisallow

	// Get Chroma lexers
	lexers := chromaLexers.Names(false)
	data.Lexers = &lexers

	// Load locales
	data.Locales, data.LocalesList, err = loadLocales(embFS, "data/locale")
	if err != nil {
		return nil, err
	}

	// Load themes
	data.Themes, data.ThemesList, err = loadThemes(cfg.UiThemesDir, data.LocalesList, data.UiDefaultTheme)
	if err != nil {
		return nil, err
	}

	// style.css file
	data.StyleCSS, err = textTemplate.ParseFS(embFS, "data/style.css")
	if err != nil {
		return nil, err
	}

	// main.tmpl
	data.Main, err = template.ParseFS(embFS, "data/base.tmpl", "data/main.tmpl")
	if err != nil {
		return nil, err
	}

	// main.js
	mainJS, err := embFS.ReadFile("data/main.js")
	if err != nil {
		return nil, err
	}
	data.MainJS = &mainJS

	// history.js
	data.HistoryJS, err = textTemplate.ParseFS(embFS, "data/history.js")
	if err != nil {
		return nil, err
	}

	// code.js
	data.CodeJS, err = textTemplate.ParseFS(embFS, "data/code.js")
	if err != nil {
		return nil, err
	}

	// paste.tmpl
	data.PastePage, err = template.ParseFS(embFS, "data/base.tmpl", "data/paste.tmpl")
	if err != nil {
		return nil, err
	}

	// paste.js
	data.PasteJS, err = textTemplate.ParseFS(embFS, "data/paste.js")
	if err != nil {
		return nil, err
	}

	// paste_continue.tmpl
	data.PasteContinue, err = template.ParseFS(embFS, "data/base.tmpl", "data/paste_continue.tmpl")
	if err != nil {
		return nil, err
	}

	// settings.tmpl
	data.Settings, err = template.ParseFS(embFS, "data/base.tmpl", "data/settings.tmpl")
	if err != nil {
		return nil, err
	}

	// about.tmpl
	data.About, err = template.ParseFS(embFS, "data/base.tmpl", "data/about.tmpl")
	if err != nil {
		return nil, err
	}

	// terms.tmpl
	data.TermsOfUse, err = template.ParseFS(embFS, "data/base.tmpl", "data/terms.tmpl")
	if err != nil {
		return nil, err
	}

	// authors.tmpl
	data.Authors, err = template.ParseFS(embFS, "data/base.tmpl", "data/authors.tmpl")
	if err != nil {
		return nil, err
	}

	// license.tmpl
	data.License, err = template.ParseFS(embFS, "data/base.tmpl", "data/license.tmpl")
	if err != nil {
		return nil, err
	}

	// source_code.tmpl
	data.SourceCodePage, err = template.ParseFS(embFS, "data/base.tmpl", "data/source_code.tmpl")
	if err != nil {
		return nil, err
	}

	// docs.tmpl
	data.Docs, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs.tmpl")
	if err != nil {
		return nil, err
	}

	// docs_apiv1.tmpl
	data.DocsApiV1, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs_apiv1.tmpl")
	if err != nil {
		return nil, err
	}

	// docs_api_libs.tmpl
	data.DocsApiLibs, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs_api_libs.tmpl")
	if err != nil {
		return nil, err
	}

	// error.tmpl
	data.ErrorPage, err = template.ParseFS(embFS, "data/base.tmpl", "data/error.tmpl")
	if err != nil {
		return nil, err
	}

	// emb.tmpl
	data.EmbeddedPage, err = template.ParseFS(embFS, "data/emb.tmpl")
	if err != nil {
		return nil, err
	}

	// emb_help.tmpl
	data.EmbeddedHelpPage, err = template.ParseFS(embFS, "data/base.tmpl", "data/emb_help.tmpl")
	if err != nil {
		return nil, err
	}

	return &data, nil
}
