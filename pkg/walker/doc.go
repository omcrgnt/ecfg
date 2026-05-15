// Package walker traverses struct trees via Provider and Handler.
//
// RuntimeProvider walks live values (reflect); SchemaProvider walks types only
// (go/types or reflect.Type). Runtime slice/map indices use walkContent;
// schema slice/map use ElemProvider without indices.
package walker
