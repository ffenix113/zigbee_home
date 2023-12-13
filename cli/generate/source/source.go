package source

type Source struct{}

func NewSource() *Source {
	return &Source{}
}

func (s *Source) WriteTo(srcDir string) error {
	return nil
}
