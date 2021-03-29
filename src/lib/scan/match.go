package scan

type match struct {
	//match <service> <pattern> <patternopt> [<versioninfo>]
	soft        bool
	service     string
	pattern     string
	versioninfo *finger
}

func newMatch() *match {
	return &match{
		soft:        false,
		service:     "",
		pattern:     "",
		versioninfo: newFinger(),
	}
}
