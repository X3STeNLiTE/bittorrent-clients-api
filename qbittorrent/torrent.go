package qbittorrent

//Torrent ...
type Torrent struct {
	//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-Documentation#get-torrent-list
	Hash                   string  `json:"hash"`           //Torrent hash
	Name                   string  `json:"name"`           //Torrent name
	Size                   int64   `json:"size"`           //Total size (bytes) of files selected for download
	Progress               float32 `json:"progress"`       //Torrent progress (percentage/100)
	Dlspeed                int64   `json:"dlspeed"`        //Torrent download speed (bytes/s)
	Upspeed                int64   `json:"upspeed"`        //Torrent upload speed (bytes/s)
	Priority               int64   `json:"priority"`       //Torrent priority. Returns -1 if queuing is disabled or torrent is in seed mode
	NumSeeds               int64   `json:"num_seeds"`      //Number of seeds connected to
	NumComplete            int64   `json:"num_complete"`   //Number of seeds in the swarm
	NumLeechs              int64   `json:"num_leechs"`     //Number of leechers connected to
	NumIncomplete          int64   `json:"num_incomplete"` //Number of leechers in the swarm
	Ratio                  float32 `json:"ratio"`          //Torrent share ratio. Max ratio value: 9999.
	Eta                    int64   `json:"eta"`            //Torrent ETA (seconds)
	State                  string  `json:"state"`          //Torrent state. See table here below for the possible values
	SeqDl                  bool    `json:"seq_dl"`         //True if sequential download is enabled
	FirstLastPiecePriority bool    `json:"f_l_piece_prio"` //True if first last piece are prioritized
	Category               string  `json:"category"`       //Category of the torrent
	SuperSeeding           bool    `json:"super_seeding"`  //True if super seeding is enabled
	ForceStart             bool    `json:"force_start"`    //True if force start is enabled for this torrent
	DateAdded              uint64  `json:"added_on"`       //Not On Wiki: Timestamp of Added Time
	DateCompleted          uint64  `json:"completion_on"`  //Not On Wiki: Timestamp of Completed Time
	LastActivity           uint64  `json:"last_activity"`  //Not On Wiki: Timestamp of LastActivity Time
}

//Completed true if Progress == 1
func (T Torrent) Completed() bool {
	return T.Progress == 1
}
