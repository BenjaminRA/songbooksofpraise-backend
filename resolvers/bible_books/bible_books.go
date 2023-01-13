package bible_books

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

func GetAllBibleBooks(p graphql.ResolveParams) (interface{}, error) {
	reader_code := p.Context.Value("language").(string)
	bible_books, err := new(models.BibleBook).GetAllBibleBooks(reader_code)
	if err != nil {
		return nil, err
	}

	return bible_books, nil
}

func GetBibleBookByCode(p graphql.ResolveParams) (interface{}, error) {
	reader_code := p.Context.Value("language").(string)
	code, ok := p.Args["code"].(string)
	if !ok {
		return nil, nil
	}

	bible_book, err := new(models.BibleBook).GetBibleBookByCode(reader_code, code)
	if err != nil {
		return nil, err
	}

	return bible_book, nil
}
