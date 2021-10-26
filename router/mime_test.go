package router

import "testing"

func TestMimeExtension(t *testing.T) {
	getter := &extensionMimeGetterType{}

	check := func(path string, expectedMime string) {
		mime := getter.GetMime(path)
		if mime != expectedMime {
			t.Errorf("Incorrect mime for file %s. Expected %s got %s", path, expectedMime, mime)
		}
	}

	check("index.html", "text/html")
	check("document.pdf", "application/pdf")
	check("image.jpg", "image/jpeg")
	check("foo", "application/octet-stream")
}
