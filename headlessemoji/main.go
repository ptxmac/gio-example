// SPDX-License-Identifier: Unlicense OR MIT

package main

// Render color emoji to a PNG using a headless GPU window (no app.Window).

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	colEmoji "eliasnaur.com/font/noto/emoji/color"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/gpu/headless"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

const emojiSample = "🥳🧁\n🥳 🧁"

func main() {
	out := flag.String("o", "emoji.png", "output PNG path")
	flag.Parse()

	if err := writePNG(*out); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote", *out)
}

func writePNG(path string) error {
	const wPx, hPx = 400, 200
	win, err := headless.NewWindow(wPx, hPx)
	if err != nil {
		return err
	}
	defer win.Release()

	collection := gofont.Collection()
	faces, err := opentype.ParseCollection(colEmoji.TTF)
	if err != nil {
		return err
	}
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(append(collection, faces...)))

	sz := image.Pt(wPx, hPx)
	var ops op.Ops
	gtx := layout.Context{
		Ops:         &ops,
		Constraints: layout.Exact(sz),
	}
	layout.UniformInset(5).Layout(gtx, material.H4(th, emojiSample).Layout)

	if err := win.Frame(&ops); err != nil {
		return err
	}
	img := image.NewRGBA(image.Rectangle{Max: sz})
	if err := win.Screenshot(img); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}
