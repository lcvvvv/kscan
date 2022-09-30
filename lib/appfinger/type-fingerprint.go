package appfinger

type FingerPrint struct {
	ProductName []string
	Hostname    string
	Domain      string
	MACAddr     string
}

var emptyProductName []string

func New() *FingerPrint {
	return &FingerPrint{emptyProductName, "", "", ""}
}

func (f *FingerPrint) AddProduct(productName string) {
	f.ProductName = append(f.ProductName, productName)
}
