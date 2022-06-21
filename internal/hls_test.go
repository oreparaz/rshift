package internal

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

const PlaylistUrlTesting = "http://example.com/playlist.m3u8"

func setupMockHttp() {
	httpmock.Activate()
	httpmock.Reset()

	httpmock.RegisterResponder("GET",
		PlaylistUrlTesting,
		httpmock.NewStringResponder(
			200,
			`#EXTM3U
## Created with Golumi Video Platform

#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:5
#EXT-X-MEDIA-SEQUENCE:4421244
#EXTINF:4.992
./GL0/34_2021_12_29_01_31_53.ts
#EXTINF:4.992
./GL0/34_2021_12_29_01_31_58.ts
#EXTINF:5.013
./GL0/34_2021_12_29_01_32_03.ts
#EXTINF:4.992
./GL0/34_2021_12_29_01_32_08.ts`,
		))

	httpmock.RegisterResponder("GET", `http://example.com/./GL0/34_2021_12_29_01_31_53.ts`,
		httpmock.NewStringResponder(200, `some segment body`))
	httpmock.RegisterResponder("GET", `http://example.com/./GL0/34_2021_12_29_01_31_58.ts`,
		httpmock.NewStringResponder(200, `some segment body`))
	httpmock.RegisterResponder("GET", `http://example.com/./GL0/34_2021_12_29_01_32_03.ts`,
		httpmock.NewStringResponder(200, `some segment body`))
	httpmock.RegisterResponder("GET", `http://example.com/./GL0/34_2021_12_29_01_32_08.ts`,
		httpmock.NewStringResponder(200, `some segment body`))
}

var _ = Describe("Parsing a simple M3U file", func() {
	var h Hls
	BeforeEach(func() {
		setupMockHttp()
		BackendFs = afero.NewMemMapFs()

		h = Hls{}
		h.fetchPlaylist(PlaylistUrlTesting)
		h.parseSegments()
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	It("should extract segment URLs", func() {
		Expect(h.segmentUrls[4421244]).To(Equal("http://example.com/./GL0/34_2021_12_29_01_31_53.ts"))
		Expect(h.segmentUrls[4421245]).To(Equal("http://example.com/./GL0/34_2021_12_29_01_31_58.ts"))
		Expect(h.segmentUrls[4421246]).To(Equal("http://example.com/./GL0/34_2021_12_29_01_32_03.ts"))
		Expect(h.segmentUrls[4421247]).To(Equal("http://example.com/./GL0/34_2021_12_29_01_32_08.ts"))
		Expect(len(h.segmentUrls)).To(Equal(4))
	})

	It("should download each segment URL", func() {
		var fs FileStorage
		h.fetchAndSaveSegments(&fs)

		httpCalls := httpmock.GetCallCountInfo()
		Expect(httpCalls["GET http://example.com/./GL0/34_2021_12_29_01_31_53.ts"]).To(Equal(1))
		Expect(httpCalls["GET http://example.com/./GL0/34_2021_12_29_01_31_58.ts"]).To(Equal(1))
		Expect(httpCalls["GET http://example.com/./GL0/34_2021_12_29_01_32_03.ts"]).To(Equal(1))
		Expect(httpCalls["GET http://example.com/./GL0/34_2021_12_29_01_32_08.ts"]).To(Equal(1))
		Expect(len(httpCalls)).To(Equal(5))
	})

	It("should download everything", func() {
		var fs FileStorage
		h.fetchAndSaveAll(&fs)

		httpCalls := httpmock.GetCallCountInfo()
		Expect(httpCalls["GET http://example.com/./GL0/34_2021_12_29_01_31_53.ts"]).To(Equal(1))
		Expect(httpCalls["GET http://example.com/./GL0/34_2021_12_29_01_31_58.ts"]).To(Equal(1))
		Expect(httpCalls["GET http://example.com/./GL0/34_2021_12_29_01_32_03.ts"]).To(Equal(1))
		Expect(httpCalls["GET http://example.com/./GL0/34_2021_12_29_01_32_08.ts"]).To(Equal(1))
		Expect(len(httpCalls)).To(Equal(5))
	})
})
