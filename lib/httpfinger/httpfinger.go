package httpfinger

import (
	"encoding/json"
)

//r["FaviconHash"] , r["KeywordFinger"]
func Init() map[string]int {
	_ = json.Unmarshal(faviconHashByte, &FaviconHash)
	_ = json.Unmarshal(keywordFingerByte, &KeywordFinger)

	r := make(map[string]int)
	r["FaviconHash"] = len(FaviconHash)
	r["KeywordFinger"] = len(KeywordFinger)
	return r
}
