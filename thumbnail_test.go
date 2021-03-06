package dandler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/sebdah/goldie"

	"github.com/stretchr/testify/assert"
)

func init() {
	goldie.FixtureDir = "testdata/fixtures"
}

func TestLoadThumbnail(t *testing.T) {
	testData := []struct {
		imageName string
		size      int64
	}{
		{
			imageName: "accidentally_save_file.gif",
			size:      17861,
		}, {
			imageName: "blocked_us.png",
			size:      44937,
		}, {
			imageName: "carlton_pls.jpg",
			size:      22806,
		}, {
			imageName: "lemur_pudding_cups.jpg",
			size:      72852,
		}, {
			imageName: "spooning_a_barret.png",
			size:      47306,
		}, {
			imageName: "whats_in_the_case.gif",
			size:      48763,
		},
	}

	tempdir, err := ioutil.TempDir("", "sp9k1-")
	if err != nil {
		t.Fatalf("failed creating test directory: %s", err)
	}

	h := thumbnailHandler{x: 200, y: 200, raw: "testdata", thumbExt: "png", thumbs: tempdir}

	for id, test := range testData {
		h.loadThumbnail(test.imageName)
		info, err := os.Stat(h.generateThumbPath(test.imageName))
		if err != nil {
			t.Logf("#%d - failed to stat thumbnail [%s] tempdir [%s]: %s",
				id, test.imageName, tempdir, err)
			t.Fail()
			continue
		}
		assert.Equal(t, test.size, info.Size(),
			"#%d - size does not match - tempDir [%s]", id, tempdir)
	}
}

func TestThumbnail(t *testing.T) {
	var testData = []struct {
		uri         string
		code        int
		contentType string
	}{
		{
			uri:         "/accidentally_save_file.gif.png",
			code:        200,
			contentType: "image/png",
		}, {
			uri:         "/blocked_us.png.png",
			code:        200,
			contentType: "image/png",
		}, {
			uri:         "/carlton_pls.jpg.png",
			code:        200,
			contentType: "image/png",
		}, {
			uri:         "/lemur_pudding_cups.jpg.png",
			code:        200,
			contentType: "image/png",
		}, {
			uri:         "/spooning_a_barret.png.png",
			code:        200,
			contentType: "image/png",
		}, {
			uri:         "/whats_in_the_case.gif.png",
			code:        200,
			contentType: "image/png",
		}, {
			uri:         "/bad.target.png",
			code:        404,
			contentType: "",
		}, {
			uri:         "/accidentally_save_file.gif.png",
			code:        200,
			contentType: "image/png",
		}, {
			uri:         "/blocked_us.png.png",
			code:        200,
			contentType: "image/png",
		}, {
			uri:         "/carlton_pls.jpg.png",
			code:        200,
			contentType: "image/png",
		},
	}

	tempdir, err := ioutil.TempDir("", "sp9k1-")
	if err != nil {
		t.Fatalf("failed creating test directory: %s", err)
	}

	logger := log.New(ioutil.Discard, "", 0)
	ts := httptest.NewServer(Thumbnail(logger, 300, 250, "./testdata/", tempdir, "png"))
	defer ts.Close()

	baseURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("failed to parse url: %s", err)
	}

	for testID, test := range testData {
		t.Run(fmt.Sprintf("TestThumbnail-%d", testID), func(t *testing.T) {
			uri, err := url.Parse(test.uri)
			if err != nil {
				t.Errorf("bad URI path: [%s]", test.uri)
				return
			}

			res, err := http.Get(baseURL.ResolveReference(uri).String())
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, test.code, res.StatusCode, "status code does not match: ")
			if test.code != 200 {
				if res.StatusCode != test.code {
					t.Logf("the response returned: \n%#v\n", res)
				}
				return
			}
			assert.Equal(t, test.contentType, res.Header.Get("Content-Type"), "Content-Type does not match: ")

			// body, err := ioutil.ReadAll(res.Body)
			// res.Body.Close()
			// require.NoError(t, err)
			// goldie.Assert(t, t.Name(), body)
		})
	}
}

func TestThumbnailJPG(t *testing.T) {
	var testData = []struct {
		uri           string
		code          int
		md5           string
		contentLength int64
		contentType   string
	}{
		{
			uri:           "/accidentally_save_file.gif",
			code:          200,
			md5:           "2aa9ba78ec27dc96a3f5603e9e8eb646",
			contentLength: 12489,
			contentType:   "image/",
		}, {
			uri:           "/blocked_us.png",
			code:          200,
			md5:           "2fc5189bea70182964bf9126bcb3f0be",
			contentLength: 10887,
			contentType:   "image/",
		}, {
			uri:           "/carlton_pls.jpg",
			code:          200,
			md5:           "950e11dcdbbe9e27781aed1e815ff83f",
			contentLength: 5081,
			contentType:   "image/",
		}, {
			uri:           "/lemur_pudding_cups.jpg",
			code:          200,
			md5:           "b5a688f25e0c248a6b101467957fc989",
			contentLength: 17019,
			contentType:   "image/",
		}, {
			uri:           "/spooning_a_barret.png",
			code:          200,
			md5:           "b62b31ec6cfc5fd85dec71a3592373a8",
			contentLength: 10705,
			contentType:   "image/",
		}, {
			uri:           "/whats_in_the_case.gif",
			code:          200,
			md5:           "806a2539113d46547dbc0fe779e5c4f3",
			contentLength: 7574,
			contentType:   "image/",
		}, {
			uri:           "/bad.target.png",
			code:          404,
			md5:           "",
			contentLength: 0,
			contentType:   "",
		}, {
			uri:           "/accidentally_save_file.gif",
			code:          200,
			md5:           "2aa9ba78ec27dc96a3f5603e9e8eb646",
			contentLength: 12489,
			contentType:   "image/",
		}, {
			uri:           "/blocked_us.png",
			code:          200,
			md5:           "2fc5189bea70182964bf9126bcb3f0be",
			contentLength: 10887,
			contentType:   "image/",
		}, {
			uri:           "/carlton_pls.jpg",
			code:          200,
			md5:           "950e11dcdbbe9e27781aed1e815ff83f",
			contentLength: 5081,
			contentType:   "image/",
		},
	}

	for _, ext := range []string{"jpg", "jpeg"} {
		tempdir, err := ioutil.TempDir("", "sp9k1-")
		if err != nil {
			t.Fatalf("failed creating test directory: %s", err)
		}

		logger := log.New(ioutil.Discard, "", 0)
		ts := httptest.NewServer(Thumbnail(logger, 300, 250, "./testdata/", tempdir, ext))
		defer ts.Close()

		baseURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("failed to parse url: %s", err)
		}

		for testID, test := range testData {
			t.Run(fmt.Sprintf("TestThumbnail-%s-%d-", ext, testID), func(t *testing.T) {
				uri, err := url.Parse(test.uri + "." + ext)
				if err != nil {
					t.Errorf("bad URI path: [%s]", test.uri)
					return
				}

				res, err := http.Get(baseURL.ResolveReference(uri).String())
				if err != nil {
					t.Error(err)
					return
				}

				assert.Equal(t, test.code, res.StatusCode, "status code does not match: ")
				if test.code != 200 {
					if res.StatusCode != test.code {
						t.Logf("the response returned: \n%#v\n", res)
					}
					return
				}
				assert.Equal(t, test.contentType+ext, res.Header.Get("Content-Type"), "Content-Type does not match: ")

				// body, err := ioutil.ReadAll(res.Body)
				// res.Body.Close()
				// require.NoError(t, err)
				// goldie.Assert(t, t.Name(), body)
			})
		}
	}
}

func TestGeneratePaths(t *testing.T) {
	testData := []struct {
		imageName string
		rawPath   string
		thumbPath string
	}{
		{
			imageName: "accidentally_save_file.gif",
			rawPath:   "testdata/accidentally_save_file.gif",
			thumbPath: "output/accidentally_save_file.gif.jpg",
		}, {
			imageName: "blocked_us.png",
			rawPath:   "testdata/blocked_us.png",
			thumbPath: "output/blocked_us.png.jpg",
		}, {
			imageName: "carlton_pls.jpg",
			rawPath:   "testdata/carlton_pls.jpg",
			thumbPath: "output/carlton_pls.jpg.jpg",
		}, {
			imageName: "lemur_pudding_cups.jpg",
			rawPath:   "testdata/lemur_pudding_cups.jpg",
			thumbPath: "output/lemur_pudding_cups.jpg.jpg",
		}, {
			imageName: "spooning_a_barret.png",
			rawPath:   "testdata/spooning_a_barret.png",
			thumbPath: "output/spooning_a_barret.png.jpg",
		}, {
			imageName: "whats_in_the_case.gif",
			rawPath:   "testdata/whats_in_the_case.gif",
			thumbPath: "output/whats_in_the_case.gif.jpg",
		},
	}

	h := thumbnailHandler{raw: "testdata", thumbs: "output", thumbExt: "jpg"}

	for id, test := range testData {
		assert.Equal(t, test.rawPath, h.generateRawPath(test.imageName), "#%d - wrong raw path", id)
		assert.Equal(t, test.thumbPath, h.generateThumbPath(test.imageName), "#%d - wrong thumb path", id)
	}
}
