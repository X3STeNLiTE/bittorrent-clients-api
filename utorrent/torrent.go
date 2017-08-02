package utorrent

//Compatible with uTorrent 3.4.5

import (
	"sort"
)

//UTorrent file status
/*
Status:
	1 = Started
	2 = Checking
	4 = Start after check
	8 = Checked
	16 = Error
	32 = Paused
	64 = Queued
	128 = Loaded
*/
const (
	TorrentStatusStarted         uint8 = 1 << 0
	TorrentStatusChecking        uint8 = 1 << 1
	TorrentStatusStartAfterCheck uint8 = 1 << 2
	TorrentStatusChecked         uint8 = 1 << 3
	TorrentStatusError           uint8 = 1 << 4
	TorrentStatusPaused          uint8 = 1 << 5
	TorrentStatusQueued          uint8 = 1 << 6
	TorrentStatusLoaded          uint8 = 1 << 7
)

//TorrentFieldIndex ...
type TorrentFieldIndex int

//Defined torrents index of each `json: "torrents"` field
const (
	TorrentsIndexHash           TorrentFieldIndex = iota // 0: Hash (string),
	TorrentsIndexStatus                                  // 1: Status* (integer),
	TorrentsIndexName                                    // 2: Name (string),
	TorrentsIndexSize                                    // 3: Size (integer in bytes),
	TorrentsIndexProgress                                // 4: Percent Progress (integer in per mils),
	TorrentsIndexDownloaded                              // 5: Downloaded (integer in bytes),
	TorrentsIndexUploaded                                // 6: Uploaded (integer in bytes),
	TorrentsIndexRatio                                   // 7: Ratio (integer in per mils),
	TorrentsIndexUploadSpeed                             // 8: Upload Speed (integer in bytes per second),
	TorrentsIndexDownloadSpeed                           // 9: Download Speed (integer in bytes per second),
	TorrentsIndexETA                                     // 10: ETA (integer in seconds),
	TorrentsIndexLabel                                   // 11: Label (string),
	TorrentsIndexPeerConnected                           // 12: Peers Connected (integer),
	TorrentsIndexPeersInSwarm                            // 13: Peers In Swarm (integer),
	TorrentsIndexSeedsConnected                          // 14: Seeds Connected (integer),
	TorrentsIndexSeedsInSwarm                            // 15: Seeds In Swarm (integer),
	TorrentsIndexAvailability                            // 16: Availability (integer in 1/65535ths),
	TorrentsIndexQueueOrder                              // 17: Torrent Queue Order (integer),
	TorrentsIndexRemaining                               // 18: Remaining (integer in bytes)
	TorrentsIndexDownloadURL                             // 19: Download URL
	TorrentsIndexRssFeedURL                              // 20: Rss Feed URL
	TorrentsIndexStatusMessage                           // 21: Status Message
	TorrentsIndexStreamID                                // 22: Stream ID
	TorrentsIndexDateAdded                               // 23: Remaining (integer in bytes)
	TorrentsIndexDateCompleted                           // 24: Remaining (integer in bytes)
	TorrentsIndexAppUpdateURL                            // 25: App Update URL
	TorrentsIndexSavePath                                // 26: Remaining (integer in bytes)
)

//Torrent ...
type Torrent struct {
	Hash              string
	Status            uint8
	Name              string
	Size              uint64
	Progress          uint32
	Downloaded        uint64
	Uploaded          uint64
	Ratio             uint64
	UlSpeed           uint32
	DlSpeed           uint32
	Eta               uint32
	Label             string
	PeerConnected     uint32
	PeerInSwarm       uint32
	SeedConnected     uint32
	SeedInSwarm       uint32
	Availability      uint32
	TorrentQueueOrder uint32
	Remaining         uint32
	DownloadURL       string
	StatusMessage     string
	DateAdded         uint64
	DateCompleted     uint64
	SavePath          string
}

//PercentCompleted Torrent.Progress devided by 10..
func (t Torrent) PercentCompleted() float64 {
	return float64(t.Progress) / 10.0
}

//AvailabilityRatio Torrent.Availability devided by 65535
func (t Torrent) AvailabilityRatio() float64 {
	return float64(t.Availability) / 65535.0
}

//Completed return is torrent completed
func (t Torrent) Completed() bool {
	return t.Progress == 1000
}

func (t Torrent) Error() bool {
	return t.Status&TorrentStatusError != 0
}

//StatusMap ...
type StatusMap map[int]string

var status = StatusMap{
	int(TorrentStatusStarted):         "STA",
	int(TorrentStatusChecking):        "CHK",
	int(TorrentStatusStartAfterCheck): "SAC",
	int(TorrentStatusChecked):         "CKD",
	int(TorrentStatusError):           "ERR",
	int(TorrentStatusPaused):          "PAU",
	int(TorrentStatusQueued):          "QUE",
	int(TorrentStatusLoaded):          "LOD",
}

func (m StatusMap) concat(status *string, flag int) {
	v, ok := m[flag]
	sp := ""
	if ok {
		sp = ","
	}
	if len(*status) > 0 {
		*status += sp
	}
	*status += v
	return
}

//StatusStr ...
func (t Torrent) StatusStr() string {
	r := ""
	//
	status.concat(&r, int(t.Status&TorrentStatusStarted))
	status.concat(&r, int(t.Status&TorrentStatusChecking))
	status.concat(&r, int(t.Status&TorrentStatusStartAfterCheck))
	status.concat(&r, int(t.Status&TorrentStatusChecked))
	status.concat(&r, int(t.Status&TorrentStatusError))
	status.concat(&r, int(t.Status&TorrentStatusPaused))
	status.concat(&r, int(t.Status&TorrentStatusQueued))
	status.concat(&r, int(t.Status&TorrentStatusLoaded))
	//
	return r
}

//TorrentSortAsc ...
type TorrentSortAsc func(a, b *Torrent) bool

type torrentCollectionSorter struct {
	torrents []Torrent
	by       TorrentSortAsc
}

//Len ...
func (t torrentCollectionSorter) Len() int {
	return len(t.torrents)
}

//Less ...
func (t torrentCollectionSorter) Less(i, j int) bool {
	return t.by(&t.torrents[i], &t.torrents[j])
}

//Swap ...
func (t torrentCollectionSorter) Swap(i, j int) {
	t.torrents[i], t.torrents[j] = t.torrents[j], t.torrents[i]
}

//ByName sort []Torrent by Name
var ByName = TorrentSortAsc(func(a, b *Torrent) bool {
	return a.Name < b.Name
})

//ByProgress sort []Torrent by progress
var ByProgress = TorrentSortAsc(func(a, b *Torrent) bool {
	return a.Progress < b.Progress
})

//ByQueueOrder sort []Torrent by Queue Order #
var ByQueueOrder = TorrentSortAsc(func(a, b *Torrent) bool {
	return a.TorrentQueueOrder < b.TorrentQueueOrder
})

//ByAdded sort []Torrent by Added Date
var ByAdded = TorrentSortAsc(func(a, b *Torrent) bool {
	return a.DateAdded < b.DateAdded
})

//ByFinished sort []Torrent by Finished Date
var ByFinished = TorrentSortAsc(func(a, b *Torrent) bool {
	return a.DateCompleted < b.DateCompleted
})

//Sort sort []Torrent by provided sorter
func (by TorrentSortAsc) sort(torrents []Torrent) {
	sorter := &torrentCollectionSorter{
		torrents: torrents,
		by:       by,
	}
	sort.Sort(sorter)
}

//TorrentCollection slice of Torrent
type TorrentCollection []Torrent

//Sort sort []Torrent by Name
func (t TorrentCollection) Sort(by TorrentSortAsc, desc bool) {
	by.sort(t)
}
