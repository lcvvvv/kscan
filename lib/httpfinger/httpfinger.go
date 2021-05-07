package httpfinger

import (
	"encoding/json"
)

//r["FaviconHash"] , r["KeywordFinger"]
func Init() map[string]int {
	_ = json.Unmarshal(faviconHashByte, &FaviconHash)
	_ = json.Unmarshal(keywordFingerSourceByte, &KeywordFinger)
	var keywordFingerFofa keywordFinger
	_ = json.Unmarshal(keywordFingerFofaByte, &keywordFingerFofa)
	KeywordFinger = append(KeywordFinger, keywordFingerFofa...)

	r := make(map[string]int)
	r["FaviconHash"] = len(FaviconHash)
	r["KeywordFinger"] = len(KeywordFinger)
	return r
}
