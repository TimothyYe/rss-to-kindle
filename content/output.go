package content

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"rss-to-kindle/utils"
	"strings"

	"github.com/nfnt/resize"
)

func createTempDir() string {
	dir, err := ioutil.TempDir("", "example")
	utils.ExitIfErr(err)
	return dir
}

func makeArticlesFile(dir string, feed Feed) {
	for _, section := range feed.Sections {
		for _, article := range section.Articles {
			utils.WriteFile(dir, fmt.Sprintf("%d.html", article.ID), []byte(GenerateArticle(article)))
		}
	}
}

func makeContentsFile(dir string, feed Feed) {
	utils.WriteFile(dir, "contents.html", []byte(GenerateContents(feed)))
}

func resizeImageFile(path string) {
	if regexp.MustCompile(`.jpg$|.png$`).MatchString(path) == false {
		return
	}

	file, err := os.Open(path)
	utils.ExitIfErr(err)

	var img image.Image
	switch filepath.Ext(path) {
	case ".jpg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	}
	utils.ExitIfErr(err)

	file.Close()
	os.Remove(path)

	m := resize.Resize(640, 0, img, resize.Bicubic)

	out, err := os.Create(strings.Replace(path, ".png", ".jpg", -1))
	utils.ExitIfErr(err)
	defer out.Close()

	err = jpeg.Encode(out, m, nil)
	utils.ExitIfErr(err)
}

func makeImageFile(dir string, images []string) {
	for _, image := range images {
		res, err := http.Get(image)
		utils.ExitIfErr(err)

		content, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		utils.ExitIfErr(err)

		filename := utils.NormalizeImageFilename(utils.GetFilename(image, true))
		utils.WriteFile(dir, filename, content)

		resizeImageFile(filepath.Join(dir, filename))
	}
}

func makeOpfFile(dir string, feed Feed) []string {
	manifest, images := GenerateManifest(feed)
	spine := GenerateSpine(feed)

	makeImageFile(dir, images)

	utils.WriteFile(dir, "main.opf", []byte(GenerateOpf(feed, manifest, spine)))
	return images
}

func makeNavContentsFile(dir string, feed Feed) {
	utils.WriteFile(dir, "nav-contents.ncx", []byte(GenerateNavMain(feed)))
}

//MakeMain ...
func MakeMain(feed Feed) string {
	dir := "./output"

	makeContentsFile(dir, feed)
	makeOpfFile(dir, feed)
	makeArticlesFile(dir, feed)
	makeNavContentsFile(dir, feed)

	return dir
}
