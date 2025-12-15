package ax

import (
	"strings"

	"github.com/vapstack/htm"
)

// Data defines a chunk of HTML as an Alpine component and provides the reactive data for that component.
func Data(value ...htm.TypedValue) htm.Mod {
	if len(value) > 0 {
		v := value[0]
		return htm.AttrValue("x-data", v)
	}
	return htm.Attr("x-data")
}

// SetData sets the x-data attribute on the node.
func SetData(n *htm.Node, value ...htm.TypedValue) {
	if len(value) > 0 {
		v := value[0]
		n.AttrValue("x-data", v)
	} else {
		n.Attr("x-data")
	}
}

// Init runs an expression when the component is initialized.
func Init(value ...htm.TypedValue) htm.Mod {
	if len(value) > 0 {
		v := value[0]
		return htm.AttrValue("x-init", v)
	}
	return htm.Attr("x-init")
}

// SetInit sets the x-init attribute on the node.
func SetInit(n *htm.Node, value ...htm.TypedValue) {
	if len(value) > 0 {
		v := value[0]
		n.AttrValue("x-init", v)
	} else {
		n.Attr("x-init")
	}
}

// Show toggles the visibility of an element based on the truthiness of a JavaScript expression.
func Show(v string) htm.Mod {
	return htm.Attr("x-show", v)
}

// SetShow sets the x-show attribute on the node.
func SetShow(n *htm.Node, v string) {
	n.Attr("x-show", v)
}

/**/

// Bind sets the value of an attribute to the result of a JavaScript expression.
func Bind(attr string, v string) htm.Mod {
	return htm.Attr(":"+attr, v)
	// return htm.Attr("x-bind:"+attr, js)
}

// SetBind sets an x-bind:[attr] attribute on the node.
func SetBind(n *htm.Node, attr string, v string) {
	n.Attr(":"+attr, v)
	// n.Attr("x-bind:"+attr, js)
}

// BindValue applies an object of attributes to the element (e.g. x-bind="{...}").
func BindValue(value htm.TypedValue) htm.Mod {
	return htm.AttrValue("x-bind", value)
}

// SetBindValue sets the x-bind attribute (object syntax) on the node.
func SetBindValue(n *htm.Node, value htm.TypedValue) {
	n.AttrValue("x-bind", value)
}

// On attaches an event listener to the element.
func On(event string, v string) htm.Mod {
	return htm.Attr("@"+event, v)
	// return htm.Attr("x-on:"+event, js)
}

// SetOn sets an x-on:[event] attribute on the node.
func SetOn(n *htm.Node, event string, v string) {
	n.Attr("@"+event, v)
	// n.Attr("x-on:"+event, js)
}

/**/

// Text sets the text content of an element to the result of a JavaScript expression.
func Text(v string) htm.Mod {
	return htm.Attr("x-text", v)
}

// SetText sets the x-text attribute on the node.
func SetText(n *htm.Node, v string) {
	n.Attr("x-text", v)
}

// HTML sets the inner HTML of an element to the result of a JavaScript expression.
func HTML(v string) htm.Mod {
	return htm.Attr("x-html", v)
}

// SetHTML sets the x-html attribute on the node.
func SetHTML(n *htm.Node, v string) {
	n.Attr("x-html", v)
}

/**/

// Model binds the value of an input element to a specific data property.
func Model(v string) htm.Mod {
	return htm.Attr("x-model", v)
}

// SetModel sets the x-model attribute on the node.
func SetModel(n *htm.Node, v string) {
	n.Attr("x-model", v)
}

// ModelValue binds the value of an input element to a specific data property.
// It also allows to pass optional modifiers.
func ModelValue(modifiers string, v htm.TypedValue) htm.Mod {
	if modifiers != "" {
		return htm.AttrValue("x-model."+modifiers, v)
	}
	return htm.AttrValue("x-model."+modifiers, v)
}

// SetModelValue sets the x-model attribute on the node.
// It also allows to pass optional modifiers.
func SetModelValue(n *htm.Node, modifiers string, v htm.TypedValue) {
	if modifiers != "" {
		n.AttrValue("x-model."+modifiers, v)
	} else {
		n.AttrValue("x-model."+modifiers, v)
	}
}

// Modelable allows to expose any Alpine property as the target of an x-model directive.
func Modelable(v string) htm.Mod {
	return htm.Attr("x-modelable", v)
}

// SetModelable sets the x-modelable attribute on the node.
func SetModelable(n *htm.Node, v string) {
	n.Attr("x-modelable", v)
}

/**/

// For allows to iterate over an array or object to create DOM elements.
// Must be used on a <template> tag.
func For(v string) htm.Mod {
	return htm.Attr("x-for", v)
}

// SetFor sets the x-for attribute on the node.
func SetFor(n *htm.Node, v string) {
	n.Attr("x-for", v)
}

// If allows to toggle an element in and out of the DOM based on a JS expression.
// Must be used on a <template> tag.
func If(v string) htm.Mod {
	return htm.Attr("x-if", v)
}

// SetIf sets the x-if attribute on the node.
func SetIf(n *htm.Node, v string) {
	n.Attr("x-if", v)
}

/**/

// Effect allows to react to changes in data by running a callback.
func Effect(v string) htm.Mod {
	return htm.Attr("x-effect", v)
}

// SetEffect sets the x-effect attribute on the node.
func SetEffect(n *htm.Node, v string) {
	n.Attr("x-effect", v)
}

// Ref allows to reference elements directly in JavaScript using $refs.
func Ref(name string) htm.Mod {
	return htm.Attr("x-ref", name)
}

// SetRef sets the x-ref attribute on the node.
func SetRef(n *htm.Node, name string) {
	n.Attr("x-ref", name)
}

// Cloak hides the element until Alpine has fully initialized.
func Cloak(v ...bool) htm.Mod {
	return htm.AttrBool("x-cloak", v...)
}

// SetCloak sets the x-cloak attribute on the node.
func SetCloak(n *htm.Node, v ...bool) {
	n.AttrBool("x-cloak", v...)
}

// Ignore prevents Alpine from initializing elements within a block of HTML.
func Ignore(v ...bool) htm.Mod {
	return htm.AttrBool("x-ignore", v...)
}

// SetIgnore sets the x-ignore attribute on the node.
func SetIgnore(n *htm.Node, v ...bool) {
	n.AttrBool("x-ignore", v...)
}

// ID allows to declare a new scope for $id().
func ID(v string) htm.Mod {
	return htm.Attr("x-id", v)
}

// SetID sets the x-id attribute on the node.
func SetID(n *htm.Node, v string) {
	n.Attr("x-id", v)
}

// Teleport allows to transport part of template to another part of the DOM.
// selector: the DOM selector to append the content to (e.g. "body").
func Teleport(selector string) htm.Mod {
	return htm.Attr("x-teleport", selector)
}

// SetTeleport sets the x-teleport attribute on the node.
func SetTeleport(n *htm.Node, selector string) {
	n.Attr("x-teleport", selector)
}

/**/

// Transition enables standard transitions on an element.
func Transition(modifier ...string) htm.Mod {
	switch len(modifier) {
	case 0:
		return htm.Attr("x-transition")
	case 1:
		return htm.Attr("x-transition." + modifier[0])
	default:
		return htm.Attr("x-transition." + strings.Join(modifier, "."))
	}
}

// SetTransition sets the x-transition attribute on the node.
func SetTransition(n *htm.Node, modifier ...string) {
	switch len(modifier) {
	case 0:
		n.Attr("x-transition")
	case 1:
		n.Attr("x-transition." + modifier[0])
	default:
		n.Attr("x-transition." + strings.Join(modifier, "."))
	}
}

// TransitionStage applies specific transition classes (e.g. x-transition:enter="...").
// Stages: enter, enter-start, enter-end, leave, leave-start, leave-end.
func TransitionStage(stage string, classes string) htm.Mod {
	return htm.Attr("x-transition:"+stage, classes)
}

// SetTransitionStage applies specific transition classes (e.g. x-transition:enter="...").
// Stages: enter, enter-start, enter-end, leave, leave-start, leave-end.
func SetTransitionStage(n *htm.Node, stage string, classes string) {
	n.Attr("x-transition:"+stage, classes)
}

/**/

func Key(v string) htm.Mod {
	return htm.Attr(":key", v)
}

func SetKey(n *htm.Node, v string) {
	n.Attr(":key", v)
}
