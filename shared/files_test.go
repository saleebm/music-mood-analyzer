package shared

import (
	"log"
	"testing"
)

func TestTestServer(t *testing.T) {
	t.Run("file exists", func(t *testing.T) {
		filename := "/Users/minasaleeb/go/saleebm/music-mood-analyzer/tmp/temp-46502.json"
		log.Println(Exists(filename))
	})
}
