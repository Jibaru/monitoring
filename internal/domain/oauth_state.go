package domain

type OAuthState struct {
	id    ID
	state string
}

func NewOAuthState(id ID, state string) (*OAuthState, error) {
	return &OAuthState{
		id:    id,
		state: state,
	}, nil
}

func (o *OAuthState) ID() ID {
	return o.id
}

func (o *OAuthState) State() string {
	return o.state
}
