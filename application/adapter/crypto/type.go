package crypto

// ICryptoAdapter define o contrato genérico de operações criptográficas usadas pela aplicação.
// Implementações concretas (ex.: AES, libs de hash, etc.) ficam em infra/crypto/*.
type ICryptoAdapter interface {
	// OTP gera um código de uso único (one-time password / token),
	// normalmente aleatório e de curta duração, representado como string.
	// Não faz persistência nem validação, apenas gera o valor.
	OTP() string

	// Hash gera um hash não reversível do texto em claro.
	// Deve ser usado para senhas ou dados sensíveis onde apenas a comparação é necessária.
	// A implementação é responsável por usar algoritmo, salt e custo adequados.
	Hash(plainText string) string

	// Encrypt serializa e criptografa o valor informado.
	//
	// - plainText pode ser qualquer valor serializável (string, struct, map, etc.).
	// - optionalIV, quando vazio, permite que a implementação gere um IV aleatório.
	//   Quando informado, a implementação pode usar esse IV de forma determinística.
	//
	// O retorno é o texto cifrado em formato string (ex.: base64).
	Encrypt(plainText any, optionalIV ...string) string

	// Decrypt descriptografa o texto cifrado e reconstrói o valor original.
	//
	// - cipherText deve ser o valor retornado previamente por Encrypt.
	// - o retorno é um valor genérico (any); o chamador deve fazer type assertion
	//   para o tipo esperado.
	//
	// Retorna erro quando o conteúdo não puder ser descriptografado ou validado.
	Decrypt(cipherText string) (any, error)
}
