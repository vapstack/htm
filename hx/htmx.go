package hx

import "github.com/vapstack/htm"

/**/

// Get causes an element to issue a GET request to the specified URL.
func Get(url string) htm.Mod {
	return htm.Attr("hx-get", url)
}

// SetGet sets the hx-get attribute on the node.
func SetGet(n *htm.Node, url string) {
	n.Attr("hx-get", url)
}

// Post causes an element to issue a POST request to the specified URL.
func Post(url string) htm.Mod {
	return htm.Attr("hx-post", url)
}

// SetPost sets the hx-post attribute on the node.
func SetPost(n *htm.Node, url string) {
	n.Attr("hx-post", url)
}

// Put causes an element to issue a PUT request to the specified URL.
func Put(url string) htm.Mod {
	return htm.Attr("hx-put", url)
}

// SetPut sets the hx-put attribute on the node.
func SetPut(n *htm.Node, url string) {
	n.Attr("hx-put", url)
}

// Patch causes an element to issue a PATCH request to the specified URL.
func Patch(url string) htm.Mod {
	return htm.Attr("hx-patch", url)
}

// SetPatch sets the hx-patch attribute on the node.
func SetPatch(n *htm.Node, url string) {
	n.Attr("hx-patch", url)
}

// Delete causes an element to issue a DELETE request to the specified URL.
func Delete(url string) htm.Mod {
	return htm.Attr("hx-delete", url)
}

// SetDelete sets the hx-delete attribute on the node.
func SetDelete(n *htm.Node, url string) {
	n.Attr("hx-delete", url)
}

/**/

// Trigger specifies the event that triggers the request.
func Trigger(trigger string) htm.Mod {
	return htm.Attr("hx-trigger", trigger)
}

// SetTrigger sets the hx-trigger attribute on the node.
func SetTrigger(n *htm.Node, trigger string) {
	n.Attr("hx-trigger", trigger)
}

// Target specifies the target element to be swapped.
func Target(selector string) htm.Mod {
	return htm.Attr("hx-target", selector)
}

// SetTarget sets the hx-target attribute on the node.
func SetTarget(n *htm.Node, selector string) {
	n.Attr("hx-target", selector)
}

// Swap controls how the content is swapped in (e.g., "outerHTML", "beforeend").
func Swap(strategy string) htm.Mod {
	return htm.Attr("hx-swap", strategy)
}

// SetSwap sets the hx-swap attribute on the node.
func SetSwap(n *htm.Node, strategy string) {
	n.Attr("hx-swap", strategy)
}

// SwapOOB allows to specify that some content in a response should be swapped into the DOM somewhere other than the target.
func SwapOOB(v string) htm.Mod {
	return htm.Attr("hx-swap-oob", v)
}

// SetSwapOOB sets the hx-swap-oob attribute on the node.
func SetSwapOOB(n *htm.Node, v string) {
	n.Attr("hx-swap-oob", v)
}

// Select selects the content to be swapped in from a response.
func Select(selector string) htm.Mod {
	return htm.Attr("hx-select", selector)
}

// SetSelect sets the hx-select attribute on the node.
func SetSelect(n *htm.Node, selector string) {
	n.Attr("hx-select", selector)
}

// SelectOOB selects the content to be swapped in from a response, out of band.
func SelectOOB(selector string) htm.Mod {
	return htm.Attr("hx-select-oob", selector)
}

// SetSelectOOB sets the hx-select-oob attribute on the node.
func SetSelectOOB(n *htm.Node, selector string) {
	n.Attr("hx-select-oob", selector)
}

// Indicator specifies the element that is indicated during the request (e.g. a loading spinner).
func Indicator(selector string) htm.Mod {
	return htm.Attr("hx-indicator", selector)
}

// SetIndicator sets the hx-indicator attribute on the node.
func SetIndicator(n *htm.Node, selector string) {
	n.Attr("hx-indicator", selector)
}

/**/

// Vals allows to add to the parameters that will be submitted with the request.
// Typically accepts a JSON object.
func Vals(v htm.TypedValue) htm.Mod {
	return htm.AttrValue("hx-vals", v)
}

// SetVals sets the hx-vals attribute on the node.
func SetVals(n *htm.Node, v htm.TypedValue) {
	n.AttrValue("hx-vals", v)
}

// Params filters the parameters that will be submitted with a request.
// Values: "*", "none", "not <list>", or "<list>".
func Params(v string) htm.Mod {
	return htm.Attr("hx-params", v)
}

// SetParams sets the hx-params attribute on the node.
func SetParams(n *htm.Node, v string) {
	n.Attr("hx-params", v)
}

// Include includes additional values in a request.
func Include(selector string) htm.Mod {
	return htm.Attr("hx-include", selector)
}

// SetInclude sets the hx-include attribute on the node.
func SetInclude(n *htm.Node, selector string) {
	n.Attr("hx-include", selector)
}

// Headers adds to the headers that will be submitted with the request.
// Typically accepts a JSON object.
func Headers(v htm.TypedValue) htm.Mod {
	return htm.AttrValue("hx-headers", v)
}

// SetHeaders sets the hx-headers attribute on the node.
func SetHeaders(n *htm.Node, v htm.TypedValue) {
	n.AttrValue("hx-headers", v)
}

// Encoding allows to switch the request encoding from the usual application/x-www-form-urlencoded to multipart/form-data.
func Encoding(v string) htm.Mod {
	return htm.Attr("hx-encoding", v)
}

// SetEncoding sets the hx-encoding attribute on the node.
func SetEncoding(n *htm.Node, v string) {
	n.Attr("hx-encoding", v)
}

// Request allows to configure various aspects of the request.
// Typically accepts a JSON object.
func Request(v htm.TypedValue) htm.Mod {
	return htm.AttrValue("hx-request", v)
}

// SetRequest sets the hx-request attribute on the node.
func SetRequest(n *htm.Node, v htm.TypedValue) {
	n.AttrValue("hx-request", v)
}

// Sync allows to synchronize AJAX requests between multiple elements.
func Sync(selector string) htm.Mod {
	return htm.Attr("hx-sync", selector)
}

// SetSync sets the hx-sync attribute on the node.
func SetSync(n *htm.Node, selector string) {
	n.Attr("hx-sync", selector)
}

// Validate allows to force an element to validate itself before sending a request.
func Validate(v ...bool) htm.Mod {
	return htm.AttrBool("hx-validate", v...)
}

// SetValidate sets the hx-validate attribute on the node.
func SetValidate(n *htm.Node, v ...bool) {
	n.AttrBool("hx-validate", v...)
}

/**/

// PushURL pushes a new URL into the browser location history.
// Accepts: "true", "false", or a URL string.
func PushURL(v string) htm.Mod {
	return htm.Attr("hx-push-url", v)
}

// SetPushURL sets the hx-push-url attribute on the node.
func SetPushURL(n *htm.Node, v string) {
	n.Attr("hx-push-url", v)
}

// ReplaceURL replaces the current URL in the browser location history.
// Accepts: "true", "false", or a URL string.
func ReplaceURL(v string) htm.Mod {
	return htm.Attr("hx-replace-url", v)
}

// SetReplaceURL sets the hx-replace-url attribute on the node.
func SetReplaceURL(n *htm.Node, v string) {
	n.Attr("hx-replace-url", v)
}

// History prevents sensitive data from being saved to the history cache.
// Usually used as hx-history="false".
func History(v bool) htm.Mod {
	if !v {
		return htm.Attr("hx-history", "false")
	}
	return htm.AttrValue("hx-history", htm.Unset)
}

// SetHistory sets the hx-history attribute on the node.
func SetHistory(n *htm.Node, v bool) {
	if !v {
		n.Attr("hx-history", "false")
	}
	n.AttrValue("hx-history", htm.Unset)
}

// HistoryElt specifies the element to snapshot and restore during history navigation.
func HistoryElt() htm.Mod {
	return htm.Attr("hx-history-elt")
}

// SetHistoryElt sets the hx-history-elt attribute on the node.
func SetHistoryElt(n *htm.Node) {
	n.Attr("hx-history-elt")
}

/**/

// Boost progressively enhances anchors and forms to use AJAX.
func Boost(value ...bool) htm.Mod {
	if len(value) > 0 {
		if value[0] {
			return htm.Attr("hx-boost", "true")
		} else {
			return htm.Attr("hx-boost", "false")
		}
	}
	return htm.Attr("hx-boost", "true")
}

// SetBoost sets the hx-boost attribute on the node.
func SetBoost(n *htm.Node, value ...bool) {
	if len(value) > 0 {
		if value[0] {
			n.Attr("hx-boost", "true")
		} else {
			n.Attr("hx-boost", "false")
		}
	} else {
		n.Attr("hx-boost", "true")
	}
}

// Confirm shows a confirm() dialog before issuing a request.
func Confirm(msg string) htm.Mod {
	return htm.Attr("hx-confirm", msg)
}

// SetConfirm sets the hx-confirm attribute on the node.
func SetConfirm(n *htm.Node, msg string) {
	n.Attr("hx-confirm", msg)
}

// Preserve specifies that an element should be kept invariant across requests (e.g. video players).
func Preserve(v ...bool) htm.Mod {
	return htm.AttrBool("hx-preserve", v...)
}

// SetPreserve sets the hx-preserve attribute on the node.
func SetPreserve(n *htm.Node, v ...bool) {
	n.AttrBool("hx-preserve", v...)
}

// Disinherit allows to control the inheritance of htmx attributes.
func Disinherit(v string) htm.Mod {
	return htm.Attr("hx-disinherit", v)
}

// SetDisinherit sets the hx-disinherit attribute on the node.
func SetDisinherit(n *htm.Node, v string) {
	n.Attr("hx-disinherit", v)
}

// Disable disables htmx processing for a given element and its children.
func Disable(v ...bool) htm.Mod {
	return htm.AttrBool("hx-disable", v...)
}

// SetDisable sets the hx-disable attribute on the node.
func SetDisable(n *htm.Node, v ...bool) {
	n.AttrBool("hx-disable", v...)
}

// Ext enables an htmx extension for an element and all its children.
func Ext(extensions string) htm.Mod {
	return htm.Attr("hx-ext", extensions)
}

// SetExt sets the hx-ext attribute on the node.
func SetExt(n *htm.Node, extensions string) {
	n.Attr("hx-ext", extensions)
}

// On sets event handlers (hx-on:event).
func On(event string, js string) htm.Mod {
	return htm.Attr("hx-on:"+event, js)
}

// SetOn sets hx-on:[event] attribute on the node.
func SetOn(n *htm.Node, event string, js string) {
	n.Attr("hx-on:"+event, js)
}
