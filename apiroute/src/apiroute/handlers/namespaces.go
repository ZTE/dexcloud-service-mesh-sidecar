package handlers

import (
	"apiroute/managers"
	"apiroute/managers/configmanager"
)

func GetAllNamespaces() []string {
	if configmanager.CfgM.GetConfigInfo().EnableTest {
		return MockGetAllNamespaces()
	}
	nm := managers.GetNamespaceManager()
	return nm.GetNamespaceList()
}

func MockGetAllNamespaces() []string {
	namespaces := []string{"default", "zenap", "openo", "test"}
	return namespaces
}
