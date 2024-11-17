package banner

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

func ShowBanner() {
	ban, readErr := os.ReadFile("banner.txt")
	if readErr != nil {
		zap.L().Warn("gohttpd: banner.txt cannot read", zap.String("banner file read", readErr.Error()))
		return
	}
	fmt.Println(string(ban) + "\n")
}
