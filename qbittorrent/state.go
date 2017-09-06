package qbittorrent

//State ...
type State map[string]string

//
const (
	StateError       = "Error"                           //Some error occurred, applies to paused torrents
	StatePausedUP    = "Upload Paused"                   //Torrent is paused and has finished downloading
	StatePausedDL    = "Download Paused"                 //Torrent is paused and has NOT finished downloading
	StateQueuedUP    = "Upload Queued"                   //Queuing is enabled and torrent is queued for upload
	StateQueuedDL    = "Download Queued"                 //Queuing is enabled and torrent is queued for download
	StateUploading   = "Uploading"                       //Torrent is being seeded and data is being transfered
	StateStalledUP   = "Uploading, connection waiting"   //Torrent is being seeded, but no connection were made
	StateCheckingUP  = "Upload Checking"                 //Torrent has finished downloading and is being checked; this status also applies to preallocation (if enabled) and checking resume data on qBt startup
	StateCheckingDL  = "Download File Checking"          //Same as checkingUP, but torrent has NOT finished downloading
	StateDownloading = "Downloading"                     //Torrent is being downloaded and data is being transfered
	StateStalledDL   = "Donwloading, connection waiting" //Torrent is being downloaded, but no connection were made
	StateMetaDL      = "Fetching Metadata"               //Torrent has just started downloading and is fetching metadata
)

//TorrentState ...
var TorrentState = &State{
	"error":       StateError,
	"pausedUP":    StatePausedUP,
	"pausedDL":    StatePausedDL,
	"queuedUP":    StateQueuedUP,
	"queuedDL":    StateQueuedDL,
	"uploading":   StateUploading,
	"stalledUP":   StateStalledUP,
	"checkingUP":  StateCheckingUP,
	"checkingDL":  StateCheckingDL,
	"downloading": StateDownloading,
	"stalledDL":   StateStalledDL,
	"metaDL":      StateMetaDL,
}
