package analyser

type HTMLVersion string

// Summary represents a summary of HTML content.
type Summary struct {
	version              string              // version represents the HTML version.
	title                string              // title represents the HTML document title.
	headerCount          map[string]int      // headerCount represents the count of each header type.
	internalLinksMap     map[string]struct{} // internalLinksMap represents internal links found in the HTML.
	externalLinksMap     map[string]struct{} // externalLinksMap represents external links found in the HTML.
	inaccessibleLinksMap map[string]struct{} // inaccessibleLinksMap represents inaccessible links found in the HTML.
	hasLoginForm         bool                // hasLoginForm represents if the HTML contains a login form.
}

// NewSummary creates a new instance of Summary
func NewSummary() *Summary {
	return &Summary{
		headerCount:          make(map[string]int),
		internalLinksMap:     map[string]struct{}{},
		externalLinksMap:     map[string]struct{}{},
		inaccessibleLinksMap: map[string]struct{}{},
	}
}

// SetVersion sets the HTML version
func (s *Summary) SetVersion(version string) {
	s.version = version
}

// SetTitle sets the HTML document title
func (s *Summary) SetTitle(title string) {
	s.title = title
}

// IncrementHeaderCount increments the count for the specified header type
func (s *Summary) IncrementHeaderCount(header string) {
	s.headerCount[header]++
}

// AddExternalLink adds a link to the externalLinks
func (s *Summary) AddExternalLink(link string) {
	if _, ok := s.externalLinksMap[link]; !ok {
		s.externalLinksMap[link] = struct{}{}
	}
}

// AddInternalLink adds a link to the internalLinks
func (s *Summary) AddInternalLink(link string) {
	if _, ok := s.internalLinksMap[link]; !ok {
		s.internalLinksMap[link] = struct{}{}
	}
}

// AddInaccessibleLink adds a link to the inaccessibleLinks
func (s *Summary) AddInaccessibleLink(link string) {
	if _, ok := s.inaccessibleLinksMap[link]; !ok {
		s.inaccessibleLinksMap[link] = struct{}{}
	}
}

// SetHasLoginForm sets the has login form flag
func (s *Summary) SetHasLoginForm(status bool) {
	s.hasLoginForm = status
}
