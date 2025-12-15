package aria

import "github.com/vapstack/htm"

// Label defines a string value that labels the current element.
func Label(v string) htm.Mod {
	return htm.Attr("aria-label", v)
}

// SetLabel defines a string value that labels the current element.
func SetLabel(n *htm.Node, v string) {
	n.Attr("aria-label", v)
}

// LabelledBy identifies the element (or elements) that labels the current element.
func LabelledBy(id string) htm.Mod {
	return htm.Attr("aria-labelledby", id)
}

// SetLabelledBy identifies the element (or elements) that labels the current element.
func SetLabelledBy(n *htm.Node, id string) {
	n.Attr("aria-labelledby", id)
}

// DescribedBy identifies the element (or elements) that describes the object.
func DescribedBy(id string) htm.Mod {
	return htm.Attr("aria-describedby", id)
}

// SetDescribedBy identifies the element (or elements) that describes the object.
func SetDescribedBy(n *htm.Node, id string) {
	n.Attr("aria-describedby", id)
}

// Details identifies the element that provides a detailed, extended description for the object.
func Details(id string) htm.Mod {
	return htm.Attr("aria-details", id)
}

// SetDetails identifies the element that provides a detailed, extended description for the object.
func SetDetails(n *htm.Node, id string) {
	n.Attr("aria-details", id)
}

// RoleDescription defines a human-readable, author-localized description for the role of an element.
func RoleDescription(v string) htm.Mod {
	return htm.Attr("aria-roledescription", v)
}

// SetRoleDescription defines a human-readable, author-localized description for the role of an element.
func SetRoleDescription(n *htm.Node, v string) {
	n.Attr("aria-roledescription", v)
}

/**/

// Hidden indicates whether the element is exposed to an accessibility API.
func Hidden(v ...bool) htm.Mod {
	return htm.AttrBool("aria-hidden", v...)
}

// SetHidden indicates whether the element is exposed to an accessibility API.
func SetHidden(n *htm.Node, v ...bool) {
	n.AttrBool("aria-hidden", v...)
}

// Disabled indicates that the element is perceivable but disabled.
func Disabled(v ...bool) htm.Mod {
	return htm.AttrBool("aria-disabled", v...)
}

// SetDisabled indicates that the element is perceivable but disabled.
func SetDisabled(n *htm.Node, v ...bool) {
	n.AttrBool("aria-disabled", v...)
}

// Expanded indicates whether the element, or another grouping element it controls, is currently expanded or collapsed.
func Expanded(v ...bool) htm.Mod {
	return htm.AttrBool("aria-expanded", v...)
}

// SetExpanded indicates whether the element, or another grouping element it controls, is currently expanded or collapsed.
func SetExpanded(n *htm.Node, v ...bool) {
	n.AttrBool("aria-expanded", v...)
}

// HasPopup indicates the availability and type of interactive popup element.
func HasPopup(v string) htm.Mod {
	return htm.Attr("aria-haspopup", v)
}

// SetHasPopup indicates the availability and type of interactive popup element.
func SetHasPopup(n *htm.Node, v string) {
	n.Attr("aria-haspopup", v)
}

// Pressed indicates the current "pressed" state of toggle buttons.
func Pressed(v ...bool) htm.Mod {
	return htm.AttrBool("aria-pressed", v...)
}

// SetPressed indicates the current "pressed" state of toggle buttons.
func SetPressed(n *htm.Node, v ...bool) {
	n.AttrBool("aria-pressed", v...)
}

// Checked indicates the current "checked" state of checkboxes, radio buttons, and other widgets.
func Checked(v ...bool) htm.Mod {
	return htm.AttrBool("aria-checked", v...)
}

// SetChecked indicates the current "checked" state of checkboxes, radio buttons, and other widgets.
func SetChecked(n *htm.Node, v ...bool) {
	n.AttrBool("aria-checked", v...)
}

// Selected indicates the current "selected" state of various widgets.
func Selected(v ...bool) htm.Mod {
	return htm.AttrBool("aria-selected", v...)
}

// SetSelected indicates the current "selected" state of various widgets.
func SetSelected(n *htm.Node, v ...bool) {
	n.AttrBool("aria-selected", v...)
}

// Modal indicates whether an element is modal when displayed.
func Modal(v ...bool) htm.Mod {
	return htm.AttrBool("aria-modal", v...)
}

// SetModal indicates whether an element is modal when displayed.
func SetModal(n *htm.Node, v ...bool) {
	n.AttrBool("aria-modal", v...)
}

// Current indicates the element that represents the current item within a container or set of related elements.
func Current(v string) htm.Mod {
	return htm.Attr("aria-current", v)
}

// SetCurrent indicates the element that represents the current item within a container or set of related elements.
func SetCurrent(n *htm.Node, v string) {
	n.Attr("aria-current", v)
}

// Required indicates that user input is required on the element before a form may be submitted.
func Required(v ...bool) htm.Mod {
	return htm.AttrBool("aria-required", v...)
}

// SetRequired indicates that user input is required on the element before a form may be submitted.
func SetRequired(n *htm.Node, v ...bool) {
	n.AttrBool("aria-required", v...)
}

// ReadOnly indicates that the element is not editable, but is otherwise operable.
func ReadOnly(v ...bool) htm.Mod {
	return htm.AttrBool("aria-readonly", v...)
}

// SetReadOnly indicates that the element is not editable, but is otherwise operable.
func SetReadOnly(n *htm.Node, v ...bool) {
	n.AttrBool("aria-readonly", v...)
}

// Placeholder defines a short hint (a word or short phrase) intended to aid the user with data entry.
func Placeholder(v string) htm.Mod {
	return htm.Attr("aria-placeholder", v)
}

// SetPlaceholder defines a short hint (a word or short phrase) intended to aid the user with data entry.
func SetPlaceholder(n *htm.Node, v string) {
	n.Attr("aria-placeholder", v)
}

/**/

// ValueMin defines the minimum allowed value for a range widget.
func ValueMin(v float64) htm.Mod {
	return htm.AttrValue("aria-valuemin", htm.Float(v))
}

// SetValueMin defines the minimum allowed value for a range widget.
func SetValueMin(n *htm.Node, v float64) {
	n.AttrValue("aria-valuemin", htm.Float(v))
}

// ValueMax defines the maximum allowed value for a range widget.
func ValueMax(v float64) htm.Mod {
	return htm.AttrValue("aria-valuemax", htm.Float(v))
}

// SetValueMax defines the maximum allowed value for a range widget.
func SetValueMax(n *htm.Node, v float64) {
	n.AttrValue("aria-valuemax", htm.Float(v))
}

// ValueNow defines the current value for a range widget.
func ValueNow(v float64) htm.Mod {
	return htm.AttrValue("aria-valuenow", htm.Float(v))
}

// SetValueNow defines the current value for a range widget.
func SetValueNow(n *htm.Node, v float64) {
	n.AttrValue("aria-valuenow", htm.Float(v))
}

// ValueText defines the human readable text alternative of aria-valuenow for a range widget.
func ValueText(v string) htm.Mod {
	return htm.Attr("aria-valuetext", v)
}

// SetValueText defines the human readable text alternative of aria-valuenow for a range widget.
func SetValueText(n *htm.Node, v string) {
	n.Attr("aria-valuetext", v)
}

/**/

// Controls identifies the element (or elements) whose contents or presence are controlled by the current element.
func Controls(id string) htm.Mod {
	return htm.Attr("aria-controls", id)
}

// SetControls identifies the element (or elements) whose contents or presence are controlled by the current element.
func SetControls(n *htm.Node, id string) {
	n.Attr("aria-controls", id)
}

// Owns identifies an element (or elements) in order to define a visual, functional, or contextual parent/child relationship.
func Owns(id string) htm.Mod {
	return htm.Attr("aria-owns", id)
}

// SetOwns identifies an element (or elements) in order to define a visual, functional, or contextual parent/child relationship.
func SetOwns(n *htm.Node, id string) {
	n.Attr("aria-owns", id)
}

// ActiveDescendant identifies the currently active element when focus is on a composite widget.
func ActiveDescendant(id string) htm.Mod {
	return htm.Attr("aria-activedescendant", id)
}

// SetActiveDescendant identifies the currently active element when focus is on a composite widget.
func SetActiveDescendant(n *htm.Node, id string) {
	n.Attr("aria-activedescendant", id)
}

// FlowTo identifies the next element (or elements) in an alternate reading order.
func FlowTo(id string) htm.Mod {
	return htm.Attr("aria-flowto", id)
}

// SetFlowTo identifies the next element (or elements) in an alternate reading order.
func SetFlowTo(n *htm.Node, id string) {
	n.Attr("aria-flowto", id)
}

/**/

// Live indicates that an element will be updated.
func Live(v string) htm.Mod {
	return htm.Attr("aria-live", v)
}

// SetLive indicates that an element will be updated.
func SetLive(n *htm.Node, v string) {
	n.Attr("aria-live", v)
}

// Atomic indicates whether assistive technologies will present all changes.
func Atomic(v ...bool) htm.Mod {
	return htm.AttrBool("aria-atomic", v...)
}

// SetAtomic indicates whether assistive technologies will present all changes.
func SetAtomic(n *htm.Node, v ...bool) {
	n.AttrBool("aria-atomic", v...)
}

// Relevant indicates what notifications the user agent will trigger.
func Relevant(v string) htm.Mod {
	return htm.Attr("aria-relevant", v)
}

// SetRelevant indicates what notifications the user agent will trigger.
func SetRelevant(n *htm.Node, v string) {
	n.Attr("aria-relevant", v)
}

// Busy indicates an element is being modified.
func Busy(v ...bool) htm.Mod {
	return htm.AttrBool("aria-busy", v...)
}

// SetBusy indicates an element is being modified.
func SetBusy(n *htm.Node, v ...bool) {
	n.AttrBool("aria-busy", v...)
}

/**/

// ColCount defines the total number of columns.
func ColCount(v int) htm.Mod {
	return htm.AttrValue("aria-colcount", htm.Int(v))
}

// SetColCount defines the total number of columns.
func SetColCount(n *htm.Node, v int) {
	n.AttrValue("aria-colcount", htm.Int(v))
}

// ColIndex defines an element's column index.
func ColIndex(v int) htm.Mod {
	return htm.AttrValue("aria-colindex", htm.Int(v))
}

// SetColIndex defines an element's column index.
func SetColIndex(n *htm.Node, v int) {
	n.AttrValue("aria-colindex", htm.Int(v))
}

// ColSpan defines the number of columns spanned by a cell.
func ColSpan(v int) htm.Mod {
	return htm.AttrValue("aria-colspan", htm.Int(v))
}

// SetColSpan defines the number of columns spanned by a cell.
func SetColSpan(n *htm.Node, v int) {
	n.AttrValue("aria-colspan", htm.Int(v))
}

// RowCount defines the total number of rows.
func RowCount(v int) htm.Mod {
	return htm.AttrValue("aria-rowcount", htm.Int(v))
}

// SetRowCount defines the total number of rows.
func SetRowCount(n *htm.Node, v int) {
	n.AttrValue("aria-rowcount", htm.Int(v))
}

// RowIndex defines an element's row index.
func RowIndex(v int) htm.Mod {
	return htm.AttrValue("aria-rowindex", htm.Int(v))
}

// SetRowIndex defines an element's row index.
func SetRowIndex(n *htm.Node, v int) {
	n.AttrValue("aria-rowindex", htm.Int(v))
}

// RowSpan defines the number of rows spanned.
func RowSpan(v int) htm.Mod {
	return htm.AttrValue("aria-rowspan", htm.Int(v))
}

// SetRowSpan defines the number of rows spanned.
func SetRowSpan(n *htm.Node, v int) {
	n.AttrValue("aria-rowspan", htm.Int(v))
}

// Level defines the hierarchical level of an element.
func Level(v int) htm.Mod {
	return htm.AttrValue("aria-level", htm.Int(v))
}

// SetLevel defines the hierarchical level of an element.
func SetLevel(n *htm.Node, v int) {
	n.AttrValue("aria-level", htm.Int(v))
}

// PosInSet defines an element's position in the set.
func PosInSet(v int) htm.Mod {
	return htm.AttrValue("aria-posinset", htm.Int(v))
}

// SetPosInSet defines an element's position in the set.
func SetPosInSet(n *htm.Node, v int) {
	n.AttrValue("aria-posinset", htm.Int(v))
}

// SetSize defines the number of items in the set.
func SetSize(v int) htm.Mod {
	return htm.AttrValue("aria-setsize", htm.Int(v))
}

// SetSetSize defines the number of items in the set.
func SetSetSize(n *htm.Node, v int) {
	n.AttrValue("aria-setsize", htm.Int(v))
}

/**/

// Orientation indicates whether the element's orientation is horizontal or vertical.
func Orientation(v string) htm.Mod {
	return htm.Attr("aria-orientation", v)
}

// SetOrientation indicates whether the element's orientation is horizontal or vertical.
func SetOrientation(n *htm.Node, v string) {
	n.Attr("aria-orientation", v)
}

// Sort indicates if items are sorted.
func Sort(v string) htm.Mod {
	return htm.Attr("aria-sort", v)
}

// SetSort indicates if items are sorted.
func SetSort(n *htm.Node, v string) {
	n.Attr("aria-sort", v)
}

// KeyShortcuts indicates keyboard shortcuts.
func KeyShortcuts(v string) htm.Mod {
	return htm.Attr("aria-keyshortcuts", v)
}

// SetKeyShortcuts indicates keyboard shortcuts.
func SetKeyShortcuts(n *htm.Node, v string) {
	n.Attr("aria-keyshortcuts", v)
}

// Autocomplete indicates autocomplete behavior.
func Autocomplete(v string) htm.Mod {
	return htm.Attr("aria-autocomplete", v)
}

// SetAutocomplete indicates autocomplete behavior.
func SetAutocomplete(n *htm.Node, v string) {
	n.Attr("aria-autocomplete", v)
}

// Multiline indicates whether a text box accepts multiple lines.
func Multiline(v ...bool) htm.Mod {
	return htm.AttrBool("aria-multiline", v...)
}

// SetMultiline indicates whether a text box accepts multiple lines.
func SetMultiline(n *htm.Node, v ...bool) {
	n.AttrBool("aria-multiline", v...)
}

// Multiselectable indicates that the user may select more than one item.
func Multiselectable(v ...bool) htm.Mod {
	return htm.AttrBool("aria-multiselectable", v...)
}

// SetMultiselectable indicates that the user may select more than one item.
func SetMultiselectable(n *htm.Node, v ...bool) {
	n.AttrBool("aria-multiselectable", v...)
}

// Invalid indicates the entered value does not conform to the expected format.
func Invalid(v htm.TypedValue) htm.Mod {
	return htm.AttrValue("aria-invalid", v)
}

// SetInvalid indicates the entered value does not conform to the expected format.
func SetInvalid(n *htm.Node, v htm.TypedValue) {
	n.AttrValue("aria-invalid", v)
}
