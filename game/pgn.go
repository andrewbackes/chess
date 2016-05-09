package game

func EmptyTags() map[string]string {
	tags := make(map[string]string)
	tags["Event"] = ""
	tags["Site"] = ""
	tags["Date"] = ""
	tags["Round"] = ""
	tags["White"] = ""
	tags["Black"] = ""
	tags["Result"] = ""
	return tags
}
