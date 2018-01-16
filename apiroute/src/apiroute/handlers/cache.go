package handlers

import (
	"apiroute/managers"
	"apiroute/util"
)

var GetNameList = func() (nameList []*util.NameAndNamespace) {
	nameList, _ = managers.GetDataSyncManager().GetNameList()
	return nameList
}

var GetVersionsAndPaths = func(namespace string, serviceName string) (versions map[string]string) {
	return managers.GetDataSyncManager().GetVersionsAndPaths(namespace, serviceName)
}
