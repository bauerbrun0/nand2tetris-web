package chips

type Chip struct {
	Inputs  map[string]IO
	Outputs map[string]IO
}

type IO struct {
	Width int
}
