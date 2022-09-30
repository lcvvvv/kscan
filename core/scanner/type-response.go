package scanner

type Banner struct {
	Header   string
	Body     string
	Response string
	Cert     string
	Title    string
	Hash     string
	Icon     string
}

type FingerPrint struct {
	ProductName []string
	Hostname    string
	Domain      string
	MACAddr     string
}
