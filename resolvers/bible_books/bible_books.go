package bible_books

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

func GetAllBibleBooks(p graphql.ResolveParams) (interface{}, error) {
	reader_code := p.Context.Value("language").(string)
	bible_books := new(models.BibleBook).GetAllBibleBooks(reader_code)

	return bible_books, nil
}

func GetBibleBookByCode(p graphql.ResolveParams) (interface{}, error) {
	reader_code := p.Context.Value("language").(string)
	code, ok := p.Args["code"].(string)
	if !ok {
		return nil, nil
	}

	bible_book := new(models.BibleBook).GetBibleBookByCode(reader_code, code)

	return bible_book, nil
}
