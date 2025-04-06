package service

// проверка подключения к БД
func (s *Service) HealthCheck() error {
	return s.storage.Ping()
}
