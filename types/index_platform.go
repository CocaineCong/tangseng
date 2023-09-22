package types

import (
	"github.com/RoaringBitmap/roaring"
)

type InvertedInfo struct {
	Token  string          `json:"token"`
	DocIds *roaring.Bitmap `json:"doc_ids"`
}
