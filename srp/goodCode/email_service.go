package goodCode

type EmailService struct {
	repository EmailRepository
	sender     EmailSender
}

func NewEmailService(repository EmailRepository, sender EmailSender) *EmailService {
	return &EmailService{
		repository: repository,
		sender:     sender,
	}
}

func (s *EmailService) Send(from string, to string, subject string, message string) error {
	err := s.repository.Save(from, to, subject, message)
	if err != nil {
		return err
	}

	return s.sender.Send(from, to, subject, message)
}