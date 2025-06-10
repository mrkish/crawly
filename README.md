# Crawly

## Design

### CLI

### Steps
- Pass a root URL and depth
- The root URL will be fetched
- Once the HTML is received, it will be parsed for any anchor tags
  - Store the pages that have been parsed
  - External anchor tags will be ignored
  - Relative (page-internal) links will be ignored
  - Duplicate links will not be traversed twice


Get HTML
Parse HTML body for links
Store any links found on the page
Return a list of links that can be parsed (depth + 1)
Limit the parsing up to the specified depth

## Issues
- 
