package ui

import "github.com/aditya-K2/gspt/spt"

func OpenDeviceMenu() {
	m := NewMenu()
	cc := []string{}
	// TODO: Better Error Handling
	devices, err := spt.UserDevices()
	if err != nil {
		SendNotification(err.Error())
		return
	}
	for _, v := range devices {
		cc = append(cc, v.Name)
	}
	m.Content(cc)
	m.Title("Choose A Device")
	m.SetSelectionHandler(func(s int) {
		if err := spt.TransferPlayback(devices[s].ID); err != nil {
			SendNotification(err.Error())
		} else {
			RefreshProgress(true)
		}
	})
	Ui.Root.AddCenteredWidget(m)
}
