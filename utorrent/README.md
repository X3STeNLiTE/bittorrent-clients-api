# Utorrent Client #

## Compatible with uTorrent 3.4.5 ##

### Listing file ###

```
client := utorrent.Clients("{UTORRENT_URL}", "{ADMIN}", "{PASSWORD}")
```

### Actions ###

```
// Start
func (T Torrent) Start(hash []string) error

// Stop
func (T Torrent) Stop(hash []string) error

// Pause
func (T Torrent) Pause(hash []string) error

// ForceStart
func (T Torrent) ForceStart(hash []string) error

// Unpause
func (T Torrent) Unpause(hash []string) error

// Recheck
func (T Torrent) Recheck(hash []string) error

// Remove
func (T Torrent) Remove(hash []string) error

// RemoveData
func (T Torrent) RemoveData(hash []string) error
```