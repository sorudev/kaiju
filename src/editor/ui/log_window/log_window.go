/******************************************************************************/
/* log_window.go                                                              */
/******************************************************************************/
/*                           This file is part of:                            */
/*                                KAIJU ENGINE                                */
/*                          https://kaijuengine.org                           */
/******************************************************************************/
/* MIT License                                                                */
/*                                                                            */
/* Copyright (c) 2023-present Kaiju Engine authors (AUTHORS.md).              */
/* Copyright (c) 2015-present Brent Farris.                                   */
/*                                                                            */
/* May all those that this source may reach be blessed by the LORD and find   */
/* peace and joy in life.                                                     */
/* Everyone who drinks of this water will be thirsty again; but whoever       */
/* drinks of the water that I will give him shall never thirst; John 4:13-14  */
/*                                                                            */
/* Permission is hereby granted, free of charge, to any person obtaining a    */
/* copy of this software and associated documentation files (the "Software"), */
/* to deal in the Software without restriction, including without limitation  */
/* the rights to use, copy, modify, merge, publish, distribute, sublicense,   */
/* and/or sell copies of the Software, and to permit persons to whom the      */
/* Software is furnished to do so, subject to the following conditions:       */
/*                                                                            */
/* The above copyright, blessing, biblical verse, notice and                  */
/* this permission notice shall be included in all copies or                  */
/* substantial portions of the Software.                                      */
/*                                                                            */
/* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS    */
/* OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF                 */
/* MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.     */
/* IN NO EVENT SHALL THE /* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY    */
/* CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT  */
/* OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE      */
/* OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.                              */
/******************************************************************************/

package log_window

import (
	"kaiju/editor/cache/editor_cache"
	"kaiju/editor/ui/editor_window"
	"kaiju/engine"
	"kaiju/host_container"
	"kaiju/klib"
	"kaiju/markup"
	"kaiju/markup/document"
	"kaiju/systems/logging"
	"kaiju/ui"
	"slices"
	"strconv"
	"strings"
	"time"
)

type viewGroup = int

const (
	viewGroupAll viewGroup = iota
	viewGroupInfo
	viewGroupWarn
	viewGroupError
	viewGroupSelected
)

type visibleMessage struct {
	Time     string
	Message  string
	Trace    string
	Data     map[string]string
	Category string
}

func newVisibleMessage(msg string, trace []string, cat string) visibleMessage {
	mapping := logging.ToMap(msg)
	t, _ := time.Parse(time.RFC3339, mapping["time"])
	message := mapping["msg"]
	delete(mapping, "time")
	delete(mapping, "msg")
	return visibleMessage{
		Time:     t.Format(time.StampMilli),
		Message:  message,
		Trace:    strings.Join(trace, "\n"),
		Data:     mapping,
		Category: cat,
	}
}

type LogWindow struct {
	doc        *document.Document
	container  *host_container.Container
	Group      viewGroup
	all        []visibleMessage
	lastReload engine.FrameId
	logStream  *logging.LogStream
	infoEvtId  logging.EventId
	warnEvtId  logging.EventId
	errEvtId   logging.EventId
}

func (l *LogWindow) All() []visibleMessage {
	res := slices.Clone(l.all)
	slices.Reverse(res)
	return res
}

func (l *LogWindow) filter(typeName string) []visibleMessage {
	res := make([]visibleMessage, 0, len(l.all))
	for i := range l.all {
		if l.all[i].Category == typeName {
			res = append(res, l.all[i])
		}
	}
	return res
}

func (l *LogWindow) Infos() []visibleMessage {
	res := l.filter("info")
	slices.Reverse(res)
	return res
}

func (l *LogWindow) Warnings() []visibleMessage {
	res := l.filter("warn")
	slices.Reverse(res)
	return res
}

func (l *LogWindow) Errors() []visibleMessage {
	res := l.filter("error")
	slices.Reverse(res)
	return res
}

func New(logStream *logging.LogStream) *LogWindow {
	l := &LogWindow{
		lastReload: engine.InvalidFrameId,
		all:        make([]visibleMessage, 0),
		logStream:  logStream,
	}
	l.infoEvtId = logStream.OnInfo.Add(func(msg string) {
		l.all = append(l.all, newVisibleMessage(msg, []string{}, "info"))
		if l.container != nil {
			l.container.RunFunction(l.reloadUI)
		}
	})
	l.warnEvtId = logStream.OnWarn.Add(func(msg string, trace []string) {
		l.all = append(l.all, newVisibleMessage(msg, trace, "warn"))
		if l.container != nil {
			l.container.RunFunction(l.reloadUI)
		}
	})
	l.errEvtId = logStream.OnError.Add(func(msg string, trace []string) {
		l.all = append(l.all, newVisibleMessage(msg, trace, "error"))
		if l.container != nil {
			l.container.RunFunction(l.reloadUI)
		}
	})
	return l
}

func (l *LogWindow) Show(listing *editor_window.Listing) {
	if l.container != nil {
		l.container.Host.Window.Focus()
		return
	}
	l.container = host_container.New("Log Window", nil)
	editor_window.OpenWindow(l, engine.DefaultWindowWidth,
		engine.DefaultWindowWidth/3, -1, -1)
	listing.Add(l)
}

func (l *LogWindow) Init() {
	l.reloadUI()
}

func (l *LogWindow) Closed() {
	l.logStream.OnInfo.Remove(l.infoEvtId)
	l.logStream.OnWarn.Remove(l.warnEvtId)
	l.logStream.OnError.Remove(l.errEvtId)
	l.container = nil
	l.lastReload = engine.InvalidFrameId
}

func (l *LogWindow) Tag() string                          { return editor_cache.LogWindow }
func (l *LogWindow) Container() *host_container.Container { return l.container }

func (l *LogWindow) clearAll(e *document.DocElement) {
	l.all = l.all[:0]
	l.reloadUI()
}

func (l *LogWindow) deactivateGroups() {
	all, _ := l.doc.GetElementById("all")
	info, _ := l.doc.GetElementById("info")
	warn, _ := l.doc.GetElementById("warn")
	err, _ := l.doc.GetElementById("error")
	selected, _ := l.doc.GetElementById("selected")
	all.UI.Entity().Deactivate()
	info.UI.Entity().Deactivate()
	warn.UI.Entity().Deactivate()
	err.UI.Entity().Deactivate()
	selected.UI.Entity().Deactivate()
	ab, _ := l.doc.GetElementById("allBtn")
	ib, _ := l.doc.GetElementById("infoBtn")
	wb, _ := l.doc.GetElementById("warningsBtn")
	eb, _ := l.doc.GetElementById("errorsBtn")
	sb, _ := l.doc.GetElementById("selectedBtn")
	ab.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("normal")
	ib.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("normal")
	wb.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("normal")
	eb.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("normal")
	sb.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("normal")
}

func (l *LogWindow) showCurrent() {
	switch l.Group {
	case viewGroupAll:
		l.showAll(nil)
	case viewGroupInfo:
		l.showInfos(nil)
	case viewGroupWarn:
		l.showWarns(nil)
	case viewGroupError:
		l.showErrors(nil)
	case viewGroupSelected:
		l.showSelected(nil)
	}
}

func (l *LogWindow) showAll(*document.DocElement) {
	l.Group = viewGroupAll
	l.deactivateGroups()
	e, _ := l.doc.GetElementById("all")
	b, _ := l.doc.GetElementById("allBtn")
	e.UI.Entity().Activate()
	b.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("bolder")
}

func (l *LogWindow) showInfos(*document.DocElement) {
	l.Group = viewGroupInfo
	l.deactivateGroups()
	e, _ := l.doc.GetElementById("info")
	b, _ := l.doc.GetElementById("infoBtn")
	e.UI.Entity().Activate()
	b.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("bolder")
}

func (l *LogWindow) showWarns(*document.DocElement) {
	l.Group = viewGroupWarn
	l.deactivateGroups()
	e, _ := l.doc.GetElementById("warn")
	b, _ := l.doc.GetElementById("warningsBtn")
	e.UI.Entity().Activate()
	b.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("bolder")
}

func (l *LogWindow) showErrors(*document.DocElement) {
	l.Group = viewGroupError
	l.deactivateGroups()
	e, _ := l.doc.GetElementById("error")
	b, _ := l.doc.GetElementById("errorsBtn")
	e.UI.Entity().Activate()
	b.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("bolder")
}

func (l *LogWindow) showSelected(*document.DocElement) {
	l.Group = viewGroupSelected
	l.deactivateGroups()
	e, _ := l.doc.GetElementById("selected")
	b, _ := l.doc.GetElementById("selectedBtn")
	e.UI.Entity().Activate()
	b.HTML.Children[0].DocumentElement.UI.(*ui.Label).SetFontWeight("bolder")
}

func (l *LogWindow) selectEntry(e *document.DocElement) {
	if id, err := strconv.Atoi(e.HTML.Attribute("data-entry")); err == nil {
		var target []visibleMessage
		switch l.Group {
		case viewGroupAll:
			target = l.all
		case viewGroupInfo:
			target = l.filter("info")
		case viewGroupWarn:
			target = l.filter("warn")
		case viewGroupError:
			target = l.filter("error")
		}
		if id >= 0 && id < len(target) {
			// The lists are printed in reverse order, so we invert the index
			id = len(target) - id - 1
			selectedElm, _ := l.doc.GetElementById("selected")
			lbl := selectedElm.HTML.Children[0].DocumentElement.UI.(*ui.Label)
			sb := strings.Builder{}
			sb.WriteString(target[id].Time)
			sb.WriteRune('\n')
			sb.WriteString(target[id].Message)
			sb.WriteRune('\n')
			for k, v := range target[id].Data {
				sb.WriteString(k)
				sb.WriteRune('=')
				sb.WriteString(v)
				sb.WriteRune('\n')
			}
			sb.WriteString(target[id].Trace)
			lbl.SetText(sb.String())
			l.showSelected(e)
		}
	}
}

func (l *LogWindow) reloadUI() {
	if l.container == nil {
		return
	}
	for _, e := range l.container.Host.Entities() {
		e.Destroy()
	}
	frame := l.container.Host.Frame()
	if l.lastReload == frame {
		return
	}
	l.lastReload = frame
	html := klib.MustReturn(l.container.Host.AssetDatabase().ReadText("editor/ui/log_window.html"))
	l.doc = markup.DocumentFromHTMLString(l.container.Host, html, "", l, map[string]func(*document.DocElement){
		"clearAll":     l.clearAll,
		"showAll":      l.showAll,
		"showInfos":    l.showInfos,
		"showWarns":    l.showWarns,
		"showErrors":   l.showErrors,
		"showSelected": l.showSelected,
		"selectEntry":  l.selectEntry,
	})
	l.showCurrent()
}
