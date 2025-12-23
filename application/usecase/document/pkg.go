package document

import "src/application/usecase/document/query/search_signed_document"

func Register() {
	search_signed_document.Register()
}
