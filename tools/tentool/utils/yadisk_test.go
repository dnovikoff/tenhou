package utils

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYaDiskLink(t *testing.T) {
	assert.Equal(t, "https://cloud-api.yandex.net/v1/disk/public/resources?public_key=https%3A%2F%2Fyadi.sk%2Fd%2FD_0tO_135zIKPw",
		YaDiskResourcesLink("https://yadi.sk/d/D_0tO_135zIKPw"))
}

func TestYaDiskParse(t *testing.T) {
	input := loadTestData(t, "yadisk.json")
	items, err := YaDiskParseItems(input)
	require.NoError(t, err)
	// for _, v := range items {
	// 	fmt.Printf(`{
	// 		Type: "%v",
	// 		File: "%v",
	// 		Path: "%v",
	// 		PublicURL: "%v",
	// 	},`, v.Type, v.File, v.Path, v.PublicURL)
	// }
	assert.Equal(t,
		[]YADiskItem{
			{
				Type:      "file",
				File:      "https://downloader.disk.yandex.ru/disk/c6c13546d1fdaaa28ea1bd5f4c260672589e00243f759ef9309ad6e029b70da6/5bc36f10/1_wL5UdiKftVnHai1l69oC0JxPIk-JA6XrFfGNPWIgtxi-ur0zC9Nxm4ag_5HeFCyCemlLI8ke6js4PAlP0_3Q%3D%3D?uid=0&filename=2009.zip&disposition=attachment&hash=&limit=0&content_type=application%2Fx-zip-compressed&fsize=271490715&hid=55eb069447de1d24113a70dbdad444a5&media_type=compressed&tknv=v2&etag=3d383b391b5dee2b5d280d063f10d841",
				Path:      "/2009.zip",
				PublicURL: "https://yadi.sk/d/LGFi6N4CWJ2drg",
			}, {
				Type:      "file",
				File:      "https://downloader.disk.yandex.ru/disk/25515cb82b9ee8a6d9719e1ad7b965cc6e8f40783105b92a44a15324647ca3b5/5bc36f10/1_wL5UdiKftVnHai1l69oKN-yk_7ZUJvr4-PNp1nV17BxrI-4M4F-xnPr7isgxk6itifHYttfcVYvPr1xmEkKQ%3D%3D?uid=0&filename=2010.zip&disposition=attachment&hash=&limit=0&content_type=application%2Fx-zip-compressed&fsize=962886929&hid=1ff742d5a379dd09fa99b986c45e84bb&media_type=compressed&tknv=v2&etag=af84a2c25ba6ec35c2cabc3fb4f74bed",
				Path:      "/2010.zip",
				PublicURL: "https://yadi.sk/d/D_0tO_135zIKPw/2010.zip",
			}, {
				Type:      "file",
				File:      "https://downloader.disk.yandex.ru/disk/663d882f7640cfa8eb869743003df068a14f92c897f41750087e4cee7f1c0e62/5bc36f10/1_wL5UdiKftVnHai1l69oA6FTxzrwyvZCKzy05HCgIZwUb8Ovis8jbDfYSWjddSJyIb7oHUYZj6lsQ9BHJ-Qbg%3D%3D?uid=0&filename=2011.zip&disposition=attachment&hash=&limit=0&content_type=application%2Fx-zip-compressed&fsize=1482133817&hid=8508fa0cd8d69cb4bef7a72f0f0241a5&media_type=compressed&tknv=v2&etag=bd05f20da6757cd58195123378676fc9",
				Path:      "/2011.zip",
				PublicURL: "https://yadi.sk/d/D_0tO_135zIKPw/2011.zip",
			}, {
				Type:      "file",
				File:      "https://downloader.disk.yandex.ru/disk/92701c385a8196a17d87f1653ecfef447ccc05b8dc6c48eddfe3799e280db408/5bc36f10/1_wL5UdiKftVnHai1l69oJtj3nXHfs7aKfSiScpy_Y0Jx9EsG9Vb-Qco4Ael8Z9YAL71wdRzja3DG-KKY0MbDw%3D%3D?uid=0&filename=2012.zip&disposition=attachment&hash=&limit=0&content_type=application%2Fx-zip-compressed&fsize=1785179658&hid=83d617d4f30521d6ee01545eebd68f86&media_type=compressed&tknv=v2&etag=c5353fe089528931eefa79cbc806bcd7",
				Path:      "/2012.zip",
				PublicURL: "https://yadi.sk/d/D_0tO_135zIKPw/2012.zip",
			}, {
				Type:      "file",
				File:      "https://downloader.disk.yandex.ru/disk/cb236a44e05d9dda9ab446846ca5142f585aa37d975b8d337f4ebb6b5242d549/5bc36f10/1_wL5UdiKftVnHai1l69oJpx6NlB9vi12j_D8qX8zB9c7H85IABPhBONUeHezZzTx8bXv9oNo28Vr3ZQSW7ipw%3D%3D?uid=0&filename=2013.zip&disposition=attachment&hash=&limit=0&content_type=application%2Fx-zip-compressed&fsize=1934631204&hid=a40ca2a5340a7c8223fbe933b49795f1&media_type=compressed&tknv=v2&etag=d579ad0d664a40fc884b7d2bb9107206",
				Path:      "/2013.zip",
				PublicURL: "https://yadi.sk/d/D_0tO_135zIKPw/2013.zip",
			}, {
				Type:      "file",
				File:      "https://downloader.disk.yandex.ru/disk/e3b38a909bb1fb0ce5cdab8442fa0ec9ab2a82946534031308dfd512d4d39329/5bc36f10/1_wL5UdiKftVnHai1l69oHLH3QklMAhQH5pWtllefEBl2sDWd3GMMVxXC0gGJx7mkYZrqBW12-C4dcvXKd4ZiA%3D%3D?uid=0&filename=2014.zip&disposition=attachment&hash=&limit=0&content_type=application%2Fx-zip-compressed&fsize=1841563615&hid=a1f4483abda0082a67b6a25d64750bf9&media_type=compressed&tknv=v2&etag=6eaa33f6b2e1e9b02798f328b1034886",
				Path:      "/2014.zip",
				PublicURL: "https://yadi.sk/d/D_0tO_135zIKPw/2014.zip",
			}, {
				Type:      "file",
				File:      "https://downloader.disk.yandex.ru/disk/ea34a2cc0c8a7160350b07f42a39ee5af177f5fec3848129c5cb7d2fe19e73a6/5bc36f10/1_wL5UdiKftVnHai1l69oLARlonB3wlVGaicQPRPAgfV2VFYviFHWcT7P7TuVnm2wcnYp7u6HAh8rosaXFuSSw%3D%3D?uid=0&filename=2015.zip&disposition=attachment&hash=&limit=0&content_type=application%2Fx-zip-compressed&fsize=1828381069&hid=896128d146191cef34bbb3ebfbfccfbd&media_type=compressed&tknv=v2&etag=577a247106d66ed9a04051a6c2dc4fb1",
				Path:      "/2015.zip",
				PublicURL: "https://yadi.sk/d/D_0tO_135zIKPw/2015.zip",
			}}, items)
}

func loadTestData(t require.TestingT, filename string) string {
	data, err := ioutil.ReadFile("test_data/" + filename)
	require.NoError(t, err)
	return string(data)
}
