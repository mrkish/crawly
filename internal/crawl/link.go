package crawl

type Link struct {
	URL   string
	Depth int
}

func mergeMap(existing, update map[string]Link) map[string]Link {
	merged := make(map[string]Link)
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range update {
		// if existingLink, ok := existing[k]; ok {
			// v.LinksTo = mergeMap(existingLink.LinksTo, v.LinksTo)
			// v.LinksFrom = mergeMap(existingLink.LinksFrom, v.LinksFrom)
		// }
		merged[k] = v
	}
	return merged
}
