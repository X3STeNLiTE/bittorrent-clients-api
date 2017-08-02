package utorrent

import (
	"encoding/json"
)

//LabelFieldIndex ...
type LabelFieldIndex int

//
const (
	LabelIndexName  LabelFieldIndex = iota // 0: Name
	LabelIndexCount                        // 1: Count
)

//Label ...
type Label struct {
	Name  string
	Count int
}

type listFileRespContainer struct {
	Build      int
	Label      [][]json.RawMessage
	Torrents   [][]json.RawMessage
	Torrentc   string
	Rssfeeds   []string
	Rssfilters []string
	Message    []string
}

//ListFileResponse ...
type ListFileResponse struct {
	Build      int       `json:"build"`    //Utorrent Client Build #
	Label      []Label   `json:"label"`    //Labels
	Torrents   []Torrent `json:"torrents"` //Torrents
	Torrentc   string    `json:"torrentc"` //Cache ID
	Rssfeeds   []string  `json:"rssfeeds"`
	Rssfilters []string  `json:"rssfilters"`
	Messages   []string  `json:"messages"`
}

func unmarshalListFileResponse(data []byte) (*ListFileResponse, error) {
	var rsd ListFileResponse
	var container listFileRespContainer
	err := json.Unmarshal(data, &container)
	if err != nil {
		return nil, err
	}

	rsd.Build = container.Build
	rsd.Torrentc = container.Torrentc
	rsd.Torrents, _ = unmarshalTorrentList(container.Torrents)
	rsd.Label, _ = unmarshalLabelList(container.Label)
	rsd.Rssfeeds = container.Rssfeeds
	rsd.Rssfilters = container.Rssfilters

	return &rsd, nil
}

func unmarshalLabelList(data [][]json.RawMessage) ([]Label, error) {
	labels := make([]Label, 0)
	for _, item := range data {
		label := Label{}
		for i, elem := range item {
			var inf interface{}
			switch LabelFieldIndex(i) {
			case LabelIndexName:
				inf = &label.Name
				break
			case LabelIndexCount:
				inf = &label.Count
				break
			}
			json.Unmarshal(elem, inf)
		}
		labels = append(labels, label)
	}
	return labels, nil
}

func unmarshalTorrentList(data [][]json.RawMessage) ([]Torrent, error) {
	torrents := make([]Torrent, 0)
	for _, item := range data {
		torrent := Torrent{}
		for i, elem := range item {
			var inf interface{}
			switch TorrentFieldIndex(i) {
			case TorrentsIndexHash:
				inf = &torrent.Hash
				break
			case TorrentsIndexStatus:
				inf = &torrent.Status
				break
			case TorrentsIndexName:
				inf = &torrent.Name
				break
			case TorrentsIndexSize:
				inf = &torrent.Size
				break
			case TorrentsIndexProgress:
				inf = &torrent.Progress
				break
			case TorrentsIndexDownloaded:
				inf = &torrent.Downloaded
				break
			case TorrentsIndexUploaded:
				inf = &torrent.Uploaded
				break
			case TorrentsIndexRatio:
				inf = &torrent.Ratio
				break
			case TorrentsIndexUploadSpeed:
				inf = &torrent.UlSpeed
				break
			case TorrentsIndexDownloadSpeed:
				inf = &torrent.DlSpeed
				break
			case TorrentsIndexETA:
				inf = &torrent.Eta
				break
			case TorrentsIndexLabel:
				inf = &torrent.Label
				break
			case TorrentsIndexPeerConnected:
				inf = &torrent.PeerConnected
				break
			case TorrentsIndexPeersInSwarm:
				inf = &torrent.PeerInSwarm
				break
			case TorrentsIndexSeedsConnected:
				inf = &torrent.SeedConnected
				break
			case TorrentsIndexSeedsInSwarm:
				inf = &torrent.SeedInSwarm
				break
			case TorrentsIndexAvailability:
				inf = &torrent.Availability
				break
			case TorrentsIndexQueueOrder:
				inf = &torrent.TorrentQueueOrder
				break
			case TorrentsIndexRemaining:
				inf = &torrent.Remaining
				break
			case TorrentsIndexDateAdded:
				inf = &torrent.DateAdded
				break
			case TorrentsIndexDateCompleted:
				inf = &torrent.DateCompleted
				break
			case TorrentsIndexSavePath:
				inf = &torrent.SavePath
				break
			}
			json.Unmarshal(elem, inf)
		}
		torrents = append(torrents, torrent)
	}
	return torrents, nil
}
