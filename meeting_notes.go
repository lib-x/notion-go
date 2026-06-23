package notion

import "context"

// MeetingNotesService implements Notion meeting notes endpoints.
type MeetingNotesService struct {
	client *Client
}

// Query queries meeting notes.
func (s *MeetingNotesService) Query(ctx context.Context, body Object) (*MeetingNotesResponse, error) {
	var out MeetingNotesResponse
	err := s.client.post(ctx, apiPath("v1", "blocks", "meeting_notes", "query"), body, &out)
	return &out, err
}
