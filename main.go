package main

import (
	"fmt"
	"strings"
)

type RowData struct {
	RankWebsite int
	Domain      string
	TLD         string
	IDN_TLD     string
	Valid       bool
	RefIPs      int
}

func GetTLD(domain string) (TLD string, IDN_TLD string) {
	var ListIDN_TLD = map[string]string{
		".com": ".co.id",
		".org": ".org.id",
		".gov": ".go.id",
	}

	// Mengambil TLD dari domain
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return "", ""
	}
	TLD = "." + parts[len(parts)-1] // Mendapatkan TLD terakhir

	if idn, ok := ListIDN_TLD[TLD]; ok {
		return TLD, idn
	}
	return TLD, TLD
}

func ProcessGetTLD(website RowData, ch chan RowData, chErr chan error) {
	if website.Domain == "" {
		chErr <- fmt.Errorf("domain name is empty")
		return
	}
	if !website.Valid {
		chErr <- fmt.Errorf("domain not valid")
		return
	}
	if website.RefIPs < 0 {
		chErr <- fmt.Errorf("domain RefIPs not valid")
		return
	}

	// Mendapatkan TLD dan IDN_TLD
	TLD, IDN_TLD := GetTLD(website.Domain)
	website.TLD = TLD
	website.IDN_TLD = IDN_TLD

	// Mengirim data yang sudah diisi ke channel
	ch <- website
}

// Gunakan variable ini sebagai goroutine di fungsi FilterAndFillData
var FuncProcessGetTLD = ProcessGetTLD

func FilterAndFillData(TLD string, data []RowData) ([]RowData, error) {
	ch := make(chan RowData, len(data))
	errCh := make(chan error)

	// Menggunakan goroutine untuk memproses setiap RowData
	for _, website := range data {
		go FuncProcessGetTLD(website, ch, errCh)
	}

	var result []RowData
	for i := 0; i < len(data); i++ {
		select {
		case res := <-ch:
			// Memfilter berdasarkan TLD
			if res.TLD == TLD {
				result = append(result, res)
			}
		case err := <-errCh:
			// Mengembalikan error jika ada
			return nil, err
		}
	}

	return result, nil
}

// gunakan untuk melakukan debugging
func main() {
	rows, err := FilterAndFillData(".com", []RowData{
		{1, "google.com", "", "", true, 100},
		{2, "facebook.com", "", "", true, 100},
		{3, "golang.org", "", "", true, 100},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rows)
}
