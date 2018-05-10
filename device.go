package godevices

import "github.com/reimashi/godevices/interfaces"

type Device interface {
	GetModel() string
	GetVendor() string
	GetInterfaceType() interfaces.Type
}