package internal

import (
	"bytes"
	"github.com/grafov/m3u8"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"path/filepath"
)

func getMp(playlistBytes []byte) *m3u8.MediaPlaylist {
	mp := m3u8.MediaPlaylist{}
	err := mp.Decode(*bytes.NewBuffer(playlistBytes), false)
	Expect(err).Should(BeNil())
	return &mp
}

func readPlaylistAndCheck(retrieve int, expectedSeqNo int, fs FileStorage) {
	playlistBytes, err := fs.ReadPlaylistNear(retrieve)
	Expect(err).Should(BeNil())
	Expect(int(getMp(playlistBytes).SeqNo)).Should(Equal(expectedSeqNo))
}

var _ = Describe("Saving a simple playlist", func() {
	BeforeEach(func() {
		BackendFs = afero.NewMemMapFs()
	})
	It("saves the index file", func() {
		var fs FileStorage
		fs.SavePlaylist(&Hls{
			Mp: &m3u8.MediaPlaylist{
				SeqNo: 123,
				Segments: []*m3u8.MediaSegment{},
			},
			downloadTime: 456,
			playlistUrl: PlaylistUrlTesting,
		})
		_, err := BackendFs.Stat(filepath.Join(OutPath, "m3u8", "playlist.m3u8---456"))
		Expect(err).Should(BeNil())
	})
	It("saves a segment file", func() {
		var fs FileStorage
		err := fs.SaveSegment(
			&Hls{
				Mp: &m3u8.MediaPlaylist{
					SeqNo: 123,
					Segments: []*m3u8.MediaSegment{},
				},
				downloadTime: 456,
				playlistUrl: PlaylistUrlTesting,
			},
		"http://foo/segment_001.ts",
			[]byte{0x00, 0x01},
		)
		Expect(err).Should(BeNil())
		_, err = BackendFs.Stat(filepath.Join(OutPath, "ts", "segment_001.ts"))
		Expect(err).Should(BeNil())
	})

	It("retrieves the right playlist", func() {
		var fs FileStorage
		var err error

		for _, i := range []int {10, 20, 30} {
			err = fs.SavePlaylist(
				&Hls{
					Mp: &m3u8.MediaPlaylist{
						SeqNo: uint64(100 + i),
						Segments: []*m3u8.MediaSegment{},
					},
					downloadTime: int64(i),
					playlistUrl: PlaylistUrlTesting,
				})
			Expect(err).Should(BeNil())
		}

		readPlaylistAndCheck(9, 110, fs)
		readPlaylistAndCheck(10, 120, fs)
		readPlaylistAndCheck(11, 120, fs)
		readPlaylistAndCheck(21, 130, fs)
		readPlaylistAndCheck(16, 120, fs)
	})

	It("retrieves the right playlist after inserting", func() {
		var fs FileStorage
		var err error

		for _, i := range []int {10, 20, 30} {
			err = fs.SavePlaylist(
				&Hls{
					Mp: &m3u8.MediaPlaylist{
						SeqNo: uint64(100 + i),
						Segments: []*m3u8.MediaSegment{},
					},
					downloadTime: int64(i),
					playlistUrl: PlaylistUrlTesting,
				})
			Expect(err).Should(BeNil())
		}

		readPlaylistAndCheck(9, 110, fs)
		readPlaylistAndCheck(10, 120, fs)
		readPlaylistAndCheck(11, 120, fs)
		readPlaylistAndCheck(21, 130, fs)
		readPlaylistAndCheck(16, 120, fs)
	})

	It("retrieves the right playlist after inserting with colliding sequence ID", func() {
		var fs FileStorage
		var err error

		downloadTimes :=     []int64{10, 20, 30, 40, 50}
		for index, i := range []int {10, 20, 30, 10, 20} {
			err = fs.SavePlaylist(
				&Hls{
					Mp: &m3u8.MediaPlaylist{
						SeqNo: uint64(i),
						Segments: []*m3u8.MediaSegment{},
					},
					downloadTime: downloadTimes[index],
					playlistUrl: PlaylistUrlTesting,
				})
			Expect(err).Should(BeNil())
		}


		readPlaylistAndCheck(9, 10, fs)

		readPlaylistAndCheck(10, 20, fs)
		readPlaylistAndCheck(11, 20, fs)

		readPlaylistAndCheck(19, 20, fs)
		readPlaylistAndCheck(20, 30, fs)
		readPlaylistAndCheck(21, 30, fs)
		readPlaylistAndCheck(22, 30, fs)

		readPlaylistAndCheck(30, 10, fs)
		readPlaylistAndCheck(39, 10, fs)

		readPlaylistAndCheck(40, 20, fs)
		readPlaylistAndCheck(42, 20, fs)
	})
})
