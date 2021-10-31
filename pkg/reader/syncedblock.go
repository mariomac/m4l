package reader

import (
	"github.com/mariomac/msxmml/pkg/song"
)

type SyncedBlock struct {
	block    song.SyncedBlock
	counters map[string]channelCounter
	time     float64 // time in beats fractions
}

type channelCounter struct {
	time  float64 // time in beats fractions
	index int
}

func NewSyncedBlock(block song.SyncedBlock) SyncedBlock {
	counters := map[string]channelCounter{}
	for chn := range block.Channels {
		counters[chn] = channelCounter{}
	}
	return SyncedBlock{block: block, counters: counters}
}

// Next extracts the next item to be played/enqueued. Returns it as well as the channel where it belongs to.
// If there are no more items, returns empty channel
func (sbr *SyncedBlock) Next() (song.TablatureItem, string) {
	soonerChannel := ""
	for name, channel := range sbr.block.Channels {
		cnt := sbr.counters[name]
		if cnt.index >= len(channel.Items) {
			continue
		}
		if soonerChannel == "" || cnt.time < sbr.counters[soonerChannel].time {
			soonerChannel = name
		}
	}
	if soonerChannel == "" {
		return song.TablatureItem{}, ""
	} else {
		cnt := sbr.counters[soonerChannel]
		it := sbr.block.Channels[soonerChannel].Items[cnt.index]
		sbr.counters[soonerChannel] = channelCounter{
			index: cnt.index + 1,
			time:  cnt.time + itemDurationBeats(it),
		}
		return it, soonerChannel
	}
}

func itemDurationBeats(it song.TablatureItem) float64 {
	if it.Note != nil {
		return 4 / float64(it.Note.Length)
	}
	if it.Silence != nil {
		return 4 / float64(it.Silence.Length)
	}
	return 0
}
