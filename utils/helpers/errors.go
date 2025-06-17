package helpers

const (
	// Erros gerais de requisição
	ErrDecodingRequest    = "Falha ao decodificar o corpo da requisição"
	ErrInvalidRequestBody = "Corpo da requisição inválido"
	ErrInvalidPathParam   = "Parâmetro de caminho inválido"
	ErrMethodNotAllowed   = "Método não permitido"
	ErrNotFound           = "Não encontrado"

	// Erros de validação de dados
	ErrMissingRequiredFields = "Campos obrigatórios ausentes"
	ErrInvalidEmailFormat    = "Formato de email inválido"
	ErrPasswordMismatch      = "As senhas não coincidem"
	ErrInvalidID             = "ID inválido"
	ErrInvalidTokenData      = "Dados do token inválidos"
	ErrValidationCodeError   = "Erro ao buscar código de validação"
	ErrValidationCodeNoteFound = "Código de validação não encontrado"
	ErrValidationCodeInvalid = "Código de validação inválido"
	ErrValidationCodeExpired = "Código expirado, gere um novo"
	ErrValidationCodeNotRegistered = "Falha ao registrar código de validação"
	ErrValidationCodeNotDeleted = "Falha ao deletar o código de validação"
	ErrPassordNotUpdated     = "Erro ao atualizar o password do usuário"

	// Erros de autenticação e autorização
	ErrUserHasAlreadyBeenRegistered = "E-mail já registrado"
	ErrUserOrPasswordInvalid = "Usuário ou senha inválidos"
	ErrUnauthorized          = "Não autorizado"
	ErrInvalidOldPassword    = "Senha antiga inválida"
	ErrTokenExpired          = "Token expirado"
	ErrTokenMalformed        = "Token malformado"
	ErrTokenSignatureInvalid = "Assinatura do token inválida"
	ErrTokenNotGenerated	 = "Erro ao gerar token"
	ErrTemporaryHashNotRegistered = "Falha ao salvar hash temporário"
	ErrHashNotGenerated      = "Falha ao gerar o hash da senha"

	// Erro de CSRF
	ErrMissingCSRFTokenHeader = "Cabeçalho do token CSRF ausente"
	ErrMissingCSRFTokenCookie = "Cookie do token CSRF ausente"
	ErrInvalidCSRFToken       = "Token CSRF inválido"
	ErrCSRFTokenMismatch      = "Incompatibilidade de token CSRF"
	ErrGeneratingCSRFToken    = "Erro ao gerar token CSRF"

	// Erros de usuário
	ErrUserBadRequest 		  = "Erro ao buscar usuário"
	ErrUserNotFound			  = "Usuário não encontrado"
	ErrUsersNotFound		  = "Nenhum usuário encontrado"
	ErrUsersIndex			  = "Erro ao listar usuários"
	ErrEmailAlreadyExists     = "Este email já está cadastrado"
	ErrFailedToCreateUser     = "Falha ao criar usuário"
	ErrFailedToUpdatePassword = "Falha ao atualizar senha do usuário"
	ErrUserShow               = "Erro ao exibir usuário"
	ErrUserUpdate             = "Erro ao atualizar usuário"

	// Erros de banco de dados
	ErrDatabaseError  = "Erro no banco de dados"
	ErrDatabaseQuery  = "Erro ao executar consulta no banco de dados"
	ErrDatabaseInsert = "Erro ao inserir dados no banco de dados"
	ErrDatabaseUpdate = "Erro ao atualizar dados no banco de dados"
	ErrDatabaseDelete = "Erro ao deletar dados do banco de dados"
	ErrNoRows         = "Nenhum registro encontrado"

	// Erros internos do servidor
	ErrInternalServer  = "Erro interno do servidor"
	ErrPasswordHashing = "Erro ao gerar hash da senha"
	ErrTokenGeneration = "Erro ao gerar token JWT"
)