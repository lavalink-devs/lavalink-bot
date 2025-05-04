package lavalinkbot

import (
	"math/rand"
	"sync"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type RepeatMode string

const (
	RepeatModeNone  RepeatMode = "none"
	RepeatModeTrack RepeatMode = "track"
	RepeatModeQueue RepeatMode = "queue"
)

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		queues: map[snowflake.ID]*queue{},
	}
}

type PlayerManager struct {
	queues map[snowflake.ID]*queue
	mu     sync.Mutex
}

type queue struct {
	tracks    []lavalink.Track
	mode      RepeatMode
	channelID snowflake.ID
}

func (q *PlayerManager) Get(guildID snowflake.ID) (RepeatMode, []lavalink.Track) {
	q.mu.Lock()
	defer q.mu.Unlock()

	qq, ok := q.queues[guildID]
	if !ok {
		return RepeatModeNone, nil
	}
	return qq.mode, qq.tracks
}

func (q *PlayerManager) Delete(guildID snowflake.ID) {
	q.mu.Lock()
	defer q.mu.Unlock()

	delete(q.queues, guildID)
}

func (q *PlayerManager) ChannelID(guildID snowflake.ID) snowflake.ID {
	q.mu.Lock()
	defer q.mu.Unlock()

	qu, ok := q.queues[guildID]
	if !ok {
		return 0
	}
	return qu.channelID
}

func (q *PlayerManager) Add(guildID snowflake.ID, channelID snowflake.ID, tracks ...lavalink.Track) {
	q.mu.Lock()
	defer q.mu.Unlock()

	qq, ok := q.queues[guildID]
	if !ok {
		qq = &queue{
			channelID: channelID,
		}
		q.queues[guildID] = qq
	}
	qq.tracks = append(qq.tracks, tracks...)
}

func (q *PlayerManager) Remove(guildID snowflake.ID, from int, to int) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	qq, ok := q.queues[guildID]
	if !ok {
		return false
	}

	queueLen := len(qq.tracks)
	if from >= queueLen || to >= queueLen {
		return false
	}

	if to == 0 {
		to = from + 1
	}

	qq.tracks = append(qq.tracks[:from], qq.tracks[to:]...)
	return true
}

func (q *PlayerManager) Clear(guildID snowflake.ID) {
	q.mu.Lock()
	defer q.mu.Unlock()

	delete(q.queues, guildID)
}

func (q *PlayerManager) Shuffle(guildID snowflake.ID) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	qq, ok := q.queues[guildID]
	if !ok {
		return false
	}

	if len(q.queues[guildID].tracks) >= 1 {
		return false
	}

	for i := range qq.tracks {
		j := i + rand.Intn(len(qq.tracks)-i)
		qq.tracks[i], qq.tracks[j] = qq.tracks[j], qq.tracks[i]
	}

	return true
}

func (q *PlayerManager) SetRepeatMode(guildID snowflake.ID, mode RepeatMode) {
	q.mu.Lock()
	defer q.mu.Unlock()

	qq, ok := q.queues[guildID]
	if !ok {
		return
	}
	qq.mode = mode
}

func (q *PlayerManager) Next(guildID snowflake.ID) (lavalink.Track, bool) {
	return q.NextCount(guildID, 1)
}

func (q *PlayerManager) NextCount(guildID snowflake.ID, count int) (lavalink.Track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	qq, ok := q.queues[guildID]
	if !ok {
		return lavalink.Track{}, false
	}
	if len(qq.tracks) < count {
		return lavalink.Track{}, false
	}

	track := qq.tracks[count-1]
	if qq.mode != RepeatModeTrack {
		if qq.mode == RepeatModeQueue {
			qq.tracks = append(qq.tracks, track)
		}
		qq.tracks = qq.tracks[count:]
	}
	return track, true
}
