package service

//import (
//	"avito_shop/internal/repository"
//	"fmt"
//)
//
//type TransactionService struct {
//	userRepo        repository.UserRepository
//	transactionRepo repository.TransactionRepository
//}
//
//func NewTransactionService(userRepo repository.UserRepository, transactionRepo repository.TransactionRepository) *TransactionService {
//	return &TransactionService{
//		userRepo:        userRepo,
//		transactionRepo: transactionRepo,
//	}
//}
//
//// Отправить монеты от одного пользователя к другому
//func (s *TransactionService) SendCoins(senderID, recipientID int, amount float64) error {
//	// Получаем отправителя
//	sender, err := s.userRepo.GetUserByID(senderID)
//	if err != nil {
//		return fmt.Errorf("sender not found: %w", err)
//	}
//
//	// Проверяем баланс отправителя
//	if sender.Balance < amount {
//		return fmt.Errorf("insufficient balance")
//	}
//
//	// Получаем получателя
//	recipient, err := s.userRepo.GetUserByID(recipientID)
//	if err != nil {
//		return fmt.Errorf("recipient not found: %w", err)
//	}
//
//	// Обновляем баланс отправителя и получателя
//	sender.Balance -= amount
//	recipient.Balance += amount
//
//	// Сохраняем обновления в БД
//	err = s.userRepo.UpdateUser(sender)
//	if err != nil {
//		return fmt.Errorf("failed to update sender balance: %w", err)
//	}
//	err = s.userRepo.UpdateUser(recipient)
//	if err != nil {
//		return fmt.Errorf("failed to update recipient balance: %w", err)
//	}
//
//	// Записываем транзакцию
//	transaction := &models.Transaction{
//		SenderID:    senderID,
//		RecipientID: recipientID,
//		Amount:      amount,
//	}
//	err = s.transactionRepo.CreateTransaction(transaction)
//	if err != nil {
//		return fmt.Errorf("failed to create transaction: %w", err)
//	}
//
//	return nil
//}

//func (s *storeService) GetTransactionHistory(ctx context.Context, userID int) (*response.CoinHistory, error) {
//	// Получаем отправленные транзакции (userID -> другие пользователи)
//	sent, err := s.transRepo.GetSentTransactions(ctx, userID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Получаем полученные транзакции (другие пользователи -> userID)
//	received, err := s.transRepo.GetReceivedTransactions(ctx, userID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Формируем структуру ответа
//	history := &response.CoinHistory{
//		Sent:     make([]response.CoinTransaction, len(sent)),
//		Received: make([]response.CoinTransaction, len(received)),
//	}
//
//	// Заполняем отправленные транзакции
//	for i, t := range sent {
//		history.Sent[i] = response.CoinTransaction{
//			FromUser: strconv.Itoa(t.ToUserID), // Можно заменить на имя пользователя
//			Amount:   t.Amount,
//		}
//	}
//
//	// Заполняем полученные транзакции
//	for i, t := range received {
//		history.Received[i] = response.CoinTransaction{
//			FromUser: strconv.Itoa(t.FromUserID),
//			Amount:   t.Amount,
//		}
//	}
//
//	return history, nil
//}

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

// SendCoin обновление балансов пользователей
//func (s *storeService) SendCoin(ctx context.Context, userID, recipientID int, amount int) error {
//	// Получаем отправителя
//	user, err := s.userRepo.GetUserByID(ctx, userID)
//	if err != nil {
//		return fmt.Errorf("user not found: %w", err)
//	}
//
//	// Проверяем, достаточно ли средств у отправителя
//	if user.Balance < amount {
//		return fmt.Errorf("not enough balance")
//	}
//
//	// Получаем получателя
//	recipient, err := s.userRepo.GetUserByID(ctx, recipientID)
//	if err != nil {
//		return fmt.Errorf("recipient not found: %w", err)
//	}
//
//	// Обновляем баланс отправителя и получателя
//	user.Balance -= amount
//	recipient.Balance += amount
//
//	// Сохраняем обновленные данные пользователей в базе
//	err = s.userRepo.UpdateUser(ctx, user)
//	if err != nil {
//		return fmt.Errorf("failed to update sender: %w", err)
//	}
//
//	err = s.userRepo.UpdateUser(ctx, recipient)
//	if err != nil {
//		return fmt.Errorf("failed to update recipient: %w", err)
//	}
//
//	// Логика для записи транзакции в таблицу транзакций
//	return s.transactionRepo.CreateTransaction(ctx, userID, recipientID, amount)
//}
