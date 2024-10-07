package dbg

import (
	"net/url"
	"strings"
)

var (
	fail_url_tags = []string{"//", ".."}
)

func CutUrlSuffixSlash(url_text string) (string, error) {
	url_text = strings.TrimSpace(url_text)
	if _, err := url.Parse(url_text); err != nil {
		return "", err
	}
	url_text = strings.TrimSuffix(url_text, "/")

	v := strings.Replace(url_text, "://", "", 1)
	//fmt.Println(v)
	for _, tag := range fail_url_tags {
		if strings.Contains(v, tag) {
			return "", Error("(", url_text, ") contains :", tag)
		}
	}

	return url_text, nil
}

func CutUrlSuffixSlashP(url_text_ptr *string) error {
	url_text, err := CutUrlSuffixSlash(*url_text_ptr)
	if err != nil {
		return err
	}
	*url_text_ptr = url_text
	return nil
}
