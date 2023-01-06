package todo

func (c *Collection) EqualsTo(other *Collection) bool {

	if c == nil || other == nil || len(c.Items) != len(other.Items) {
		return false
	}

	for _, item := range c.Items {
		if otherItem, ok := other.Items[item.Id]; !ok || !otherItem.EqualsTo(item) {
			return false
		}
	}
	for _, otherItem := range other.Items {
		if item, ok := c.Items[otherItem.Id]; !ok || !item.EqualsTo(otherItem) {
			return false
		}
	}
	return true
}
