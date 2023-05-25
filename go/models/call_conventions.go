package models

type Callconvention uint8

const (
	Cdecl Callconvention = iota
	Fastcall
	Thiscall
)
