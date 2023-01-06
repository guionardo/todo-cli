package todo

func ParseTags(tags []string) []string {
	tagsMap := make(map[string]struct{})
	for _, tag := range tags {
		if len(tag) > 0 {
			tagsMap[tag] = struct{}{}
		}
	}
	parsedTags := make([]string, len(tagsMap))
	index := 0
	for k := range tagsMap {
		parsedTags[index] = k
		index++
	}

	return parsedTags
}

func EqualTags(tags1 []string, tags2 []string) bool {
	if len(tags1) != len(tags2) {
		return false
	}

	for _, tag := range tags1 {
		found := false
		for _, tag2 := range tags2 {
			if tag == tag2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
