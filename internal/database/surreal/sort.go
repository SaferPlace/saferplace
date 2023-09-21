package surreal

import "api.safer.place/incident/v1"

// ByCommentTimestamp sorts the comments by timestamp, from the oldest to the newest
type ByCommentTimestamp []*incident.Comment

func (s ByCommentTimestamp) Len() int           { return len(s) }
func (s ByCommentTimestamp) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByCommentTimestamp) Less(i, j int) bool { return s[i].Timestamp < s[j].Timestamp }
