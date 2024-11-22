package node

import (
	"fmt"
)

func manifestURLFromCID(cid string) string {
	return fmt.Sprintf("https://%s.ipfs.w3s.link/manifest.json", cid)
}
