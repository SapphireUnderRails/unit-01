package main

func stringInArray(str string, list []string) bool {
	for _, i := range list {
		if i == str {
			return true
		}
	}

	return false
}

func removeStringFromArray(str string, list []string) []string {
	for i, j := range list {
		if j == str {
			return append(list[:i], list[i+1:]...)
		}
	}

	return list
}
