// Package pickertest is a harness for testing file-chooser components: an
// in-memory filesystem, a driver that sends keys and settles the commands a
// component answers with, and a script runner whose golden screens are
// readable by a person.
//
// It knows nothing of any one component: it needs only Init, Update, and
// View, so the same scripts can characterize one chooser and then hold
// another to it.
package pickertest
