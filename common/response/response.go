package response

import (
	"fmt"
	"log"
	"net/http"
)

func CheckErrNonFatal(err error) error {
	if err != nil {
		log.Printf("❌ Error encountered: %v", err)
		return err
	}
	return nil
}

func CheckCodeNonFatal(res *http.Response) error {
	if res.StatusCode != 200 {
		if res.StatusCode == 403 {
			errMsg := fmt.Sprintf("🚫 403 Forbidden error for URL: %s", res.Request.URL.String())
			log.Printf("%s", errMsg)
			return fmt.Errorf("403 Forbidden: %s", res.Request.URL.String())
		} else {
			errMsg := fmt.Sprintf("❌ HTTP %d error for URL: %s", res.StatusCode, res.Request.URL.String())
			log.Printf("%s", errMsg)
			return fmt.Errorf("HTTP %d: %s", res.StatusCode, res.Request.URL.String())
		}
	}
	return nil
}
