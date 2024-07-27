package wikipedia

type PageRevisionSlotMain struct {
	ContentModel  string `json:"contentmodel"`
	ContentFormat string `json:"contentformat"`
	Content       string `json:"content"`
}

type PageRevisionSlot struct {
	Main PageRevisionSlotMain `json:"main"`
}

type PageRevision struct {
	Slots PageRevisionSlot `json:"slots"`
}

type Page struct {
	PageId    int            `json:"pageid"`
	Title     string         `json:"title"`
	Revisions []PageRevision `json:"revisions"`
}

type Query struct {
	Pages []Page `json:"pages"`
}
