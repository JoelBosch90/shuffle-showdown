package wikipedia_models

type RevisionSlotMain struct {
	ContentModel  string `json:"contentmodel"`
	ContentFormat string `json:"contentformat"`
	Content       string `json:"content"`
}

type RevisionSlot struct {
	Main RevisionSlotMain `json:"main"`
}

type Revision struct {
	Slots RevisionSlot `json:"slots"`
}
