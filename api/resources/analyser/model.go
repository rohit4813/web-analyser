package analyser

// Summary represents a summary of HTML page.
type Summary struct {
	URL                  string              // URL represents the URL which was analysed.
	Version              string              // Version represents the HTML Version.
	Title                string              // Title represents the HTML page Title.
	HeadersCount         map[string]int      // HeadersCount represents the count of each header type.
	InternalLinksMap     map[string]struct{} // InternalLinksMap represents internal links found in the HTML page.
	ExternalLinksMap     map[string]struct{} // ExternalLinksMap represents external links found in the HTML page.
	InaccessibleLinksMap map[string]struct{} // InaccessibleLinksMap represents inaccessible links found in the HTML page.
	HasLoginForm         bool                // HasLoginForm represents if the HTML page contains a login form.
}

// NewSummary creates a new instance of Summary
func NewSummary(url string) *Summary {
	return &Summary{
		URL:                  url,
		HeadersCount:         make(map[string]int),
		InternalLinksMap:     map[string]struct{}{},
		ExternalLinksMap:     map[string]struct{}{},
		InaccessibleLinksMap: map[string]struct{}{},
	}
}

// SetURL sets the URL
func (s *Summary) SetURL(url string) {
	s.URL = url
}

// SetVersion sets the HTML Version
func (s *Summary) SetVersion(version string) {
	s.Version = version
}

// SetTitle sets the HTML page Title
func (s *Summary) SetTitle(title string) {
	s.Title = title
}

// IncrementHeadersCount increments the count for the specified header type
func (s *Summary) IncrementHeadersCount(header string) {
	s.HeadersCount[header]++
}

// AddExternalLink adds a link to the externalLinks
func (s *Summary) AddExternalLink(link string) {
	if _, ok := s.ExternalLinksMap[link]; !ok {
		s.ExternalLinksMap[link] = struct{}{}
	}
}

// AddInternalLink adds a link to the internalLinks
func (s *Summary) AddInternalLink(link string) {
	if _, ok := s.InternalLinksMap[link]; !ok {
		s.InternalLinksMap[link] = struct{}{}
	}
}

// AddInaccessibleLink adds a link to the inaccessibleLinks
func (s *Summary) AddInaccessibleLink(link string) {
	if _, ok := s.InaccessibleLinksMap[link]; !ok {
		s.InaccessibleLinksMap[link] = struct{}{}
	}
}

// SetHasLoginForm sets the has login form flag
func (s *Summary) SetHasLoginForm(status bool) {
	s.HasLoginForm = status
}
