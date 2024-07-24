package vo

import "time"

type ListMemoReq struct {
	Page            int        `json:"page,omitempty"`
	Size            int        `json:"size,omitempty"`
	Tag             string     `json:"tag,omitempty"`
	Source          string     `json:"source,omitempty"`
	Username        string     `json:"username,omitempty"`
	Start           *time.Time `json:"start,omitempty"`
	End             *time.Time `json:"end,omitempty"`
	ContentContains string     `json:"contentContains,omitempty"`
	ShowType        *int       `json:"showType,omitempty"`
}

type MemoExt struct {
	YoutubeUrl    string `json:"youtube_url,omitempty"`
	VideoUrl      string `json:"video_url,omitempty"`
	LocalVideoUrl string `json:"local_video_url,omitempty"`
}

type SaveMemoReq struct {
	ID              int      `json:"id,omitempty"`
	Content         string   `json:"content,omitempty"`
	Ext             MemoExt  `json:"ext"`
	Pinned          *bool    `json:"pinned,omitempty"`
	ShowType        int      `json:"show_type,omitempty"`
	ExternalFavicon string   `json:"externalFavicon,omitempty"`
	ExternalTitle   string   `json:"externalTitle,omitempty"`
	ExternalUrl     string   `json:"externalUrl,omitempty"`
	Imgs            []string `json:"imgs,omitempty"`
	Music163Url     string   `json:"music163Url,omitempty"`
	BilibiliUrl     string   `json:"bilibiliUrl,omitempty"`
	Location        string   `json:"location,omitempty"`
}