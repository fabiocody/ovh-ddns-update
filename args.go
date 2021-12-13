package main

type ArgsType struct {
	Database    string `default:".ovh-ddns-update.sqlite3" help:"database file"`
	Domain      string `arg:"positional,required" help:"OVH DynHost domain"`
	OvhId       string `arg:"positional,required" help:"OVH DynHost identifier"`
	OvhPassword string `arg:"positional,required" help:"OVH DynHost identifier password"`
}

var Args ArgsType
