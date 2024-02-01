package properties

import (
	"errors"
	"kaiju/engine"
	"kaiju/ui"
	"kaiju/uimarkup/css/helpers"
	"kaiju/uimarkup/css/rules"
	"kaiju/uimarkup/markup"
)

func setChildrenFontSize(elm markup.DocElement, size string, host *engine.Host) {
	if elm.HTML.IsText() {
		lbl := elm.UI.(*ui.Label)
		size := helpers.NumFromLengthWithFont(size, host.Window,
			host.FontCache().EMSize(lbl.FontFace()))
		lbl.SetFontSize(size)
	} else {
		for _, child := range elm.HTML.Children {
			setChildrenFontSize(*child.DocumentElement, size, host)
		}
	}
}

func (p FontSize) Process(panel *ui.Panel, elm markup.DocElement, values []rules.PropertyValue, host *engine.Host) error {
	if len(values) != 1 {
		return errors.New("FontSize requires exactly 1 value")
	}
	setChildrenFontSize(elm, values[0].Str, host)
	return nil
}
