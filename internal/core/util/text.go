package util

type text struct{}

var Text text

func (text) Coalesce(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
