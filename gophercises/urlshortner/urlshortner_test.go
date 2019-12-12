package urlshortner

import (
	"fmt"
	"testing"
)

func TestParseYaml(t *testing.T) {

	yml := `
- path: /go-doc
  url: https://godoc.org/
- path: /go-testing
  url: https://godoc.org/testing
`
	pathURLs, err := ParseYaml([]byte(yml))
	if err != nil {
		t.Error(err)
	}
	t.Errorf("%v", pathURLs)
	// fmt.Println(pathURLs)
}

func TestBuildMap(t *testing.T) {
	pathurls := `[{/go-doc https://godoc.org/} {/go-testing https://godoc.org/testing}]`
	pathsToUrls := BuildMap(pathurls)
	fmt.Println(pathsToUrls)
}
