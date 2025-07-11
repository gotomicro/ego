package transport

var customHeaderKeyStore = contextKeyStore{
	keyArr: make([]string, 0),
}

// SetHeaderKeys overrides custom keys with provided array.
func SetHeaderKeys(arr []string) {
	//length := len(arr)
	customHeaderKeyStore.keyArr = arr
	//customHeaderKeyStore.length = length
}

// CustomHeaderKeys returns custom content key list
func CustomHeaderKeys() []string {
	return customHeaderKeyStore.keyArr
}
