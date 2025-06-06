package crawl

import (
	// "log/slog"
	"net/url"
	// "strconv"
)

type Page struct {
	Link
	Root  *url.URL
	Links []Link
}

// type Pages []Page
//
// func (p Pages) LogValue() slog.Value {
// 	return slog.GroupValue(
// 		slog.Attr{Key: "pages", Value: slog.GroupValue(func(g *slog.Group) {
// 			for i, page := range p {
// 				g.Add(
// 					strconv.Itoa(i),
// 					slog.StringValue(page.URL),
// 				)
// 			}
// 		})})
// }

// slog.String("url", p.URL),
// slog.Int("depth", p.Depth),
// slog.Attr{Key: "page", value: func() slog.Attr{
//
// 	for i, page := range p {
//
// 	}
//
// func (ms MySlice) LogValue() slog.Value {
//     // Create a Group to hold the slice elements
//     group := slog.GroupValue(
//         slog.Attr{Key: "slice_elements", Value: slog.GroupValue(func(g *slog.Group) {
//             for i, val := range ms {
//                 g.Add(strconv.Itoa(i), slog.StringValue(val))
//             }
//         })},
//     )
//     return group
// }

type PageBuilder func(url string, depth int) Page

func pageFactory(root *url.URL) PageBuilder {
	return func(url string, depth int) Page {
		return Page{
			Root: root,
			Link: Link{
				URL:   url,
				Depth: depth,
			},
		}
	}
}
