package updater

import "gitlab-metrics/gitlab/dto"

var (
	InterestedActions = [...]string{
		dto.EventActionNameAccepted,
		dto.EventActionNameClosed,
		dto.EventActionNameCommented,
		dto.EventActionNameCreated,
		dto.EventActionNameMerged,
		dto.EventActionNameOpened,
		dto.EventActionNamePushedNew,
		dto.EventActionNamePushedTo,
		dto.EventActionNameReopened,
		dto.EventActionNameUpdated,
	}
)
