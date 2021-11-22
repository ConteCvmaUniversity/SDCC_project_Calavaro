package main

import (
	"fmt"

	"github.com/IlConteCvma/SDCC_Project/menu"
)

func comm_seq(args ...string) error {
	//check args
	if len(args) == 0 {
		// set communicationState
		communicationState = Sequencer
		title := "\033[2J\033[H" + "\n" + "...SEQUENCER MENU..."
		commandOptions := []menu.CommandOption{
			menu.CommandOption{Command: "send", Description: "Send message use: send arg1 arg2 ...", Function: sendMessages},
			menu.CommandOption{Command: "show", Description: "Show message received from sequencer", Function: showMsg},
			menu.CommandOption{Command: "menu", Description: "Show the menu options", Function: nil},
			menu.CommandOption{Command: "quit or exit", Description: "Close sequencer newMenu", Function: nil},
		}
		menuOptions := menu.NewMenuOptions("Insert command > ", 0)
		newMenu := menu.NewMenu(title, commandOptions, menuOptions)
		newMenu.Start()
	} else {
		fmt.Println("expected only seq command")
	}

	return nil
}

func comm_scalar_clock(args ...string) error {
	if len(args) == 0 {
		// set communicationState
		communicationState = Scalar
		title := "\033[2J\033[H" + "\n" + "...SCALAR CLOCK MENU..."
		commandOptions := []menu.CommandOption{
			menu.CommandOption{Command: "send", Description: "Send message use: send arg1 arg2 ...", Function: sendMessages},
			menu.CommandOption{Command: "show", Description: "Show message received from scalar communication", Function: showMsg},
			menu.CommandOption{Command: "menu", Description: "Show the menu options", Function: nil},
			menu.CommandOption{Command: "quit or exit", Description: "Close scalar clock newMenu", Function: nil},
			menu.CommandOption{Command: "test", Description: "Send message use: send arg1 arg2 ...", Function: test},
		}
		menuOptions := menu.NewMenuOptions("Insert command > ", 0)
		newMenu := menu.NewMenu(title, commandOptions, menuOptions)
		newMenu.Start()
	} else {
		fmt.Println("expected only sc command")
	}

	return nil
}

func comm_vector_clock(args ...string) error {
	if len(args) == 0 {
		// set communicationState
		communicationState = Vector
		title := "\033[2J\033[H" + "\n" + "...VECTOR CLOCK MENU..."
		commandOptions := []menu.CommandOption{
			menu.CommandOption{Command: "send", Description: "Send message use: send arg1 arg2 ...", Function: sendMessages},
			menu.CommandOption{Command: "show", Description: "Show message received from vector communication", Function: showMsg},
			menu.CommandOption{Command: "menu", Description: "Show the newMenu options", Function: nil},
			menu.CommandOption{Command: "quit or exit", Description: "Close vector clock newMenu", Function: nil},
		}
		menuOptions := menu.NewMenuOptions("Insert command > ", 0)
		newMenu := menu.NewMenu(title, commandOptions, menuOptions)
		newMenu.Start()

	} else {
		fmt.Println("expected only vc command")
	}

	return nil
}

func open_menu() {

	//fmt.Println("\033[2J\033[H")

	commandOptions := []menu.CommandOption{
		menu.CommandOption{Command: "seq", Description: "Run communication whit sequencer", Function: comm_seq},
		menu.CommandOption{Command: "sc", Description: "Run scalar clock communication", Function: comm_scalar_clock},
		menu.CommandOption{Command: "vc", Description: "Run vector clock communication", Function: comm_vector_clock},
		menu.CommandOption{Command: "menu", Description: "Show the menu options", Function: nil},
		menu.CommandOption{Command: "quit or exit", Description: "Close application", Function: nil},
	}
	menuOptions := menu.NewMenuOptions("Insert command > ", 0)
	newMenu := menu.NewMenu("...MAIN MENU...", commandOptions, menuOptions)
	newMenu.Start()

}
