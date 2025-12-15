package la

import "github.com/vapstack/htm"

// See https://icons8.com/line-awesome

func Icon(name string, mods ...htm.Mod) *htm.Node {
	return htm.I().Class("la las").Class("la-" + name).Apply(mods)
}
