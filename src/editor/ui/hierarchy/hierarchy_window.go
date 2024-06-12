/******************************************************************************/
/* hierarchy.go                                                               */
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

package hierarchy

import (
	"kaiju/editor/cache/editor_cache"
	"kaiju/editor/selection"
	"kaiju/editor/ui/drag_datas"
	"kaiju/engine"
	"kaiju/klib"
	"kaiju/markup"
	"kaiju/markup/document"
	"kaiju/matrix"
	"kaiju/systems/events"
	"kaiju/ui"
	"log/slog"
	"strings"
)

const sizeConfig = "hierarchyWindowSize"

type Hierarchy struct {
	host      *engine.Host
	selection *selection.Selection
	doc       *document.Document
	input     *ui.Input
	query     string
	uiGroup   *ui.Group
}

type entityEntry struct {
	Entity          *engine.Entity
	ShowingChildren bool
}

type hierarchyData struct {
	Entries []entityEntry
	Query   string
}

func (e entityEntry) Depth() int {
	depth := 0
	p := e.Entity
	for p.Parent != nil {
		depth++
		p = p.Parent
	}
	return depth
}

func New(host *engine.Host, selection *selection.Selection, uiGroup *ui.Group) *Hierarchy {
	h := &Hierarchy{
		host:      host,
		selection: selection,
		uiGroup:   uiGroup,
	}
	h.host.OnClose.Add(func() {
		if h.doc != nil {
			h.doc.Destroy()
		}
	})
	h.Reload()
	h.selection.Changed.Add(h.onSelectionChanged)
	return h
}

func (h *Hierarchy) Toggle() {
	if h.doc == nil {
		h.Show()
	} else {
		if h.doc.Elements[0].UI.Entity().IsActive() {
			h.Hide()
		} else {
			h.Show()
		}
	}
}

func (h *Hierarchy) Show() {
	if h.doc == nil {
		h.Reload()
	} else {
		h.doc.Activate()
	}
}

func (h *Hierarchy) Hide() {
	if h.doc != nil {
		h.doc.Deactivate()
	}
}

func (h *Hierarchy) orderEntitiesVisually() []entityEntry {
	allEntities := h.host.Entities()
	entries := make([]entityEntry, 0, len(allEntities))
	roots := make([]*engine.Entity, 0, len(allEntities))
	for _, entity := range allEntities {
		if entity.IsRoot() && !entity.EditorBindings.IsDeleted {
			roots = append(roots, entity)
		}
	}
	var addChildren func(*engine.Entity)
	addChildren = func(entity *engine.Entity) {
		if entity.EditorBindings.IsDeleted {
			return
		}
		entries = append(entries, entityEntry{entity, false})
		for _, c := range entity.Children {
			addChildren(c)
		}
	}
	for _, r := range roots {
		addChildren(r)
	}
	return entries
}

func (h *Hierarchy) filter(entries []entityEntry) []entityEntry {
	if h.query == "" {
		return entries
	}
	filtered := make([]entityEntry, 0, len(entries))
	// TODO:  Append the entire path to the entity kf not already appended
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Entity.Name()), h.query) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func (h *Hierarchy) Reload() {
	isActive := false
	if h.doc != nil {
		isActive = h.doc.Elements[0].UI.Entity().IsActive()
		h.doc.Destroy()
	}
	data := hierarchyData{
		Entries: h.filter(h.orderEntitiesVisually()),
		Query:   h.query,
	}
	host := h.host
	host.CreatingEditorEntities()
	h.doc = klib.MustReturn(markup.DocumentFromHTMLAsset(
		host, "editor/ui/hierarchy_window.html", data,
		map[string]func(*document.Element){
			"selectedEntity": h.selectedEntity,
			"dragStart":      h.dragStart,
			"drop":           h.drop,
			"dragEnter":      h.dragEnter,
			"dragExit":       h.dragExit,
			"resizeHover":    h.resizeHover,
			"resizeExit":     h.resizeExit,
			"resizeStart":    h.resizeStart,
			"resizeStop":     h.resizeStop,
		}))
	h.doc.SetGroup(h.uiGroup)
	host.DoneCreatingEditorEntities()
	if elm, ok := h.doc.GetElementById("searchInput"); !ok {
		slog.Error(`Failed to locate the "searchInput" for the hierarchy`)
	} else {
		h.input = elm.UI.(*ui.Input)
		h.input.AddEvent(ui.EventTypeSubmit, h.submit)
	}
	h.doc.Clean()
	if s, ok := editor_cache.EditorConfigValue(sizeConfig); ok {
		w, _ := h.doc.GetElementById("window")
		w.UIPanel.Layout().ScaleWidth(matrix.Float(s.(float64)))
	}
	if !isActive {
		h.doc.Deactivate()
	}
}

func (h *Hierarchy) submit() {
	h.query = strings.ToLower(strings.TrimSpace(h.input.Text()))
	h.Reload()
}

func (h *Hierarchy) onSelectionChanged() {
	elm, ok := h.doc.GetElementById("list")
	if !ok {
		slog.Error("Could not find hierarchy list, reopen the hierarchy window")
		return
	}
	for i := range elm.Children {
		elm.Children[i].UnEnforceColor()
	}
	for _, c := range elm.Children {
		id := engine.EntityId(c.Attribute("id"))
		for _, se := range h.selection.Entities() {
			if se.Id() == id {
				c.EnforceColor(matrix.ColorDarkBlue())
				break
			}
		}
	}
}

func (h *Hierarchy) selectedEntity(elm *document.Element) {
	id := engine.EntityId(elm.Attribute("id"))
	if e, ok := h.host.FindEntity(id); !ok {
		slog.Error("Could not find entity", slog.String("id", string(id)))
	} else {
		kb := &h.host.Window.Keyboard
		if kb.HasCtrl() {
			h.selection.Toggle(e)
		} else if kb.HasShift() {
			h.selection.Add(e)
		} else {
			h.selection.Set(e)
		}
		h.onSelectionChanged()
	}
}

func (h *Hierarchy) drop(elm *document.Element) {
	elm.UnEnforceColor()
	from := h.host.Window.Mouse.DragData().(*drag_datas.EntityIdDragData)
	if f, ok := h.host.FindEntity(from.EntityId); ok {
		toId := elm.Attribute("id")
		if toId != "" {
			to := engine.EntityId(toId)
			if t, ok := h.host.FindEntity(to); ok {
				f.SetParent(t)
				h.Reload()
			} else {
				slog.Error("Could not find drop target entity", slog.String("id", string(to)))
			}
		} else {
			f.SetParent(nil)
			h.Reload()
		}
	} else {
		slog.Error("Could not find drag entity", slog.String("id", string(from.EntityId)))
	}
}

func (h *Hierarchy) dragStart(elm *document.Element) {
	id := engine.EntityId(elm.Attribute("id"))
	h.host.Window.CursorSizeAll()
	h.host.Window.Mouse.SetDragData(&drag_datas.EntityIdDragData{id})
	elm.EnforceColor(matrix.ColorPurple())
	var eid events.Id
	eid = h.host.Window.Mouse.OnDragStop.Add(func() {
		h.host.Window.CursorStandard()
		h.host.Window.Mouse.OnDragStop.Remove(eid)
		elm.UnEnforceColor()
	})
}

func (h *Hierarchy) dragEnter(elm *document.Element) {
	myId := engine.EntityId(elm.Attribute("id"))
	if dd, ok := h.host.Window.Mouse.DragData().(*drag_datas.EntityIdDragData); !ok {
		return
	} else {
		if myId != dd.EntityId {
			elm.EnforceColor(matrix.ColorOrange())
		}
	}
}

func (h *Hierarchy) dragExit(elm *document.Element) {
	myId := engine.EntityId(elm.Attribute("id"))
	if dd, ok := h.host.Window.Mouse.DragData().(*drag_datas.EntityIdDragData); !ok {
		return
	} else {
		if myId != dd.EntityId {
			elm.UnEnforceColor()
		}
	}
}

func (h *Hierarchy) resizeHover(e *document.Element) {
	h.host.Window.CursorSizeWE()
}

func (h *Hierarchy) resizeExit(e *document.Element) {
	dd := h.host.Window.Mouse.DragData()
	if dd != h {
		h.host.Window.CursorStandard()
	}
}

func (h *Hierarchy) resizeStart(e *document.Element) {
	h.host.Window.CursorSizeWE()
	h.host.Window.Mouse.SetDragData(h)
}

func (h *Hierarchy) resizeStop(e *document.Element) {
	dd := h.host.Window.Mouse.DragData()
	if dd != h {
		return
	}
	h.host.Window.CursorStandard()
	w, _ := h.doc.GetElementById("window")
	s := w.UIPanel.Layout().PixelSize().Width()
	editor_cache.SetEditorConfigValue(sizeConfig, s)
}

func (h *Hierarchy) DragUpdate() {
	win, _ := h.doc.GetElementById("window")
	x := h.host.Window.Mouse.Position().X()
	w := h.host.Window.Width()
	if int(x) < w-100 {
		win.UIPanel.Layout().ScaleWidth(x)
	}
}
