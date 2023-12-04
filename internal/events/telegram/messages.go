package telegram

const msgHelp = `I can save and keep your pages. I can offer you them to read.

To save the link, just send it to me.

To get a random link from your list, send me the command /rnd
After the /rnd the link will be removed from the list.
`

const msgHello = "Hi there! 😈\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🙅"
	msgNoSavedPages   = "You have no saved pages 😢"
	msgSaved          = "Saved! ✅"
	msgAlreadyExists  = "You already have this page in your list 🙌"
)
