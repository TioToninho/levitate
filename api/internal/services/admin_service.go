package services

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"trackable-donations/api/internal/models"
)

// AdminService gerencia operações relacionadas a administração do sistema
type AdminService struct {
	donations        []models.Donation
	ngos             []models.NGO
	ngoRegistrations []models.NGORegistration
	auditLogs        []models.AuditLog
	donationService  *DonationService
	expenseService   *ExpenseService
}

// NewAdminService cria uma nova instância do serviço de administração
func NewAdminService(donationSvc *DonationService, expenseSvc *ExpenseService) *AdminService {
	return &AdminService{
		donations:        []models.Donation{},
		ngos:             []models.NGO{},
		ngoRegistrations: []models.NGORegistration{},
		auditLogs:        []models.AuditLog{},
		donationService:  donationSvc,
		expenseService:   expenseSvc,
	}
}

// RegisterNGO inicia o processo de registro de uma nova ONG
func (s *AdminService) RegisterNGO(req models.NGORegistrationRequest) (models.NGORegistration, error) {
	// Verificar se o CNPJ já está em uso
	for _, reg := range s.ngoRegistrations {
		if reg.CNPJ == req.CNPJ {
			return models.NGORegistration{}, errors.New("CNPJ já registrado no sistema")
		}
	}

	for _, ngo := range s.ngos {
		if ngo.CNPJ == req.CNPJ {
			return models.NGORegistration{}, errors.New("CNPJ já pertence a uma ONG ativa")
		}
	}

	// Validar o formato do CNPJ
	isValid, msg := s.validateCNPJFormat(req.CNPJ)

	registrationID := uint(len(s.ngoRegistrations) + 1)
	registration := models.NGORegistration{
		ID:                registrationID,
		Name:              req.Name,
		Description:       req.Description,
		Category:          req.Category,
		CNPJ:              req.CNPJ,
		CNPJValid:         isValid,
		CNPJValidationMsg: msg,
		Email:             req.Email,
		Phone:             req.Phone,
		Address:           req.Address,
		ResponsibleID:     req.ResponsibleID,
		LogoURL:           req.LogoURL,
		Status:            models.NGOStatusPending,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	s.ngoRegistrations = append(s.ngoRegistrations, registration)

	// Registrar ação no log de auditoria
	s.logAuditAction(0, "ngo_registration_created", "ngo_registration", registrationID, "",
		fmt.Sprintf("Registro de ONG solicitado: %s (CNPJ: %s)", req.Name, req.CNPJ))

	return registration, nil
}

// validateCNPJFormat valida o formato do CNPJ (somente verificação de formato)
func (s *AdminService) validateCNPJFormat(cnpj string) (bool, string) {
	// Remover caracteres não numéricos
	re := regexp.MustCompile(`[^0-9]`)
	cnpj = re.ReplaceAllString(cnpj, "")

	// Verificar se tem 14 dígitos
	if len(cnpj) != 14 {
		return false, "CNPJ deve conter 14 dígitos"
	}

	// Verificar se todos os dígitos são iguais (caso inválido)
	allEqual := true
	for i := 1; i < len(cnpj); i++ {
		if cnpj[i] != cnpj[0] {
			allEqual = false
			break
		}
	}

	if allEqual {
		return false, "CNPJ inválido: todos os dígitos são iguais"
	}

	// Calcular primeiro dígito verificador
	sum := 0
	weight := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	for i := 0; i < 12; i++ {
		digit := int(cnpj[i] - '0')
		sum += digit * weight[i]
	}

	remainder := sum % 11
	verificationDigit1 := 0
	if remainder >= 2 {
		verificationDigit1 = 11 - remainder
	}

	if int(cnpj[12]-'0') != verificationDigit1 {
		return false, "CNPJ inválido: primeiro dígito verificador incorreto"
	}

	// Calcular segundo dígito verificador
	sum = 0
	weight = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	for i := 0; i < 13; i++ {
		digit := int(cnpj[i] - '0')
		sum += digit * weight[i]
	}

	remainder = sum % 11
	verificationDigit2 := 0
	if remainder >= 2 {
		verificationDigit2 = 11 - remainder
	}

	if int(cnpj[13]-'0') != verificationDigit2 {
		return false, "CNPJ inválido: segundo dígito verificador incorreto"
	}

	return true, "CNPJ válido"
}

// ValidateCNPJOnline realiza uma validação online do CNPJ (simulado)
func (s *AdminService) ValidateCNPJOnline(registrationID uint) (models.NGORegistration, error) {
	// Encontrar o registro
	var registration models.NGORegistration
	var index int
	found := false

	for i, reg := range s.ngoRegistrations {
		if reg.ID == registrationID {
			registration = reg
			index = i
			found = true
			break
		}
	}

	if !found {
		return models.NGORegistration{}, errors.New("registro de ONG não encontrado")
	}

	// Em um ambiente real, faria uma consulta a um serviço externo
	// Aqui, simularemos com base na validação de formato
	if registration.CNPJValid {
		// Simulando consulta online bem-sucedida
		s.ngoRegistrations[index].CNPJValid = true
		s.ngoRegistrations[index].CNPJValidationMsg = "CNPJ verificado online e válido"
		s.ngoRegistrations[index].Status = models.NGOStatusValidating
		s.ngoRegistrations[index].UpdatedAt = time.Now()

		// Registrar ação no log de auditoria
		s.logAuditAction(0, "cnpj_validated", "ngo_registration", registrationID,
			string(registration.Status), string(models.NGOStatusValidating))

		return s.ngoRegistrations[index], nil
	} else {
		return models.NGORegistration{}, errors.New(registration.CNPJValidationMsg)
	}
}

// UploadNGODocuments simula o upload de documentos para o IPFS
func (s *AdminService) UploadNGODocuments(registrationID uint, fileContent []byte) (models.NGORegistration, error) {
	// Encontrar o registro
	var registration models.NGORegistration
	var index int
	found := false

	for i, reg := range s.ngoRegistrations {
		if reg.ID == registrationID {
			registration = reg
			index = i
			found = true
			break
		}
	}

	if !found {
		return models.NGORegistration{}, errors.New("registro de ONG não encontrado")
	}

	// Verificar se o CNPJ foi validado
	if !registration.CNPJValid {
		return models.NGORegistration{}, errors.New("CNPJ deve ser validado antes do upload de documentos")
	}

	// Simular upload para IPFS
	ipfsHash := fmt.Sprintf("Qm%s", generateMockHash(46))

	// Atualizar o registro
	s.ngoRegistrations[index].DocumentsIPFS = ipfsHash
	s.ngoRegistrations[index].UpdatedAt = time.Now()

	// Registrar ação no log de auditoria
	s.logAuditAction(0, "documents_uploaded", "ngo_registration", registrationID,
		"", fmt.Sprintf("Documentos enviados para IPFS: %s", ipfsHash))

	return s.ngoRegistrations[index], nil
}

// ApproveNGO aprova o registro de uma ONG e cria a entrada na blockchain
func (s *AdminService) ApproveNGO(registrationID uint, adminID uint, comments string) (models.NGO, error) {
	// Encontrar o registro
	var registration models.NGORegistration
	var regIndex int
	found := false

	for i, reg := range s.ngoRegistrations {
		if reg.ID == registrationID {
			registration = reg
			regIndex = i
			found = true
			break
		}
	}

	if !found {
		return models.NGO{}, errors.New("registro de ONG não encontrado")
	}

	// Verificar se todos os requisitos foram cumpridos
	if !registration.CNPJValid {
		return models.NGO{}, errors.New("CNPJ não foi validado")
	}

	if registration.DocumentsIPFS == "" {
		return models.NGO{}, errors.New("documentos não foram enviados")
	}

	// Simular registro na blockchain
	blockchainRef := generateMockTransactionHash()

	// Atualizar o registro
	s.ngoRegistrations[regIndex].BlockchainRef = blockchainRef
	s.ngoRegistrations[regIndex].Status = models.NGOStatusApproved
	s.ngoRegistrations[regIndex].AdminComments = comments
	s.ngoRegistrations[regIndex].UpdatedAt = time.Now()

	// Criar uma nova ONG
	ngoID := uint(len(s.ngos) + 1)
	ngo := models.NGO{
		ID:            ngoID,
		Name:          registration.Name,
		Description:   registration.Description,
		Category:      registration.Category,
		CNPJ:          registration.CNPJ,
		Email:         registration.Email,
		Phone:         registration.Phone,
		Address:       registration.Address,
		LogoURL:       registration.LogoURL,
		DocumentsIPFS: registration.DocumentsIPFS,
		BlockchainRef: blockchainRef,
		ResponsibleID: registration.ResponsibleID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	s.ngos = append(s.ngos, ngo)

	// Adicionar a ONG ao serviço de doações
	s.donationService.ngos = append(s.donationService.ngos, ngo)

	// Registrar ação no log de auditoria
	s.logAuditAction(adminID, "ngo_approved", "ngo", ngoID,
		string(models.NGOStatusValidating), string(models.NGOStatusApproved))

	return ngo, nil
}

// RejectNGO rejeita o registro de uma ONG
func (s *AdminService) RejectNGO(registrationID uint, adminID uint, reason string) (models.NGORegistration, error) {
	// Encontrar o registro
	var registration models.NGORegistration
	var index int
	found := false

	for i, reg := range s.ngoRegistrations {
		if reg.ID == registrationID {
			registration = reg
			index = i
			found = true
			break
		}
	}

	if !found {
		return models.NGORegistration{}, errors.New("registro de ONG não encontrado")
	}

	// Atualizar o registro
	s.ngoRegistrations[index].Status = models.NGOStatusRejected
	s.ngoRegistrations[index].AdminComments = reason
	s.ngoRegistrations[index].UpdatedAt = time.Now()

	// Registrar ação no log de auditoria
	s.logAuditAction(adminID, "ngo_rejected", "ngo_registration", registrationID,
		string(registration.Status), string(models.NGOStatusRejected))

	return s.ngoRegistrations[index], nil
}

// GetNGORegistrations retorna todos os registros de ONGs
func (s *AdminService) GetNGORegistrations() []models.NGORegistration {
	return s.ngoRegistrations
}

// GetNGORegistrationByID retorna um registro de ONG pelo ID
func (s *AdminService) GetNGORegistrationByID(registrationID uint) (models.NGORegistration, error) {
	for _, reg := range s.ngoRegistrations {
		if reg.ID == registrationID {
			return reg, nil
		}
	}
	return models.NGORegistration{}, errors.New("registro de ONG não encontrado")
}

// GetNGORegistrationsByCNPJ retorna registros de ONGs pelo CNPJ
func (s *AdminService) GetNGORegistrationsByCNPJ(cnpj string) []models.NGORegistration {
	var results []models.NGORegistration

	for _, reg := range s.ngoRegistrations {
		if reg.CNPJ == cnpj {
			results = append(results, reg)
		}
	}

	return results
}

// AuditEntity realiza auditoria em uma entidade (ONG, doação ou despesa)
func (s *AdminService) AuditEntity(req models.AuditRequest, adminID uint) (models.AuditResult, error) {
	result := models.AuditResult{
		EntityType:     req.EntityType,
		EntityID:       req.EntityID,
		ValidationDate: time.Now(),
	}

	var blockchainRef string
	var ipfsRef string
	var validationErrors []string

	switch req.EntityType {
	case "ngo":
		// Verificar se a ONG existe
		found := false
		for _, ngo := range s.ngos {
			if ngo.ID == req.EntityID {
				blockchainRef = ngo.BlockchainRef
				ipfsRef = ngo.DocumentsIPFS
				found = true
				break
			}
		}

		if !found {
			return result, errors.New("ONG não encontrada")
		}

	case "donation":
		// Verificar se a doação existe
		found := false
		for _, donation := range s.donationService.donations {
			if donation.ID == req.EntityID {
				blockchainRef = donation.TransactionHash
				found = true
				break
			}
		}

		if !found {
			return result, errors.New("doação não encontrada")
		}

		// Encontrar o recibo relacionado
		for _, receipt := range s.donationService.receipts {
			if receipt.DonationID == req.EntityID {
				ipfsRef = receipt.IPFSHash
				break
			}
		}

	case "expense":
		// Verificar se a despesa existe
		found := false
		for _, expense := range s.expenseService.expenses {
			if expense.ID == req.EntityID {
				blockchainRef = expense.BlockchainRef
				ipfsRef = expense.ReceiptIPFS
				found = true
				break
			}
		}

		if !found {
			return result, errors.New("despesa não encontrada")
		}

	default:
		return result, fmt.Errorf("tipo de entidade desconhecido: %s", req.EntityType)
	}

	// Verificar a validade na blockchain (simulado)
	blockchainValid := s.verifyBlockchainReference(blockchainRef)
	if !blockchainValid {
		validationErrors = append(validationErrors, "Referência na blockchain inválida ou não encontrada")
	}

	// Verificar a validade no IPFS (simulado)
	ipfsValid := s.verifyIPFSReference(ipfsRef)
	if !ipfsValid {
		validationErrors = append(validationErrors, "Referência no IPFS inválida ou não encontrada")
	}

	result.BlockchainValid = blockchainValid
	result.IPFSValid = ipfsValid
	result.BlockchainRef = blockchainRef
	result.IPFSRef = ipfsRef
	result.ValidationErrors = validationErrors

	// Registrar ação no log de auditoria
	comments := "Auditoria concluída com sucesso"
	if len(validationErrors) > 0 {
		comments = fmt.Sprintf("Auditoria com erros: %v", validationErrors)
	}

	s.logAuditAction(adminID, "audit_performed", req.EntityType, req.EntityID, "", comments)

	return result, nil
}

// verifyBlockchainReference verifica a validade de uma referência blockchain (simulado)
func (s *AdminService) verifyBlockchainReference(reference string) bool {
	// Em um ambiente real, verificaria a transação na blockchain
	// Aqui, verificamos apenas se o formato parece válido
	if reference == "" {
		return false
	}

	// Verificar se começa com "0x"
	if len(reference) < 2 || reference[:2] != "0x" {
		return false
	}

	// Verificar se tem o comprimento adequado (0x + 64 caracteres hexadecimais)
	if len(reference) != 66 {
		return false
	}

	// Verificar se todos os caracteres após "0x" são hexadecimais
	hexPattern := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return hexPattern.MatchString(reference)
}

// verifyIPFSReference verifica a validade de uma referência IPFS (simulado)
func (s *AdminService) verifyIPFSReference(reference string) bool {
	// Em um ambiente real, verificaria se o arquivo existe no IPFS
	// Aqui, verificamos apenas se o formato parece válido
	if reference == "" {
		return false
	}

	// Verificar se começa com "Qm"
	if len(reference) < 2 || reference[:2] != "Qm" {
		return false
	}

	// Verificar se tem o comprimento adequado (aproximadamente Qm + 44 caracteres base58)
	if len(reference) < 46 {
		return false
	}

	return true
}

// GetAuditLogs retorna os logs de auditoria
func (s *AdminService) GetAuditLogs() []models.AuditLog {
	return s.auditLogs
}

// GetAuditLogsByEntityType retorna logs de auditoria por tipo de entidade
func (s *AdminService) GetAuditLogsByEntityType(entityType string) []models.AuditLog {
	var logs []models.AuditLog

	for _, log := range s.auditLogs {
		if log.EntityType == entityType {
			logs = append(logs, log)
		}
	}

	return logs
}

// GetAuditLogsByEntityID retorna logs de auditoria por ID de entidade
func (s *AdminService) GetAuditLogsByEntityID(entityType string, entityID uint) []models.AuditLog {
	var logs []models.AuditLog

	for _, log := range s.auditLogs {
		if log.EntityType == entityType && log.EntityID == entityID {
			logs = append(logs, log)
		}
	}

	return logs
}

// logAuditAction registra uma ação de auditoria
func (s *AdminService) logAuditAction(adminID uint, action string, entityType string, entityID uint,
	previousState string, newState string) {

	logID := uint(len(s.auditLogs) + 1)
	log := models.AuditLog{
		ID:            logID,
		AdminID:       adminID,
		Action:        action,
		EntityType:    entityType,
		EntityID:      entityID,
		PreviousState: previousState,
		NewState:      newState,
		Comments:      newState, // Usando o newState como comentário para simplificar
		CreatedAt:     time.Now(),
	}

	s.auditLogs = append(s.auditLogs, log)
}
