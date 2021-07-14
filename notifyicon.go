/**
 * Copyright (c) 2021 BlockDev AG
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"log"

	"github.com/lxn/walk"
)

func createNotifyIcon() {
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	model.bus.Subscribe("exit", func() {
		ni.Dispose()
	})

	if err := ni.SetIcon(model.icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("Mysterium Network - Node Launcher"); err != nil {
		log.Fatal(err)
	}

	// When the left mouse button is pressed, bring up our balloon.
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		model.ShowMain()
	})
	ni.MessageClicked().Attach(func() {})

	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})

	openUIAction := walk.NewAction()
	if err := openUIAction.SetText("Open &UI"); err != nil {
		log.Fatal(err)
	}
	openUIAction.Triggered().Attach(func() {
		model.openNodeUI()
	})

	if err := ni.ContextMenu().Actions().Add(openUIAction); err != nil {
		log.Fatal(err)
	}
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// Run the message loop.
	mw.Run()
}
